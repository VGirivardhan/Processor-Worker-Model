package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.Println("Starting URL Processor Application with Random Task Generation")

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	workerManagerConfig := WorkerManagerConfig{
		NumWorkers:       3,
		TaskQueueBuffer:  100,
		WorkerNamePrefix: "URLProcessor", // Naming pattern for workers
	}

	workerManager := NewWorkerManager(workerManagerConfig, httpClient)

	producer := NewTaskProducer(workerManager)
	producer.Start()

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	log.Println("Application running. Press Ctrl+C to stop...")
	<-c

	log.Println("Shutting down...")
	producer.Stop()
	log.Println("Application stopped")
}
