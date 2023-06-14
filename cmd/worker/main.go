package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	tickerDuration := 15 * time.Second
	ticker := time.NewTicker(tickerDuration)
	defer ticker.Stop()

	ctx := context.Background()
	cancelCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	w := newWorker()

	log.Info().Msg("Worker starrted.....")
	for {
		select {
		case <-done:
			log.Info().Msg("Worker Shutting Down...!")
			return
		case <-ticker.C:
			w.run(cancelCtx)
			ticker.Stop()
			ticker = time.NewTicker(tickerDuration)
		}
	}

}
