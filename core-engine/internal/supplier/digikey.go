package supplier

import (
	"context"
	"fmt"
	"log"
)

type DigiKey struct {
	clientID     string
	clientSecret string
	isMock       bool
}

func NewDigiKey(clientID, clientSecret string) *DigiKey {
	isMock := false
	if clientID == "" || clientSecret == "" {
		log.Println("DigiKey credentials not provided. Using mock mode.")
		isMock = true
	}
	return &DigiKey{
		clientID:     clientID,
		clientSecret: clientSecret,
		isMock:       isMock,
	}
}

func (d *DigiKey) Name() string {
	return "DigiKey"
}

func (d *DigiKey) Search(ctx context.Context, partNumber string, qty int) ([]Result, error) {
	if d.isMock {
		// Return mock data
		return []Result{
			{
				PartNumber:   partNumber,
				UnitPrice:    0.15,
				StockQty:     10000,
				LeadTimeDays: 0,
				MOQ:          1,
				ProductURL:   fmt.Sprintf("https://www.digikey.com/en/products/result?keywords=%s", partNumber),
			},
		}, nil
	}

	// TODO: Implement actual Client Credentials OAuth2 flow and Keyword Search API
	// 1. POST https://api.digikey.com/v1/oauth2/token (grant_type=client_credentials)
	// 2. GET https://api.digikey.com/Search/v3/Products/Keyword?keywords={partNumber}

	return nil, fmt.Errorf("real digikey api not fully implemented yet")
}
