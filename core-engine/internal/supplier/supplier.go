package supplier

import (
	"context"
)

type Result struct {
	Supplier     string
	PartNumber   string
	UnitPrice    float64
	StockQty     int
	LeadTimeDays int
	MOQ          int
	ProductURL   string
}

type Supplier interface {
	Search(ctx context.Context, partNumber string, qty int) ([]Result, error)
	Name() string
}
