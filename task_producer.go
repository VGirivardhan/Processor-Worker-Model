package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

type TaskProducer struct {
	workerManager WorkerManagerInterface
	taskCounter   int
	running       bool
}

var urlPool = []string{
	"https://httpbin.org/status/200",
	"https://httpbin.org/status/404",
	"https://httpbin.org/status/500",
	"https://httpbin.org/delay/2",
	"https://httpbin.org/delay/1",
	"https://jsonplaceholder.typicode.com/posts/1",
	"https://jsonplaceholder.typicode.com/users/1",
	"https://www.google.com",
	"https://httpbin.org/status/201",
	"https://httpbin.org/status/301",
	"https://httpbin.org/timeout",
	"https://httpbin.org/get",
}

func NewTaskProducer(workerManager WorkerManagerInterface) *TaskProducer {
	return &TaskProducer{
		workerManager: workerManager,
		taskCounter:   0,
		running:       false,
	}
}

func (tp *TaskProducer) Start() {
	log.Println("Starting TaskProducer")
	tp.running = true
	tp.workerManager.Start()

	// Start generating random tasks
	go tp.generateRandomTasks()
}

func (tp *TaskProducer) generateRandomTasks() {
	ticker := time.NewTicker(1 * time.Second) // Generate a task every 2 seconds
	defer ticker.Stop()

	for tp.running {
		<-ticker.C
		tp.generateRandomTask()
	}
}

func (tp *TaskProducer) generateRandomTask() {
	tp.taskCounter++

	// Pick a random URL from the pool
	randomURL := urlPool[rand.Intn(len(urlPool))]

	task := Task{
		ID:      fmt.Sprintf("task-%d", tp.taskCounter),
		URL:     randomURL,
		StartAt: time.Now(),
	}

	log.Printf("Generated random task: %s for URL: %s", task.ID, task.URL)
	tp.workerManager.AddTask(task)
}

func (tp *TaskProducer) ProcessURL(url string) {
	tp.taskCounter++

	task := Task{
		ID:      fmt.Sprintf("task-%d", tp.taskCounter),
		URL:     url,
		StartAt: time.Now(),
	}

	log.Printf("Processing new task: %s for URL: %s", task.ID, task.URL)
	tp.workerManager.AddTask(task)
}

func (tp *TaskProducer) Stop() {
	log.Println("Stopping TaskProducer")
	tp.running = false
	tp.workerManager.Stop()
}
