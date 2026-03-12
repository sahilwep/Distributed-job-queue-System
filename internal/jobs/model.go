package jobs

import "time"

type JobStatus string

const (
	StatusQueued     JobStatus = "queued"
	StatusProcessing JobStatus = "processing"
	StatusCompleted  JobStatus = "completed"
	StatusFailed     JobStatus = "failed"
)

type Job struct {
	ID         string    `json:"id"`
	Type       string    `json:"type"`
	Payload    string    `json:"payload"`
	Status     string    `json:"status"`
	Retries    int       `json:"retries"`
	MaxRetries int       `json:"max_retries"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
