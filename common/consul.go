package common

import (
	"context"
	"fmt"
	"sync"

	capi "github.com/hashicorp/consul/api"
	"github.com/rs/zerolog"
	"golang.org/x/sync/semaphore"
)

type getServicesFromConsulConfig struct {
	Client             *capi.Client
	NomadServicePrefix string
	ServiceName        string
	AllocationId       string // will replace the pattern matching any alloc id if specified
}

// returns all of the consul services that are associated with nomad
func getServicesFromConsul(config getServicesFromConsulConfig) (map[string][]string, error) {
	var uuidPattern string
	if config.AllocationId != "" {
		uuidPattern = config.AllocationId
	} else {
		uuidPattern = uuidRegexString
	}

	serviceIdRegex := fmt.Sprintf("^%s-%s", config.NomadServicePrefix, uuidPattern)
	var serviceNameClause string

	if config.ServiceName != "" {
		serviceNameClause = fmt.Sprintf(" and ServiceName == \"%s\"", config.ServiceName)
	}

	svcMap, _, err := config.Client.Catalog().Services(&capi.QueryOptions{
		Filter: fmt.Sprintf("ServiceID matches \"^%s\"%s", serviceIdRegex, serviceNameClause),
	})
	if err != nil {
		return nil, err
	}

	return svcMap, nil
}

type populateServicesConfig struct {
	Client        *capi.Client
	MaxConcurrent int64
	Services      map[string][]string
}

func populateServices(ctx context.Context, config populateServicesConfig) ([]*capi.CatalogService, error) {
	logger := zerolog.Ctx(ctx)
	sem := semaphore.NewWeighted(config.MaxConcurrent)
	results := make(chan []*capi.CatalogService, len(config.Services))
	var wg sync.WaitGroup

	var services []*capi.CatalogService

loop:
	for key := range config.Services {
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
				go func(key string) {
					defer sem.Release(1)
					defer wg.Done()

					res, _, err := config.Client.Catalog().Service(key, "", &capi.QueryOptions{})

					if err != nil {
						logger.Warn().
							AnErr("error", err).
							Str("service", key).
							Msg("Failed to retrieve information about the service")
						return
					}

					results <- res
				}(key)
			}
		}
	}

	wg.Wait()
	close(results)

	for svc := range results {
		services = append(services, svc...)
	}

	return services, nil
}

// deregisters service from Consul
func deregisterService(client *capi.Client, serviceId string, node string) (*capi.WriteMeta, error) {
	return client.Catalog().Deregister(&capi.CatalogDeregistration{ServiceID: serviceId, Node: node}, &capi.WriteOptions{})
}
