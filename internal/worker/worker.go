package worker

import (
	"Distributed_job_queue_system/internal/jobs"
	"Distributed_job_queue_system/internal/queue"
	"errors"
	"log"
	"time"
)

type Worker struct {
	queue      *queue.RedisQueue
	repo       *jobs.Repository
	workerPool int
}

func NewWorker() *Worker {
	return &Worker{
		queue:      queue.NewRedisQueue(),
		repo:       jobs.NewRepository(),
		workerPool: 5, // number of worker available
	}
}

func (w *Worker) Start() {

	jobChan := make(chan string)

	for i := 0; i < w.workerPool; i++ {
		go w.WorkerLoop(jobChan)
	}

	for {
		jobID, err := w.queue.Pop()
		if err != nil {
			time.Sleep(1 * time.Second)
			continue
		}

		jobChan <- jobID
	}
}

// WorkerPools:
func (w *Worker) WorkerLoop(jobChan chan string) {

	for jobID := range jobChan {
		log.Println("worker processing job:", jobID)

		w.ProcessJob(jobID)
	}
}

// WOrker tries to process jobs:
func (w *Worker) ProcessJob(jobID string) {

	job, err := w.repo.GetJobByID(jobID)
	if err != nil {
		log.Println("job not found:", err)
		return
	}

	job.Status = string(jobs.StatusProcessing)
	job.UpdatedAt = time.Now()

	err = w.repo.UpdateJobStatus(job)
	if err != nil {
		log.Println("Failed updating job status")
		return
	}

	log.Println("Executing job:", job.ID)

	err = w.executeJob(job)
	if err != nil {
		job.Retries++

		if job.Retries > job.MaxRetries {
			log.Println("Job moved to DLQ:", job.ID)

			job.Status = string(jobs.StatusFailed)
			w.queue.PushDead(job.ID)
		} else {
			log.Println("Retrying job:", job.ID)

			job.Status = string(jobs.StatusQueued)
			w.queue.Push(job.ID)
		}

		job.UpdatedAt = time.Now()
		w.repo.UpdateJobStatus(job)

		return
	}

	job.Status = string(jobs.StatusCompleted)
	job.UpdatedAt = time.Now()

	w.repo.UpdateJobStatus(job)

	log.Println("Job completed:", job.ID)
}

// TYPE Check for job completion/failure
func (w *Worker) executeJob(job *jobs.Job) error {

	switch job.Type {
	case "email":
		// Call for email service.
		log.Println("Sending email job")

	case "image":
		// call for image service.
		log.Println("Processing image job")

	default:
		// This job should fails:
		return errors.New("Job failed!!!")
	}

	time.Sleep(5 * time.Second)

	return nil
}
