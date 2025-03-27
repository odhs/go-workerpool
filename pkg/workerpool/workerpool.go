// Package workerpool provides a flexible and efficient implementation of a worker pool
// for concurrent task processing. It allows users to define a pool of workers that
// process jobs concurrently and return results.
//
// The package defines the following key components:
// - Job: Represents a task to be processed.
// - Result: Represents the result of a processed task.
// - ProcessFunc: A function type that processes a Job and returns a Result.
// - WorkerPool: An interface defining the worker pool's behavior.
// - Config: Configuration options for the worker pool, including the number of workers
//   and a logger for logging messages.
// - State: Represents the current state of the worker pool (Idle, Running, Stopping).
//
// The worker pool is implemented in the workerPool struct, which manages the lifecycle
// of workers, processes jobs, and handles graceful shutdown.
//
// Key Features:
// - Configurable number of workers.
// - Graceful shutdown using context and stop signals.
// - Logging support for monitoring worker activity.
// - Thread-safe state management.
//
// Usage:
// 1. Create a new worker pool using the New function, providing a ProcessFunc and Config.
// 2. Start the worker pool by calling the Start method with a context and input channel.
// 3. Stop the worker pool gracefully using the Stop method.

package workerpool

// @file workerpool.go
// @title Worker Pool
// @version 1.0.71
// @date 2025-03-27T15:40:21.0348082-03:00

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
)

// Represents a task to be processed.
// type Job any // it could be "any", from the Go 1.18
type Job interface{}

// Represents the result of a processed task.
type Result interface{}

// ProcessFunc defines the function that processes a `Job` and returns a `Result`.
type ProcessFunc func(ctx context.Context, job Job) Result

// WorkerPool defines the interface for the worker pool.
type WorkerPool interface {
	// Starts the Pool of Workers.
	Start(ctx context.Context, inputCh <-chan Job) (<-chan Result, error)
	// Stops the pool of Workers.
	Stop() error
	// Check if the pool is running.
	IsRunning() bool
}

// State represents the current state of the worker pool.
type State int

// States of the worker pool.
// iota faz com que as constantes sejam indexadas em ordem 0,1,2... automaticamente
const (
	StateIdle     State = iota // Idle (Waiting to Start)
	StateRunning               // Running
	StateStopping              // Stopping
)

// Config contains the configuration for the worker pool.
type Config struct {
	WorkerCount int          // Number of workers to create.
	Logger      *slog.Logger // Logger to register messages.
}

// DefaultConfig returns a default configuration for the worker pool.
func DefaultConfig() Config {
	return Config{
		WorkerCount: 1,              // Default to 1 worker. @TODO: runtime.NumCPU()
		Logger:      slog.Default(), // Use the default logger.
	}
}

/// IMPLEMENTATION

// Workerpool orchestra the tasks
type workerPool struct {
	workerCount int            // Number of workers in the pool.
	processFunc ProcessFunc    // Function to process each job.
	logger      *slog.Logger   // Logger for logging messages.
	state       State          // Current state of the worker pool.
	stateMutex  sync.Mutex     // To prevent a Worker from being executed more than once
	stopCh      chan struct{}  // Channel to signal workers to stop.
	stopWg      sync.WaitGroup // WaitGroup to wait for all workers to finish.
}

// New creates a new Worker Pool with the given process function and configuration.
func New(processFunc ProcessFunc, config Config) *workerPool {
	// Ensure at least one worker is created.}
	if config.WorkerCount < 1 {
		config.WorkerCount = 1
	}

	if config.Logger == nil {
		config.Logger = slog.Default()
	}

	return &workerPool{
		processFunc: processFunc,
		workerCount: config.WorkerCount,
		stopCh:      make(chan struct{}),
		state:       StateIdle,
		logger:      config.Logger,
	}
}

// Start initializes the worker pool and starts the workers.
func (wp *workerPool) Start(ctx context.Context, inputCh <-chan Job) (<-chan Result, error) {
	wp.stateMutex.Lock()
	defer wp.stateMutex.Unlock()

	// If Workerpool is not in an idle state, we cannot start	it
	if wp.state != StateIdle {
		return nil, fmt.Errorf("workerpool is not in idle state")
	}
	wp.state = StateRunning

	// Create the results channel to send processed results.
	resultCh := make(chan Result)
	wp.stopCh = make(chan struct{}) // Reset the stop channel.

	// Add the number of workers to the WaitGroup.
	wp.stopWg.Add(wp.workerCount)

	// Start each worker in a separate goroutine.
	for i := 0; i < wp.workerCount; i++ {
		go wp.worker(ctx, i, inputCh, resultCh)
	}

	// Goroutine to close the result channel and reset the state when all workers finish.
	go func() {
		// Wait for all workers to finish
		wp.stopWg.Wait()
		// Close the result channel
		close(resultCh)
		// Set the state to idle
		wp.stateMutex.Lock()
		wp.state = StateIdle
		wp.stateMutex.Unlock()
	}()

	// Return the result channel
	// The caller can read the results from this channel
	return resultCh, nil
}

// Stop signals the worker pool to stop and waits for all workers to finish.
func (wp *workerPool) Stop() error {
	wp.stateMutex.Lock()
	defer wp.stateMutex.Unlock()

	// Ensure the pool is in the running state before stopping.
	if wp.state != StateRunning {
		return fmt.Errorf("workerpool is not in running state")
	}

	// Signal workers to stop by closing the stop channel.
	wp.state = StateStopping
	close(wp.stopCh)

	// Wait for all workers to finish.
	wp.stopWg.Wait()

	wp.state = StateStopping
	return nil
}

// IsRunning checks if the worker pool is currently running.
func (wp *workerPool) IsRunning() bool {
	wp.stateMutex.Lock()
	defer wp.stateMutex.Unlock()

	return wp.state == StateRunning
}

// worker is a function that processes the tasks
// and sends the results to the result channel.
func (wp *workerPool) worker(ctx context.Context, id int, inputCh <-chan Job, resultCh chan<- Result) {
	wp.logger.Info("worker started", "worker_id", id)
	defer wp.stopWg.Done()

	for {
		select {
		case <-wp.stopCh:
			// Stop signal received, exit the worker.
			wp.logger.Info("Worker interrupted", "worker_id", id)
			return

		case <-ctx.Done():
			// Context canceled, exit the worker.
			wp.logger.Info("Worker interrupted by context", "worker_id", id)
			return

		case job, ok := <-inputCh:
			if !ok {
				// Input channel closed, exit the worker.
				wp.logger.Info("Worker interrupted by input channel closed", "worker_id", id)
				return
			}

			// Process the job using the provided process function.
			result := wp.processFunc(ctx, job)

			// Send the result to the result channel or handle stop/context signals.
			select {
			// Successfully sent the result.
			case resultCh <- result:
			// Stop signal received, exit the worker.
			case <-wp.stopCh:
				wp.logger.Info("Worker interrupted, not sending ", "worker_id", id)
				return
			// Context canceled, exit the worker.
			case <-ctx.Done():
				wp.logger.Info("Worker interrupted by context, worker interrupted", "worker_id", id)
				return
			}
		}
	}

}
