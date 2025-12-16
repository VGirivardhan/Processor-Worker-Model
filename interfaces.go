package main

// WorkerManagerInterface for dependency injection
type WorkerManagerInterface interface {
	Start()
	AddTask(task Task)
	Stop()
}

// TaskProducerInterface for dependency injection
type TaskProducerInterface interface {
	Start()
	Stop()
	ProcessURL(url string)
}
