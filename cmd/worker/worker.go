package main

import (
	"context"

	"github.com/rs/zerolog/log"
)

type worker struct {
	persistResponses  func(ctx context.Context)
	logUniqActivities func(ctx context.Context)
}

func newWorker() *worker {
	ri := newResponseInserter()
	return &worker{
		persistResponses: ri.processInsertion,
	}
}

func (w *worker) run(ctx context.Context) {
	log.Info().Msg("Work started")
	defer func() {
		if err := recover(); err != nil {
			log.Error().Interface("error", err).Msg("worker panicked.....will retry")
		}
	}()

	w.persistResponses(ctx)
	w.logUniqActivities(ctx)
	log.Info().Msg("Work completed")
}
