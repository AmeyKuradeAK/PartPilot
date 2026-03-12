package po

import (
	"fmt"
	"time"

	"github.com/jung-kurt/gofpdf"
)

type POConfig struct {
	CompanyName string
	LogoPath    string // Optional logo slot
}

type POLineItem struct {
	Supplier   string
	PartNumber string
	Quantity   int
	UnitPrice  float64
	LineTotal  float64
	ProductURL string
}

func GeneratePDF(cfg POConfig, jobID string, items []POLineItem) (string, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Branding slot
	if cfg.LogoPath != "" {
		// e.g. pdf.Image(cfg.LogoPath, 10, 10, 30, 0, false, "", 0, "")
		pdf.SetY(40) // shift down
	}

	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, fmt.Sprintf("PURCHASE ORDER - %s", cfg.CompanyName))
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(40, 10, fmt.Sprintf("Date: %s", time.Now().Format("2006-01-02")))
	pdf.Ln(15)

	// Table Headers
	pdf.SetFont("Arial", "B", 10)
	pdf.CellFormat(30, 10, "Supplier", "1", 0, "C", false, 0, "")
	pdf.CellFormat(50, 10, "Part Number", "1", 0, "C", false, 0, "")
	pdf.CellFormat(20, 10, "Qty", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Unit Price", "1", 0, "C", false, 0, "")
	pdf.CellFormat(30, 10, "Total", "1", 0, "C", false, 0, "")
	pdf.Ln(10)

	// Table Rows
	pdf.SetFont("Arial", "", 10)
	grandTotal := 0.0

	for _, item := range items {
		pdf.CellFormat(30, 10, item.Supplier, "1", 0, "L", false, 0, "")
		pdf.CellFormat(50, 10, item.PartNumber, "1", 0, "L", false, 0, "")
		pdf.CellFormat(20, 10, fmt.Sprintf("%d", item.Quantity), "1", 0, "C", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("$%.4f", item.UnitPrice), "1", 0, "R", false, 0, "")
		pdf.CellFormat(30, 10, fmt.Sprintf("$%.2f", item.LineTotal), "1", 0, "R", false, 0, "")
		
		grandTotal += item.LineTotal
		pdf.Ln(10)
	}

	// Grand Total
	pdf.Ln(5)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(130, 10, "GRAND TOTAL: ", "", 0, "R", false, 0, "")
	pdf.CellFormat(30, 10, fmt.Sprintf("$%.2f", grandTotal), "", 0, "R", false, 0, "")

	// Save
	outputPath := fmt.Sprintf("/tmp/po_%s.pdf", jobID)
	err := pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return "", fmt.Errorf("failed generating po pdf: %w", err)
	}

	return outputPath, nil
}
