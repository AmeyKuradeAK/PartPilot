package parser

import (
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/xuri/excelize/v2"
)

func ParseExcel(reader io.Reader, userMapping *ColumnMapping) (*ParseResult, error) {
	f, err := excelize.OpenReader(reader)
	if err != nil {
		return nil, fmt.Errorf("excel open failed: %w", err)
	}
	defer f.Close()

	// Get first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in excel file")
	}

	rows, err := f.GetRows(sheets[0])
	if err != nil {
		return nil, fmt.Errorf("failed getting rows from sheet: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("excel file is empty or missing data rows")
	}

	headers := rows[0]
	mapping := userMapping
	if mapping == nil {
		mapping = DetectColumns(headers)
	}

	result := &ParseResult{
		Headers:      headers,
		MappingFound: mapping.Valid(),
	}

	if !result.MappingFound {
		return result, nil
	}

	for i := 1; i < len(rows); i++ {
		rowStr := rows[i]

		qty := 1
		if mapping.QuantityIdx != -1 && mapping.QuantityIdx < len(rowStr) {
			qStr := strings.TrimSpace(rowStr[mapping.QuantityIdx])
			if q, err := strconv.Atoi(qStr); err == nil && q > 0 {
				qty = q
			}
		}

		rawName := ""
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
