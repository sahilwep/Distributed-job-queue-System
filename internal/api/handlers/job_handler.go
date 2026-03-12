package handlers

import (
	"Distributed_job_queue_system/internal/jobs"
	"Distributed_job_queue_system/internal/queue"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Request structure:
type CreateJobRequest struct {
	Type    string      `json:"type" binding:"required"`
	Payload interface{} `json:"payload" binding:"required"` // payload as interface because it can be email-data, image-job, report-job, etc..
}

// Create job
func CreateJob(c *gin.Context) {
	var req CreateJobRequest

	// bind incoming request to req variable.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Convert the payload JSON
	payloadBytes, err := json.Marshal(req.Payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to serialize payload",
		})
	}

	// Create job struct:
	job := jobs.Job{
		ID:         uuid.New().String(),
		Type:       req.Type,
		Payload:    string(payloadBytes),
		Status:     string(jobs.StatusQueued), // while creating status of job as queued
		Retries:    0,
		MaxRetries: 3,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	// Call repository layer & store the information into DB:
	repo := jobs.NewRepository()

	err = repo.CreateJob(&job)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Redis job-pushing:
	q := queue.NewRedisQueue()

	err = q.Push(job.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "failed to push job ot queue",
		})
		return
	}

	// return job created.
	c.JSON(http.StatusOK, gin.H{
		"job_id":     job.ID,
		"job_status": job.Status,
	})
}

// fetch job
func GetJob(c *gin.Context) {
	jobID := c.Param("id")

	repo := jobs.NewRepository()

	job, err := repo.GetJobByID(jobID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "job not found",
		})
		return
	}

	c.JSON(http.StatusOK, job)
}
