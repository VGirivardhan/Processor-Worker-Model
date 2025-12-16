package main

import "time"

type Task struct {
	ID      string    `json:"id"`
	URL     string    `json:"url"`
	StartAt time.Time `json:"start_at"`
}
