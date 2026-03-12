package db

import (
	"context"
	"fmt"
)

type BOMPart struct {
	ID             string
	BOMID          string
	RowIndex       int
	RawName        string
	NormalizedName *string
	Quantity       int
	IsAINormalized bool
	AIConfirmed    bool
}

type SupplierResult struct {
	BOMPartID    string
	JobID        string
	Supplier     string
	PartNumber   string
	UnitPrice    float64
	StockQty     int
	LeadTimeDays int
	MOQ          int
	ProductURL   string
	Rank         int
}

func GetBOMPartsByBOMID(ctx context.Context, bomID string) ([]BOMPart, error) {
	query := `
		SELECT id, bom_id, row_index, raw_name, normalized_name, quantity, is_ai_normalized, ai_confirmed
		FROM bom_parts
		WHERE bom_id = $1
		ORDER BY row_index ASC
	`
	rows, err := Pool.Query(ctx, query, bomID)
	if err != nil {
		return nil, fmt.Errorf("query bom_parts failed: %w", err)
	}
	defer rows.Close()

	var parts []BOMPart
	for rows.Next() {
		var p BOMPart
		if err := rows.Scan(&p.ID, &p.BOMID, &p.RowIndex, &p.RawName, &p.NormalizedName, &p.Quantity, &p.IsAINormalized, &p.AIConfirmed); err != nil {
			return nil, fmt.Errorf("scan bom_part failed: %w", err)
		}
		parts = append(parts, p)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows iteration error: %w", rows.Err())
	}

	return parts, nil
}

func GetBOMFilename(ctx context.Context, bomID string) (string, error) {
	var filename string
	err := Pool.QueryRow(ctx, "SELECT filename FROM boms WHERE id = $1", bomID).Scan(&filename)
	if err != nil {
		return "", fmt.Errorf("failed fetching bom filename: %w", err)
	}
	return filename, nil
}

func InsertBOMParts(ctx context.Context, bomID string, parts []BOMPart) error {
	// Bulk insert
	for _, p := range parts {
		_, err := Pool.Exec(ctx, `
			INSERT INTO bom_parts (bom_id, row_index, raw_name, quantity) 
			VALUES ($1, $2, $3, $4)
		`, bomID, p.RowIndex, p.RawName, p.Quantity)
		if err != nil {
			return fmt.Errorf("insert part error: %w", err)
		}
	}
	return nil
}

func UpdateBOMPartNormalization(ctx context.Context, partID string, normalizedName string, isAI bool) error {
	_, err := Pool.Exec(ctx, `
		UPDATE bom_parts 
		SET normalized_name = $1, is_ai_normalized = $2 
		WHERE id = $3
	`, normalizedName, isAI, partID)
	
	if err != nil {
		return fmt.Errorf("failed to update bom_part normalization: %w", err)
	}
	return nil
}

func SaveSupplierResult(ctx context.Context, res SupplierResult) error {
	query := `
		INSERT INTO supplier_results 
		(bom_part_id, job_id, supplier, part_number, unit_price, stock_qty, lead_time_days, moq, product_url)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := Pool.Exec(ctx, query,
		res.BOMPartID, res.JobID, res.Supplier, res.PartNumber,
		res.UnitPrice, res.StockQty, res.LeadTimeDays, res.MOQ, res.ProductURL,
	)
	
	if err != nil {
		return fmt.Errorf("failed to insert supplier result: %w", err)
	}
	return nil
}

func UpdateSupplierResultRanks(ctx context.Context, partID string, ranks map[string]int) error {
	// Note: ranks map contains mapping from supplier name to rank for this part
	// Typically we'd do a batch update here
	for supplier, rank := range ranks {
		_, err := Pool.Exec(ctx, `
			UPDATE supplier_results 
			SET rank = $1 
			WHERE bom_part_id = $2 AND supplier = $3
		`, rank, partID, supplier)
		if err != nil {
			return fmt.Errorf("failed to update rank for supplier %s: %w", supplier, err)
		}
	}
	return nil
}
