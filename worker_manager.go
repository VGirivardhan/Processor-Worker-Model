package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type WorkerManager struct {
	workers     []*URLProcessor
	taskQueue   chan Task
	wg          sync.WaitGroup
	config      WorkerManagerConfig
	httpClient  *http.Client
	successChan chan string // Shared success channel (task ID)
	errorChan   chan string // Shared error channel (task ID)
}

// WorkerManagerConfig holds configuration for WorkerManager
type WorkerManagerConfig struct {
	NumWorkers       int
	TaskQueueBuffer  int
	WorkerNamePrefix string
}

func NewWorkerManager(config WorkerManagerConfig, httpClient *http.Client) *WorkerManager {
	wm := &WorkerManager{
		workers:     make([]*URLProcessor, config.NumWorkers),
		taskQueue:   make(chan Task, config.TaskQueueBuffer),
		config:      config,
		httpClient:  httpClient,
		successChan: make(chan string, 100), // Shared success channel
		errorChan:   make(chan string, 100), // Shared error channel
	}

	// Create workers directly
	for i := 0; i < config.NumWorkers; i++ {
		workerID := fmt.Sprintf("%s-%d", config.WorkerNamePrefix, i+1)
		worker := NewURLProcessor(workerID, 10, httpClient, wm.successChan, wm.errorChan)
		wm.workers[i] = worker
	}

	return wm
}

func (wm *WorkerManager) Start() {
	log.Printf("Starting WorkerManager with %d workers", len(wm.workers))

	// Start all workers
	for _, worker := range wm.workers {
		wm.wg.Add(1)
		go func(w *URLProcessor) {
			defer wm.wg.Done()
			w.Start()
		}(worker)
	}

	// Start task distributor
	go wm.distributeTasks()

	// Start result handlers
	go wm.handleResults()
}

func (wm *WorkerManager) distributeTasks() {
	workerIndex := 0

	for task := range wm.taskQueue {
		// Round-robin distribution to workers
		selectedWorker := wm.workers[workerIndex]
		workerIndex = (workerIndex + 1) % len(wm.workers)

		// Send task to selected worker
		select {
		case selectedWorker.TaskChan <- task:
			log.Printf("Task %s assigned to %s", task.ID, selectedWorker.ID)
		default:
			log.Printf("Worker %s is busy, task %s queued", selectedWorker.ID, task.ID)
			// If worker is busy, still send but it will block until worker is ready
			selectedWorker.TaskChan <- task
		}
	}
}

func (wm *WorkerManager) handleResults() {
	for {
		select {
		case taskID := <-wm.successChan:
			log.Printf("Task %s completed successfully", taskID)
		case taskID := <-wm.errorChan:
			log.Printf("Task %s completed with failure", taskID)
		}
	}
}

func (wm *WorkerManager) AddTask(task Task) {
	wm.taskQueue <- task
}

func (wm *WorkerManager) Stop() {
	close(wm.taskQueue)

	// Close all worker channels
	for _, worker := range wm.workers {
		close(worker.TaskChan)
	}

	wm.wg.Wait()
	log.Println("WorkerManager stopped")
}
