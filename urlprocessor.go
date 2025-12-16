package main

import (
	"log"
	"net/http"
)

type URLProcessor struct {
	ID          string
	TaskChan    chan Task
	client      *http.Client
	successChan chan string // Shared success channel (task ID)
	errorChan   chan string // Shared error channel (task ID)
}

func NewURLProcessor(id string, taskChanBuffer int, client *http.Client, successChan, errorChan chan string) *URLProcessor {
	return &URLProcessor{
		ID:          id,
		TaskChan:    make(chan Task, taskChanBuffer),
		client:      client,
		successChan: successChan,
		errorChan:   errorChan,
	}
}

func (up *URLProcessor) Start() {
	log.Printf("URLProcessor %s started", up.ID)

	for task := range up.TaskChan {
		up.processTask(task)
	}
}

func (up *URLProcessor) processTask(task Task) {

	resp, err := up.client.Get(task.URL)

	if err != nil {

		// Send task ID to shared error channel
		select {
		case up.errorChan <- task.ID:
		default:
			log.Printf("Error channel full for task %s", task.ID)
		}
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Send task ID to shared success channel
		select {
		case up.successChan <- task.ID:
		default:
			log.Printf("Success channel full for task %s", task.ID)
		}
	} else {
		// Send task ID to shared error channel
		select {
		case up.errorChan <- task.ID:
		default:
			log.Printf("Error channel full for task %s", task.ID)
		}
	}
}
