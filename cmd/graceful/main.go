package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	logger := log.Default()

	// Channel to signal that a shutdown request has been made.
	mainShutdownCh := make(chan os.Signal, 2)

	// Channel to tell the job to shutdown.
	jobShutdownChan := make(chan os.Signal)

	// Channel for the job to signal that it has completed work.
	jobCompleteChan := make(chan error)

	// Notify the shutdown channel of any request for termination.
	signal.Notify(mainShutdownCh, os.Interrupt, syscall.SIGTERM)

	// Do the work
	go doWork(logger, jobCompleteChan, jobShutdownChan)

	// Wait for the job to complete or for a shutdown signal...

	select {
	case signal := <-mainShutdownCh:
		println()
		logger.Println("recieved shutdown command")
		logger.Println("gracefully shutting down...")

		// Let the job know to shutdown
		go func(shutdownSignal os.Signal) {
			jobShutdownChan <- shutdownSignal
		}(signal)

	case jobErr := <-jobCompleteChan:

		// The job completed normally

		if jobErr != nil {
			logger.Printf("job completed with error: %s\n", jobErr.Error())
			os.Exit(1)
		}

		logger.Println("completed processing")
		os.Exit(0)

	}

	// Gracefully shutdown and wait for the job to complete

	select {
	case <-mainShutdownCh:
		println()
		logger.Println("recieved second shutdown... shutting down immediately")
		logger.Println("data loss may occur")
		os.Exit(1)

	case <-jobCompleteChan:
		// Job cleanly completed!
		logger.Println("completed processing")
	}

	// Graceful shutdown cleanup
	logger.Println("cleaning up residual artifacts...")
	time.Sleep(1 * time.Second)
	logger.Println("graceful down complete")
	os.Exit(0)

}

func doWork(logger *log.Logger, completeChan chan error, shutdownCh <-chan os.Signal) {

	messages := []string{
		"123",
		"456",
		"789",
		"101112",
		"131415",
	}

	for _, msg := range messages {

		select {

		// Should I stop?
		case <-shutdownCh:

			// Yes
			completeChan <- nil
			return
		default:
			// Keep going!
		}

		logger.Printf("processing: %s\n", msg)

		time.Sleep(2 * time.Second)

		logger.Printf("processed: %s\n", msg)

	}

	completeChan <- nil

}
