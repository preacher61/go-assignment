package main

import (
	"context"
	"preacher61/go-assignment/repository"

	"github.com/rs/zerolog/log"
)

type worker struct {
	persistResponses  func(ctx context.Context) int
	logUniqActivities func(ctx context.Context)
}

func newWorker() *worker {
	ri := newResponseInserter()
	ar := repository.NewActivityRepository()
	return &worker{
		persistResponses:  ri.processInsertion,
		logUniqActivities: ar.GetUniqActivitesCount,
	}
}

func (w *worker) run(ctx context.Context) {
	log.Info().Msg("Work started")
	defer func() {
		if err := recover(); err != nil {
			log.Error().Interface("error", err).Msg("worker panicked.....will retry")
		}
	}()

	keysProcessed := w.persistResponses(ctx)
	if keysProcessed < 1 {
		log.Info().Msg("No keys found in cache......suspending worker.!")
		return
	}
	w.logUniqActivities(ctx)
	log.Info().Msg("Work completed")
}
