package common

import (
	"context"
	"math"
	"slices"

	"github.com/hashicorp/go-set/v2"
	"github.com/rs/zerolog"

	capi "github.com/hashicorp/consul/api"
	napi "github.com/hashicorp/nomad/api"
)

// gets nomad alloc by id
func getNomadAlloc(client *napi.Client, allocId string) (*napi.Allocation, error) {
	alloc, _, err := client.Allocations().Info(allocId, &napi.QueryOptions{})
	return alloc, err
}

// checks if an allocation is in either running or pending status
func checkAllocationStatus(status string) bool {
	return slices.Contains([]string{"running", "pending"}, status)
}

// returns false if the service appears to be orphan, true otherwise (has a running/pending associated nomad alloc)
func checkServiceWithNomad(client *napi.Client, service *capi.CatalogService) (bool, error) {
	allocId := extractUUID(service.ServiceID)
	alloc, err := getNomadAlloc(client, allocId)

	if err != nil {
		return false, err
	}

	if !checkAllocationStatus(alloc.ClientStatus) {
		return false, nil
	}

	return true, nil
}

// returns a channel to which will be pushed events related to allocations that can potentially leave orphan service entries in consul
func ReadAllocationsFromEventStream(ctx context.Context, client *napi.Client) (<-chan set.Set[string], error) {
	logger := zerolog.Ctx(ctx)
	events, err := client.EventStream().Stream(
		ctx,
		map[napi.Topic][]string{
			napi.TopicAllocation: {"*"},
		},
		math.MaxInt64, // to attach to the last received event
		&napi.QueryOptions{},
	)

	if err != nil {
		return nil, err
	}

	res := make(chan set.Set[string], 1)

	go func(events <-chan *napi.Events) {
		defer close(res)

	loop:
		for {
			select {
			case batch, ok := <-events:
				{
					if !ok {
						break loop
					}
					if batch == nil {
						continue
					}

					allocIds := set.New[string](1)

					for _, event := range batch.Events {
						if event.Topic != "Allocation" { // we only care about the Allocation events
							continue
						}

						alloc, err := event.Allocation()

						if err != nil {
							continue
						}

						if event.Type == "AllocationUpdated" && alloc.ClientStatus == "complete" {
							allocIds.Insert(alloc.ID)
							logger.Info().
								Str("allocationId", alloc.ID).
								Str("jobId", alloc.JobID).
								Msg("Potential allocation to be deregistered found")
						}
					}

					if allocIds.Size() > 0 {
						res <- *allocIds
					}
				}
			case <-ctx.Done():
				break loop
			}
		}
	}(events)

	return res, nil
}
