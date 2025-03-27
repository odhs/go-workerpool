package main

// @title Worker Pool Example
// @version 1.0.13
// @date 2025-03-27T17:12:41.246373-03:00
// @description Example of implementation of a Worker Pool

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/odhs/go-workerpool/pkg/workerpool"
)

type NumberJob struct {
	Number int
}

type ResultNumber struct {
	Value     int
	WorkerID  int
	Timestamp time.Time
}

const numberOfWorkers = 3

// processNumbers é um job de processamento, ele processa números
func processNumbers(ctx context.Context, job workerpool.Job) workerpool.Result {

	number := job.(NumberJob).Number

	// Distributes tasks also among the number of works defined
	workderID := number % numberOfWorkers

	sleepTime := time.Duration(800+rand.Intn(400)) * time.Millisecond
	time.Sleep(sleepTime)

	return ResultNumber{
		Value:     number,
		WorkerID:  workderID,
		Timestamp: time.Now(),
	}

}

func main() {
	maxValue := 20
	bufferSize := 10

	pool := workerpool.New(processNumbers, workerpool.Config{WorkerCount: numberOfWorkers})

	inputCh := make(chan workerpool.Job, bufferSize)
	ctx := context.Background()

	resultCh, err := pool.Start(ctx, inputCh)
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	wg.Add(maxValue)

	fmt.Println("Starting the worker pool with", maxValue, "numbers to be processed")

	go func() {
		for i := 0; i < maxValue; i++ {
			inputCh <- NumberJob{Number: i}
		}
		close(inputCh)
	}()

	go func() {
		for result := range resultCh {
			r := result.(ResultNumber)
			fmt.Printf("Number: %d, WorkerID: %d, Timestamp: %s\n", r.Value, r.WorkerID, r.Timestamp.Format("15:04:05.000"))
			wg.Done()
		}
	}()

	wg.Wait()
	fmt.Printf("All jobs are done\nEvery the %d numbers were processed.\n", maxValue)
}
