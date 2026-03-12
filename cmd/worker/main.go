package main

import (
	"Distributed_job_queue_system/internal/worker"
	"Distributed_job_queue_system/pkg/database"
	"log"
)

func main() {
	database.ConnectPostgres()

	log.Println("Worker stated")

	w := worker.NewWorker()

	w.Start()
}
