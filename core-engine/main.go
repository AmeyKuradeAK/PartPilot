package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/partpilot/core-engine/internal/ai"
	"github.com/partpilot/core-engine/internal/config"
	"github.com/partpilot/core-engine/internal/db"
	"github.com/partpilot/core-engine/internal/orchestrator"
	"github.com/partpilot/core-engine/internal/supplier"
)

func main() {
	cfg := config.Load()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Connect to DB
	db.Connect(ctx, cfg.DatabaseURL)
	defer db.Close()

	// Initialize dependencies
	norm := ai.NewNormalizer(cfg.OpenAIAPIKey)
	suppliers := []supplier.Supplier{
		supplier.NewDigiKey(cfg.DigiKeyClientID, cfg.DigiKeyClientSecret),
		supplier.NewMouser(cfg.MouserSearchAPIKey),
	}

	orch := orchestrator.New(suppliers, norm)

	// Graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down gracefully...")
		cancel()
	}()

	log.Printf("PartPilot Core Engine started. Polling every %dms...", cfg.PollIntervalMs)
	ticker := time.NewTicker(time.Duration(cfg.PollIntervalMs) * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// 1. Check for newly pending jobs
			job, err := db.ClaimPendingJob(ctx)
			if err != nil {
				log.Printf("Error claiming pending job: %v", err)
				continue
			}

			if job != nil {
				// Process in foreground for simplicity, or launch goroutine if scaling
				if err := orch.ProcessJob(ctx, job); err != nil {
					log.Printf("Failed processing job %s: %v", job.ID, err)
					errStr := err.Error()
					_ = db.UpdateJobStatus(ctx, job.ID, db.StatusFailed, &errStr)
				}
				continue // If we found a job, look for another immediately without wait
			}

			// 2. Check for awaiting_confirmation jobs that are now fully confirmed
			resumeJob, err := db.ClaimAwaitingConfirmationJob(ctx)
			if err != nil {
				log.Printf("Error claiming awaiting_confirmation job: %v", err)
				continue
			}

			if resumeJob != nil {
				// The API layer updated `ai_confirmed` for all parts, so we can run it again.
				// The orchestrator resumes where it left off because already normalized parts are skipped.
				if err := orch.ProcessJob(ctx, resumeJob); err != nil {
					log.Printf("Failed processing resumed job %s: %v", resumeJob.ID, err)
					errStr := err.Error()
					_ = db.UpdateJobStatus(ctx, resumeJob.ID, db.StatusFailed, &errStr)
				}
			}
		}
	}
}
