package main

import (
	"context"
	"go-log-keeper/config"
	"go-log-keeper/services"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	// create a context for managing cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// recover processing queue items
	services.RecoverProcessingItems(ctx)

	// WaitGroup to wait for goroutines
	var wg sync.WaitGroup

	// define the number of goroutines to run
	numsWorkers := config.NumberOfWorkers

	// start multiple workers to process logs
	for i := 0; i < numsWorkers; i++ {
		wg.Add(1)
		go services.ProcessLog(ctx, i, &wg)
	}

	// handle or signal for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// block until a signal is received
	sig := <-signalChan
	log.Printf("received signal: %v. Initiating graceful shutdown...", sig)

	// cancel the process to stop the workers
	cancel()

	// Wait for all goroutines to stop
	wg.Wait()
	log.Println("All workers have stopped. Exiting application.")

}
