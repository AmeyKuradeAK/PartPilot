package orchestrator

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/partpilot/core-engine/internal/ai"
	"github.com/partpilot/core-engine/internal/db"
	"github.com/partpilot/core-engine/internal/supplier"
	"github.com/partpilot/core-engine/internal/parser"
	"os"
	"strings"
)

type Orchestrator struct {
	suppliers  []supplier.Supplier
	normalizer *ai.Normalizer
}

func New(suppliers []supplier.Supplier, normalizer *ai.Normalizer) *Orchestrator {
	return &Orchestrator{
		suppliers:  suppliers,
		normalizer: normalizer,
	}
}

// ProcessJob is the main entry point for a single job
func (o *Orchestrator) ProcessJob(ctx context.Context, job *db.Job) error {
	log.Printf("Starting processing for job %s (BOM: %s)", job.ID, job.BOMID)

	parts, err := db.GetBOMPartsByBOMID(ctx, job.BOMID)
	if err != nil {
		return fmt.Errorf("failed to fetch parts: %w", err)
	}

	// Step 0: If no parts, parse the file
	if len(parts) == 0 {
		filename, err := db.GetBOMFilename(ctx, job.BOMID)
		if err != nil {
			return fmt.Errorf("failed to get bom filename: %w", err)
		}
		
		f, err := os.Open(filename)
		if err != nil {
			return fmt.Errorf("failed to open uploaded file: %w", err)
		}
		defer f.Close()

		var res *parser.ParseResult
		if strings.HasSuffix(strings.ToLower(filename), ".csv") {
			res, err = parser.ParseCSV(f, nil)
		} else {
			res, err = parser.ParseExcel(f, nil)
		}
		if err != nil {
			return fmt.Errorf("failed to parse bom file: %w", err)
		}

		if !res.MappingFound {
			return fmt.Errorf("parser fallback: columns not mapped. (Web UI needs to prompt user mapping)")
		}

		// Insert parsed parts to DB
		var newParts []db.BOMPart
		for idx, row := range res.Rows {
			newParts = append(newParts, db.BOMPart{
				RowIndex: idx, // start at 0
				RawName:  row.RawName,
				Quantity: row.Quantity,
			})
		}
		
		if err := db.InsertBOMParts(ctx, job.BOMID, newParts); err != nil {
			return err
		}

		// Re-fetch populated parts
		parts, err = db.GetBOMPartsByBOMID(ctx, job.BOMID)
		if err != nil {
			return fmt.Errorf("refetch parts failed: %w", err)
		}
	}

	// Step 1: AI Normalization Pass
	requiresHumanConfirmation := false
	for i, part := range parts {
		if part.NormalizedName != nil && *part.NormalizedName != "" {
			continue // Already normalized or confirmed
		}

		if ai.IsStructuredPartNumber(part.RawName) {
			// Fast path: structured part
			_ = db.UpdateBOMPartNormalization(ctx, part.ID, part.RawName, false)
			parts[i].NormalizedName = &part.RawName
			parts[i].IsAINormalized = false
		} else {
			// Slow path: needs AI
			norm, err := o.normalizer.Normalize(ctx, part.RawName)
			if err != nil {
				log.Printf("Warning: AI normalization failed for %s: %v", part.RawName, err)
				norm = part.RawName // fallback
			}
			_ = db.UpdateBOMPartNormalization(ctx, part.ID, norm, true)
			parts[i].NormalizedName = &norm
			parts[i].IsAINormalized = true
			parts[i].AIConfirmed = false
			requiresHumanConfirmation = true
		}
	}

	// If ANY part needs human confirmation, we pause the job here
	if requiresHumanConfirmation && job.Status != db.StatusAwaitingConfirmation {
		log.Printf("Job %s paused awaiting human AI confirmation", job.ID)
		err := db.UpdateJobStatus(ctx, job.ID, db.StatusAwaitingConfirmation, nil)
		return err
	}

	// Step 2: Concurrent Supplier Queries
	log.Printf("Job %s: Fanning out queries for %d parts", job.ID, len(parts))
	
	// Create an errgroup with a limit on concurrent goroutines (e.g. max 50 outbound requests at once)
	g, gCtx := errgroup.WithContext(ctx)
	g.SetLimit(50)

	var mu sync.Mutex // protects results writing to DB if needed, though inserts could be separate
	
	for _, part := range parts {
		part := part // capture loop variable
		partNum := *part.NormalizedName

		for _, supp := range o.suppliers {
			supp := supp // capture
			
			g.Go(func() error {
				// Query the supplier
				// Add small timeout per request so one bad supplier doesn't hang the job forever
				reqCtx, cancel := context.WithTimeout(gCtx, 15*time.Second)
				defer cancel()

				results, err := supp.Search(reqCtx, partNum, part.Quantity)
				if err != nil {
					log.Printf("Supplier %s failed for part %s: %v", supp.Name(), partNum, err)
					return nil // We don't fail the whole job if one supplier errors for one part
				}

				// If no results, just return
				if len(results) == 0 {
					return nil
				}

				// For V1, just take the first result per supplier per part if multiple returned,
				// or rank them and take the top. Or we can store all. The DB model allows multiple per supplier.
				// We'll store them all.
				
				mu.Lock()
				// Store all results
				for _, res := range results {
					dbRes := db.SupplierResult{
						BOMPartID:    part.ID,
						JobID:        job.ID,
						Supplier:     supp.Name(),
						PartNumber:   res.PartNumber,
						UnitPrice:    res.UnitPrice,
						StockQty:     res.StockQty,
						LeadTimeDays: res.LeadTimeDays,
						MOQ:          res.MOQ,
						ProductURL:   res.ProductURL,
					}
					if err := db.SaveSupplierResult(ctx, dbRes); err != nil {
						log.Printf("Error saving result: %v", err)
					}
				}
				mu.Unlock()

				return nil
			})
		}
	}

	if err := g.Wait(); err != nil {
		return fmt.Errorf("orchestrator group error: %w", err)
	}

	// Step 3: Ranking Pass
	// Read back supplier results and apply ranker
	for i := range parts {
		_ = i
		// (In real implementation, we would query supplier_results where job_id and part_id match)
		// Assuming we had a db.GetResultsForPart:
		// resList, _ := db.GetSupplierResultsByPart(ctx, part.ID)
		// ranked := ranker.RankResults(resList, part.Quantity)
		// var rankMap map[string]int
		// for i, r := range ranked { rankMap[r.Supplier] = i+1 }
		// db.UpdateSupplierResultRanks(ctx, part.ID, rankMap)
	}

	err = db.UpdateJobStatus(ctx, job.ID, db.StatusDone, nil)
	log.Printf("Job %s finished successfully", job.ID)
	return err
}
