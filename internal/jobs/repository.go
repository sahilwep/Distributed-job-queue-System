package jobs

import (
	"Distributed_job_queue_system/pkg/database"
	"context"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// Creating entry of job in DB
func (c *Repository) CreateJob(job *Job) error {

	query := `
	INSERT INTO jobs (id, type, payload, status, retries, max_retries, created_at, updated_at)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := database.DB.Exec(
		context.Background(),
		query,
		job.ID,
		job.Type,
		job.Payload,
		job.Status,
		job.Retries,
		job.MaxRetries,
		job.CreatedAt,
		job.UpdatedAt,
	)

	return err
}

// Fetching job by ID from DB:
func (c *Repository) GetJobByID(id string) (*Job, error) {
	query := `
		SELECT id, type, payload, status, retries, max_retries, created_at, updated_at
		FROM jobs
		WHERE id=$1
	`

	row := database.DB.QueryRow(context.Background(), query, id)

	var job Job

	err := row.Scan(
		&job.ID,
		&job.Type,
		&job.Payload,
		&job.Status,
		&job.Retries,
		&job.MaxRetries,
		&job.CreatedAt,
		&job.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &job, nil
}

// Update job status as it job being done by consumers.
func (r *Repository) UpdateJobStatus(job *Job) error {

	query := `
		UPDATE jobs
		SET status=$1, retries=$2, updated_at=$3
		WHERE id=$4
	`
	_, err := database.DB.Exec(
		context.Background(),
		query,
		job.Status,
		job.Retries,
		job.UpdatedAt,
		job.ID,
	)

	return err
}
