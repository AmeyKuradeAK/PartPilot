package supplier

import (
	"context"
	"fmt"
	"log"
)

type Mouser struct {
	apiKey string
	isMock bool
}

func NewMouser(apiKey string) *Mouser {
	isMock := false
	if apiKey == "" {
		log.Println("Mouser API key not provided. Using mock mode.")
		isMock = true
	}
	return &Mouser{
		apiKey: apiKey,
		isMock: isMock,
	}
}

func (m *Mouser) Name() string {
	return "Mouser"
}

func (m *Mouser) Search(ctx context.Context, partNumber string, qty int) ([]Result, error) {
	if m.isMock {
		// Return mock data
		return []Result{
			{
				PartNumber:   partNumber,
				UnitPrice:    0.18,
				StockQty:     5000,
				LeadTimeDays: 0,
				MOQ:          10, // slightly higher MOQ to test ranker logic
				ProductURL:   fmt.Sprintf("https://www.mouser.com/c/?q=%s", partNumber),
			},
		}, nil
	}

	// TODO: Implement actual Mouser Search API
	// POST https://api.mouser.com/api/v2/search/partnumber?apiKey={apiKey}

	return nil, fmt.Errorf("real mouser api not fully implemented yet")
}
