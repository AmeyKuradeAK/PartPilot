package parser

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type ParsedRow struct {
	RowIndex    int
	RawName     string // MPN or Description
	Quantity    int
}

type ParseResult struct {
	Rows         []ParsedRow
	Headers      []string // raw headers for fallback mapping
	MappingFound bool
}

func ParseCSV(reader io.Reader, userMapping *ColumnMapping) (*ParseResult, error) {
	// Read all data to easily handle headers
	b, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("read csv failed: %w", err)
	}

	csvReader := csv.NewReader(bytes.NewReader(b))
	csvReader.TrimLeadingSpace = true
	csvReader.FieldsPerRecord = -1 // Allow variable number of fields

	records, err := csvReader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("parse csv failed: %w", err)
	}

	if len(records) < 2 {
		return nil, fmt.Errorf("csv file is empty or missing data rows")
	}

	headers := records[0]
	mapping := userMapping
	if mapping == nil {
		mapping = DetectColumns(headers)
	}

	result := &ParseResult{
		Headers:      headers,
		MappingFound: mapping.Valid(),
	}

	if !result.MappingFound {
		return result, nil // Stop parsing rows, return headers so API can ask user
	}

	for i := 1; i < len(records); i++ {
		rowStr := records[i]
		
		// Extract values based on mapping safely
		qty := 1 // Default
		if mapping.QuantityIdx != -1 && mapping.QuantityIdx < len(rowStr) {
			qStr := strings.TrimSpace(rowStr[mapping.QuantityIdx])
			if q, err := strconv.Atoi(qStr); err == nil && q > 0 {
				qty = q
			}
		}

		rawName := ""
		// Prefer Part Number column, fallback to Description
		if mapping.PartNumberIdx != -1 && mapping.PartNumberIdx < len(rowStr) {
			rawName = strings.TrimSpace(rowStr[mapping.PartNumberIdx])
		}
		if rawName == "" && mapping.DescIdx != -1 && mapping.DescIdx < len(rowStr) {
			rawName = strings.TrimSpace(rowStr[mapping.DescIdx])
		}

		if rawName != "" {
			result.Rows = append(result.Rows, ParsedRow{
				RowIndex: i,
				RawName:  rawName,
				Quantity: qty,
			})
		}
	}

	return result, nil
}
