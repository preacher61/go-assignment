package main

import (
	"context"
	"preacher61/go-assignment/cache"
	"preacher61/go-assignment/model"
	"preacher61/go-assignment/repository"
	"sync"

	"github.com/rs/zerolog/log"
)

type responseInserter struct {
	iterateResponses func(ctx context.Context) <-chan map[string][]*model.Activity
	insert           func(ctx context.Context, data []*model.Activity) error
	deleteKeys       func(ctx context.Context, keys []string)
}

func newResponseInserter() *responseInserter {
	rh := cache.NewRedisHandler()
	ph := repository.NewActivityRepository()
	return &responseInserter{
		iterateResponses: rh.Get,
		insert:           ph.InsertActivities,
		deleteKeys:       rh.DeleteMulti,
	}
}

func (r *responseInserter) processInsertion(ctx context.Context) int {
	log.Info().Msg("Initiating response insertion.....")
	var wg sync.WaitGroup

	keysToDelete := []string{}

	for keyResMap := range r.iterateResponses(ctx) {
		wg.Add(1)

		go func(vMap map[string][]*model.Activity) {
			defer wg.Done()

			select {
			case <-ctx.Done():
				return
			default:
				for key, val := range vMap {
					err := r.insert(ctx, val)
					if err != nil {
						log.Error().Err(err).Msgf("insertion failed for key: %s ....skipping", key)
						return
					}
					keysToDelete = append(keysToDelete, key)
				}
			}
		}(keyResMap)
	}

	wg.Wait()

	if len(keysToDelete) < 1 {
		return 0
	}

	log.Info().Msg("response insertion completed......!")
	r.deleteKeys(ctx, keysToDelete)
	return len(keysToDelete)
}
