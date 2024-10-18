package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-log-keeper/config"
	"go-log-keeper/internal/db"
	"go-log-keeper/models"
	"log"
)

func RemoveFromProcessingSet(ctx context.Context, redisClient *redis.Client, uniqueKey string) {
	if err := redisClient.HDel(ctx, config.RedisProcessingQueueName, uniqueKey).Err(); err != nil {
		log.Printf("Error removing item from set after successful reprocessing: %v", err)
	}
}

func RecoverProcessingItems(ctx context.Context) {
	redisClient := db.NewRedisClient()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	items, err := redisClient.HGetAll(ctx, config.RedisProcessingQueueName).Result()
	if err != nil {
		log.Printf("Error retrieving processing set: %v", err)
		return
	}

	fmt.Printf("Recovering %d items from processing set\n", len(items))

	// for batch processing
	pipe := redisClient.Pipeline()

	for uniqueKey, itemJSON := range items {
		var item models.RequestLog
		err := json.Unmarshal([]byte(itemJSON), &item)
		if err != nil {
			log.Printf("Error unmarshaling item: %v", err)
			continue
		}
		redisClient.RPush(ctx, config.RedisQueueName, itemJSON)

		// Use a pipeline to perform batch processing
		pipe.HDel(ctx, config.RedisProcessingQueueName, uniqueKey)
	}

	// remove all the items from HashSet
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Printf("Error deleting from the '%s' queue: %v", config.RedisProcessingQueueName, err)
	}
	items, _ = redisClient.HGetAll(ctx, config.RedisProcessingQueueName).Result()
	fmt.Printf("Remaining items in processing queue after del: %d\n", len(items))
}
