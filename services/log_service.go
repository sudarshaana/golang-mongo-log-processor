package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-log-keeper/config"
	"go-log-keeper/internal/db"
	"go-log-keeper/models"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

func ProcessLog(ctx context.Context, workerID int, wg *sync.WaitGroup) {
	defer wg.Done()

	// Initialize clients
	redisClient := db.NewRedisClient()
	mongoClient, mongoCollection := db.NewMongoClient()

	// Ping Redis to ensure connection
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
		return
	}

	// Ping MongoDB to verify connection
	if err := mongoClient.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatalf("Could not ping MongoDB: %v", err)
		return
	}

	// Disconnect MongoDB client when done
	defer func() {
		if err := mongoClient.Disconnect(ctx); err != nil {
			log.Printf("Worker %d: Error disconnecting MongoDB client: %v", workerID, err)
		}
	}()
	defer func(redisClient *redis.Client) {
		err := redisClient.Close()
		if err != nil {
			log.Printf("Worker %d: Error disconnecting MongoDB redisClient: %v", workerID, err)
		}
	}(redisClient)

	// Process logs from Redis
	for {
		select {
		case <-ctx.Done():
			// Graceful shutdown: Context is cancelled, exit the loop
			log.Printf("Worker %d: Shutting down log processor...", workerID)

			return

		default:
			result, err := redisClient.BRPop(ctx, 0*time.Second, config.RedisQueueName).Result()
			if err != nil {
				log.Printf("Worker %d: Error fetching from Redis: %v", workerID, err)
				time.Sleep(1 * time.Second) // Wait before retrying
				continue
			}
			if len(result) == 0 {
				continue
			}

			logItem := result[1] // the log item popped
			uniqueKey := fmt.Sprintf("worker:%d:%s", workerID, time.Now().Format(time.RFC3339Nano))

			// Store the item in a Redis set with a unique key
			if err := redisClient.HSet(ctx, config.RedisProcessingQueueName, uniqueKey, logItem).Err(); err != nil {
				log.Printf("Worker %d: Error storing item in set: %v", workerID, err)
			}

			// Unmarshal the log data into RequestLog struct
			var logEntry models.RequestLog
			err = json.Unmarshal([]byte(logItem), &logEntry)

			if err != nil {
				RemoveFromProcessingSet(ctx, redisClient, uniqueKey)
				// TODO: Add a retry mechanism for failed log processing
				continue
			}

			success := storeLog(ctx, mongoCollection, logEntry, workerID)
			if success {
				RemoveFromProcessingSet(ctx, redisClient, uniqueKey)

			} else {
				// Move back to queue with some retry value for failed log processing
				// MoveBackToMainQueue(ctx, redisClient, logEntry)
			}
		}
	}

}
func storeLog(ctx context.Context, mongoCollection *mongo.Collection, logEntry models.RequestLog, workerID int) bool {
	_, err := mongoCollection.InsertOne(ctx, logEntry)
	if err != nil {
		log.Printf("Error inserting log into MongoDB: %v", err)
		return false
	}
	log.Printf("Worker %d: Successfully inserted log for %s %s into MongoDB.", workerID, logEntry.Method, logEntry.Path)
	return true
}
