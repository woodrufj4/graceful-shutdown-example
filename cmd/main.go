package main

import (
	"log"
	"time"
)

func main() {

	logger := log.Default()

	messages := []string{
		"123",
		"456",
		"789",
		"101112",
		"131415",
	}

	for _, msg := range messages {
		logger.Printf("processing: %s\n", msg)

		time.Sleep(2 * time.Second)

		logger.Printf("processed: %s\n", msg)
	}

	logger.Println("graceful down complete")
}
