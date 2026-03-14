package supplier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Mouser struct {
	apiKey string
	isMock bool
	client *http.Client
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
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

func (m *Mouser) Name() string {
	return "Mouser"
}

func parseMouserPrice(priceStr string) float64 {
	clean := strings.ReplaceAll(priceStr, "$", "")
	clean = strings.ReplaceAll(clean, "£", "")
	clean = strings.ReplaceAll(clean, "€", "")
	clean = strings.ReplaceAll(clean, ",", "")
	val, _ := strconv.ParseFloat(strings.TrimSpace(clean), 64)
	return val
}

func parseMouserStock(availStr string) int {
	clean := strings.Split(availStr, " ")[0]
	clean = strings.ReplaceAll(clean, ",", "")
	val, _ := strconv.Atoi(clean)
	return val
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
				MOQ:          10,
				ProductURL:   fmt.Sprintf("https://www.mouser.com/c/?q=%s", partNumber),
			},
		}, nil
	}

	url := fmt.Sprintf("https://api.mouser.com/api/v2/search/partnumber?apiKey=%s", m.apiKey)
	
	reqBody := map[string]interface{}{
		"SearchByPartRequest": map[string]interface{}{
			"mouserPartNumber": partNumber,
		},
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal mouser request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create mouser request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := m.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("mouser http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("mouser api returned status %d", resp.StatusCode)
	}

	var mRes struct {
		Errors []interface{} `json:"Errors"`
		SearchResults struct {
			NumberOfResult int `json:"NumberOfResult"`
			Parts []struct {
				ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
				ProductDetailUrl       string `json:"ProductDetailUrl"`
				Availability           string `json:"Availability"`
				PriceBreaks []struct {
					Quantity int    `json:"Quantity"`
					Price    string `json:"Price"`
				} `json:"PriceBreaks"`
			} `json:"Parts"`
		} `json:"SearchResults"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&mRes); err != nil {
		return nil, fmt.Errorf("failed to decode mouser response: %w", err)
	}

	if len(mRes.Errors) > 0 {
		return nil, fmt.Errorf("mouser api returned internal errors: %v", mRes.Errors)
	}

	var results []Result
	for _, p := range mRes.SearchResults.Parts {
		stock := parseMouserStock(p.Availability)
		
		// Find best price break for requested qty
		bestPrice := 0.0
		moq := 1
		for _, pb := range p.PriceBreaks {
			price := parseMouserPrice(pb.Price)
			if pb.Quantity <= qty && price > 0 {
				bestPrice = price
				moq = pb.Quantity
			}
		}

		if bestPrice == 0 && len(p.PriceBreaks) > 0 {
			bestPrice = parseMouserPrice(p.PriceBreaks[0].Price)
			moq = p.PriceBreaks[0].Quantity
		}

		results = append(results, Result{
			PartNumber:   p.ManufacturerPartNumber,
			UnitPrice:    bestPrice,
			StockQty:     stock,
			LeadTimeDays: 0,
			MOQ:          moq,
			ProductURL:   p.ProductDetailUrl,
		})
	}

	return results, nil
}
