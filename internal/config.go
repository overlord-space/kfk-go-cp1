package internal

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var BootstrapServers string
var TopicName string

func init() {
	log.Print("Initializing config")
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	BootstrapServers = os.Getenv("BOOTSTRAP_SERVERS")
	TopicName = os.Getenv("TOPIC_NAME")
}
