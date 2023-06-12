package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Worker Initiated..")

	<-done
	log.Println("Worked Stopped....!")
}
