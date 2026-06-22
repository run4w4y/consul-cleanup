package common

import (
	"context"
	"sync"
	"time"

	capi "github.com/hashicorp/consul/api"
	napi "github.com/hashicorp/nomad/api"
	"github.com/rs/zerolog"
	"golang.org/x/exp/maps"
	"golang.org/x/sync/semaphore"
)

type DetermineOrphansConfig struct {
	Client        *napi.Client
	MaxConcurrent int64
	Services      []*capi.CatalogService
}

func DetermineOrphans(ctx context.Context, config DetermineOrphansConfig) []*capi.CatalogService {
	logger := zerolog.Ctx(ctx)
	sem := semaphore.NewWeighted(config.MaxConcurrent)
	results := make(chan struct {
		check   bool
		service *capi.CatalogService
	}, len(config.Services))
	var wg sync.WaitGroup

loop:
	for _, svc := range config.Services {
		select {
		case <-ctx.Done():
			break loop
		default:
			{
				if err := sem.Acquire(ctx, 1); err != nil {
					logger.Warn().
						AnErr("error", err).
						Msg("Failed to acquire semaphore")
					continue
				}

				wg.Add(1)
				go func(svc *capi.CatalogService) {
					defer sem.Release(1)
					defer wg.Done()

					logger.Info().
						Str("serviceId", svc.ServiceID).
						Msg("Checking service")
					check, err := checkServiceWithNomad(config.Client, svc)

					if err != nil {
						logger.Warn().
							AnErr("error", err).
							Str("serviceId", svc.ServiceID).
							Msg("Error when retrieving the Nomad alloc")

						results <- struct {
							check   bool
							service *capi.CatalogService
						}{true, svc}
						return
					}

					results <- struct {
						check   bool
						service *capi.CatalogService
					}{check, svc}
				}(svc)
			}
		}
	}

	wg.Wait()
	close(results)

	var orphans []*capi.CatalogService

	for res := range results {
		if !res.check {
			orphans = append(orphans, res.service)
		}
	}

	logger.Info().
		Int("count", len(orphans)).
		Msg("Orphan service entries found")

	return orphans
}

// one-off run of the cleanup
func OneshotCleanup(ctx context.Context, config OneshotCleanupConfig) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().Msg("Retrieving all of the services from Consul that are associated with a Nomad allocation")
	servicesMap, err := getServicesFromConsul(getServicesFromConsulConfig{
		Client:             config.Clients.Consul,
		NomadServicePrefix: config.NomadServicePrefix,
		ServiceName:        config.ServiceName,
	})
	if err != nil {
		return err
	}

	logger.Info().Int("count", len(servicesMap)).Msg("Retrieved services from Consul")

	services, err := populateServices(ctx, populateServicesConfig{
		Client:        config.Clients.Consul,
		Services:      servicesMap,
		MaxConcurrent: config.MaxConcurrent,
	})

	if err != nil {
		return err
	}

	orphans := DetermineOrphans(ctx, DetermineOrphansConfig{
		Client:        config.Clients.Nomad,
		MaxConcurrent: config.MaxConcurrent,
		Services:      services,
	})

	for _, svc := range orphans {
		logger.Info().
			Str("serviceId", svc.ServiceID).
			Msg("Deregistering service from Consul")

		if _, err := deregisterService(config.Clients.Consul, svc.ServiceID, svc.Node); err != nil {
			logger.Warn().
				AnErr("error", err).
				Str("serviceId", svc.ServiceID).
				Msg("Failed deregistering service from Consul")
		}
	}

	return nil
}

func CleanupListenerTask(ctx context.Context, config CleanupConfig) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().
		Msg("Attaching to the Nomad event stream to listen to stopped allocations")

	events, err := ReadAllocationsFromEventStream(ctx, config.Clients.Nomad)

	if err != nil {
		return err
	}

	var wg sync.WaitGroup

loop:
	for allocIds := range events { // TODO: make use of MaxConcurrency here? or limit to just one task running at a time?
		task := func(allocIds []string) { // TODO: add a set keeping track of allocations cleared
			defer wg.Done()

			for _, allocId := range allocIds {
				servicesMap, err := getServicesFromConsul(getServicesFromConsulConfig{
					Client:             config.Clients.Consul,
					NomadServicePrefix: config.NomadServicePrefix,
					AllocationId:       allocId,
				})

				if err != nil {
					logger.Warn().
						AnErr("error", err).
						Str("allocationId", allocId).
						Msg("Failed getting Consul services for allocation")
					continue
				}

				services, err := populateServices(ctx, populateServicesConfig{
					Client:        config.Clients.Consul,
					Services:      servicesMap,
					MaxConcurrent: config.MaxConcurrent,
				})

				if err != nil {
					logger.Warn().
						AnErr("error", err).
						Str("allocationId", allocId).
						Strs("serviceList", maps.Keys(servicesMap)).
						Msg("Failed populating Consul services")
					continue
				}

				for _, svc := range services {
					logger.Info().
						Str("serviceId", svc.ServiceID).
						Msg("Deregistering service from Consul")
					if _, err := deregisterService(config.Clients.Consul, svc.ServiceID, svc.Node); err != nil {
						logger.Warn().
							AnErr("error", err).
							Str("serviceId", svc.ServiceID).
							Msg("Failed deregistering service from Consul")
					}
				}
			}
		}

		select {
		case <-ctx.Done():
			break loop
		default:
			{
				wg.Add(1)
				go task(allocIds.Slice())
			}
		}
	}

	wg.Wait()

	return nil
}

func PeriodicCleaningTask(ctx context.Context, config PeriodicCleanupConfig) error {
	logger := zerolog.Ctx(ctx)
	logger.Info().
		Uint("interval", config.Interval).
		Msg("Started periodic cleanup task")

	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()

	done := make(chan struct{})

	go func() {
	loop:
		for {
			if err := OneshotCleanup(ctx, OneshotCleanupConfig{
				CleanupConfig: config.CleanupConfig,
			}); err != nil {
				logger.Warn().
					AnErr("error", err).
					Msg("Periodic cleanup failed")
			}

			select {
			case <-ticker.C:
				continue
			case <-ctx.Done():
				{
					done <- struct{}{}
					break loop
				}
			}
		}
	}()

	<-done

	return nil
}
