package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type JobStatus string

const (
	StatusPending              JobStatus = "pending"
	StatusAwaitingConfirmation JobStatus = "awaiting_confirmation"
	StatusProcessing           JobStatus = "processing"
	StatusDone                 JobStatus = "done"
	StatusFailed               JobStatus = "failed"
)

type Job struct {
	ID     string
	BOMID  string
	Status JobStatus
}

// ClaimPendingJob safely claims a single pending job using a CTE + RETURNING
// to avoid race conditions with multiple workers.
func ClaimPendingJob(ctx context.Context) (*Job, error) {
	query := `
		UPDATE jobs
		SET status = $1, claimed_at = NOW()
		WHERE id = (
			SELECT id FROM jobs 
			WHERE status = $2 
			ORDER BY created_at ASC 
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, bom_id, status;
	`

	var job Job
	err := Pool.QueryRow(ctx, query, StatusProcessing, StatusPending).Scan(&job.ID, &job.BOMID, &job.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // No pending jobs
		}
		return nil, fmt.Errorf("error claiming job: %w", err)
	}

	return &job, nil
}

// ClaimAwaitingConfirmationJob claims a job that was paused for human AI confirmation,
// but ONLY if all its AI parts have now been confirmed.
func ClaimAwaitingConfirmationJob(ctx context.Context) (*Job, error) {
	// We look for jobs awaiting confirmation where NO parts are unconfirmed ai parts
	query := `
		UPDATE jobs
		SET status = $1, claimed_at = NOW()
		WHERE id = (
			SELECT j.id FROM jobs j
			WHERE j.status = $2
			AND NOT EXISTS (
				SELECT 1 FROM bom_parts bp 
				WHERE bp.bom_id = j.bom_id 
				AND bp.is_ai_normalized = true 
				AND bp.ai_confirmed = false
			)
			ORDER BY j.created_at ASC
			FOR UPDATE SKIP LOCKED
			LIMIT 1
		)
		RETURNING id, bom_id, status;
	`

	var job Job
	err := Pool.QueryRow(ctx, query, StatusProcessing, StatusAwaitingConfirmation).Scan(&job.ID, &job.BOMID, &job.Status)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil // None ready to resume
		}
		return nil, fmt.Errorf("error claiming awaiting confirmation job: %w", err)
	}

	return &job, nil
}

func UpdateJobStatus(ctx context.Context, jobID string, status JobStatus, errStr *string) error {
	var err error
	if status == StatusDone || status == StatusFailed {
		_, err = Pool.Exec(ctx, `UPDATE jobs SET status = $1, completed_at = NOW(), error = $2 WHERE id = $3`, status, errStr, jobID)
	} else {
		_, err = Pool.Exec(ctx, `UPDATE jobs SET status = $1, error = $2 WHERE id = $3`, status, errStr, jobID)
	}

	if err != nil {
		return fmt.Errorf("failed to update job status: %w", err)
	}
	return nil
}
