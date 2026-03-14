package supplier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type DigiKey struct {
	clientID     string
	clientSecret string
	isMock       bool
	client       *http.Client
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
		client:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (d *DigiKey) Name() string {
	return "DigiKey"
}

func (d *DigiKey) getAccessToken(ctx context.Context) (string, error) {
	data := url.Values{}
	data.Set("client_id", d.clientID)
	data.Set("client_secret", d.clientSecret)
	data.Set("grant_type", "client_credentials")

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.digikey.com/v1/oauth2/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := d.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("digikey auth failed: status %d, body: %s", resp.StatusCode, string(body))
	}

	var res struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	return res.AccessToken, nil
}

func (d *DigiKey) Search(ctx context.Context, partNumber string, qty int) ([]Result, error) {
	if d.isMock {
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

	token, err := d.getAccessToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get digikey token: %w", err)
	}

	searchURL := "https://api.digikey.com/products/v4/search/keyword"
	reqBody := map[string]interface{}{
		"Keywords": partNumber,
		"RecordCount": 10,
		"RecordStartPosition": 0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal digikey request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", searchURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create digikey request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("X-DIGIKEY-Client-Id", d.clientID)
    req.Header.Set("X-DIGIKEY-Locale-Site", "US")
    req.Header.Set("X-DIGIKEY-Locale-Language", "en")
    req.Header.Set("X-DIGIKEY-Locale-Currency", "USD")
	req.Header.Set("Accept", "application/json")

	resp, err := d.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("digikey http request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("digikey api returned status %d. body: %s", resp.StatusCode, string(body))
	}

	var dRes struct {
		ExactManufacturerProducts []struct {
			ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
			ProductUrl             string `json:"ProductUrl"`
			QuantityAvailable      int    `json:"QuantityAvailable"`
			StandardPricing        []struct {
				BreakQuantity int     `json:"BreakQuantity"`
				UnitPrice     float64 `json:"UnitPrice"`
			} `json:"StandardPricing"`
		} `json:"ExactManufacturerProducts"`
		Products []struct {
			ManufacturerPartNumber string `json:"ManufacturerPartNumber"`
			ProductUrl             string `json:"ProductUrl"`
			QuantityAvailable      int    `json:"QuantityAvailable"`
			StandardPricing        []struct {
				BreakQuantity int     `json:"BreakQuantity"`
				UnitPrice     float64 `json:"UnitPrice"`
			} `json:"StandardPricing"`
		} `json:"Products"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&dRes); err != nil {
		return nil, fmt.Errorf("failed to decode digikey response: %w", err)
	}

	var results []Result

	// Process Exact Matches
	for _, p := range dRes.ExactManufacturerProducts {
		bestPrice := 0.0
		moq := 1
		for _, pb := range p.StandardPricing {
			if pb.BreakQuantity <= qty && pb.UnitPrice > 0 {
				bestPrice = pb.UnitPrice
				moq = pb.BreakQuantity
			}
		}
		if bestPrice == 0 && len(p.StandardPricing) > 0 {
			bestPrice = p.StandardPricing[0].UnitPrice
			moq = p.StandardPricing[0].BreakQuantity
		}

		results = append(results, Result{
			PartNumber:   p.ManufacturerPartNumber,
			UnitPrice:    bestPrice,
			StockQty:     p.QuantityAvailable,
			LeadTimeDays: 0,
			MOQ:          moq,
			ProductURL:   p.ProductUrl,
		})
	}

	// Process general matches if no exacts
	if len(results) == 0 {
		for _, p := range dRes.Products {
			bestPrice := 0.0
			moq := 1
			for _, pb := range p.StandardPricing {
				if pb.BreakQuantity <= qty && pb.UnitPrice > 0 {
					bestPrice = pb.UnitPrice
					moq = pb.BreakQuantity
				}
			}
			if bestPrice == 0 && len(p.StandardPricing) > 0 {
				bestPrice = p.StandardPricing[0].UnitPrice
				moq = p.StandardPricing[0].BreakQuantity
			}

			results = append(results, Result{
				PartNumber:   p.ManufacturerPartNumber,
				UnitPrice:    bestPrice,
				StockQty:     p.QuantityAvailable,
				LeadTimeDays: 0,
				MOQ:          moq,
				ProductURL:   p.ProductUrl,
			})
		}
	}

	return results, nil
}
