package main

import (
	"context"
	"go-log-keeper/config"
	"go-log-keeper/services"
	"go-log-keeper/internal/db"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
	"fmt"
)

func main() {
	err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	// create a context for managing cancellation
	ctx, cancel := context.WithCancel(context.Background())

	var wg sync.WaitGroup


	if config.PrintHeartbeat {
		// Add heartbeat goroutine
		wg.Add(1)
		go func() {
			defer wg.Done()
			ticker := time.NewTicker(time.Duration(config.LogPrintingDelay) * time.Second)
			defer ticker.Stop()

			for {
				select {
				case <-ticker.C:
					log.Println("Service is running - Active workers:", config.NumberOfWorkers)
				case <-ctx.Done():
					return
				}
			}
		}()
	}

	if config.StartHealthCheckServer {
		// Add HTTP server
		wg.Add(1)
		go func() {
			defer wg.Done()
			redisClient := db.NewRedisClient()
			// Define handlers
			http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
				// Check Redis connection
				redisStatus := "up"
				if err := redisClient.Ping(ctx).Err(); err != nil {
					redisStatus = "down"
				}

				// Set JSON response header
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				// Write JSON response
				fmt.Fprintf(w, `{"status": "healthy", "workers": %d, "redis": "%s"}`,
					config.NumberOfWorkers,
					redisStatus)
			})

			server := &http.Server{
				Addr: ":80",
			}

			// Start server in a goroutine
			go func() {
				if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
					log.Printf("HTTP server error: %v", err)
				}
			}()

			// Wait for context cancellation
			<-ctx.Done()

			// Create shutdown context with timeout
			shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer shutdownCancel()

			// Shutdown the server gracefully
			if err := server.Shutdown(shutdownCtx); err != nil {
				log.Printf("HTTP server shutdown error: %v", err)
			}
		}()
	}



	// recover processing queue items
	services.RecoverProcessingItems(ctx)

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
