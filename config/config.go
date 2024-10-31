package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

var (
	MongoURI                 string
	MongoDB                  string
	MongoCollection          string
	RedisHost                string
	RedisPort                string
	RedisPassword            string
	RedisQueueName           string
	RedisProcessingQueueName string
	NumberOfWorkers          int
	PrintHeartbeat           bool
	LogPrintingDelay         int
	StartHealthCheckServer   bool
)

func LoadConfig() error {

	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using environment variables.")
		return err
	}
	// Now, assign the environment variables to the config variables
	MongoURI = os.Getenv("MONGODB_URI")
	MongoDB = os.Getenv("MONGODB_DB")
	MongoCollection = os.Getenv("MONGODB_COLLECTION")

	RedisHost = os.Getenv("REDIS_HOST")
	RedisPort = os.Getenv("REDIS_PORT")
	RedisPassword = os.Getenv("REDIS_PASSWORD")
	RedisQueueName = os.Getenv("REDIS_QUEUE_NAME")
	RedisProcessingQueueName = os.Getenv("REDIS_PROCESSING_QUEUE_NAME")

	workerStr := os.Getenv("NUMBER_OF_WORKERS")
	workerInt, err := strconv.Atoi(workerStr)
	if err != nil {
		log.Printf("Error converting NUMBER_OF_WORKERS to int: %v\n", err)
		return err
	}

	NumberOfWorkers = workerInt

	PrintHeartbeatStr := os.Getenv("PRINT_HEARTBEAT")
	PrintHeartbeat, err = strconv.ParseBool(PrintHeartbeatStr)
	if err != nil {
		log.Printf("Error converting PRINT_HEARTBEAT to bool: %v\n", err)
		return err
	}

	LogPrintingDelayStr := os.Getenv("LOG_PRINTING_DELAY")
	LogPrintingDelay, err = strconv.Atoi(LogPrintingDelayStr)
	if err != nil {
		log.Printf("Error converting LOG_PRINTING_DELAY to int: %v\n", err)
		return err
	}

	StartHealthCheckServerStr := os.Getenv("START_HEALTH_CHECK_SERVER")
	StartHealthCheckServer, err = strconv.ParseBool(StartHealthCheckServerStr)
	if err != nil {
		log.Printf("Error converting START_HEALTH_CHECK_SERVER to bool: %v\n", err)
		return err
	}

	return nil
}
