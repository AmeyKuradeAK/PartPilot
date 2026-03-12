package parser

import (
	"strings"
)

// ColumnMapping represents the mapped indices of important parts
type ColumnMapping struct {
	PartNumberIdx int
	QuantityIdx   int
	DescIdx       int
}

// Ensure Valid returns true if we found at least a part number or description
func (m *ColumnMapping) Valid() bool {
	return m.PartNumberIdx != -1 || m.DescIdx != -1
}

func DetectColumns(headers []string) *ColumnMapping {
	mapping := &ColumnMapping{
		PartNumberIdx: -1,
		QuantityIdx:   -1,
		DescIdx:       -1,
	}

	for i, h := range headers {
		clean := strings.ToLower(strings.TrimSpace(h))

		if mapping.PartNumberIdx == -1 {
			if clean == "part number" || clean == "mpn" || clean == "mfr part #" || clean == "component id" || clean == "p/n" || clean == "part" {
				mapping.PartNumberIdx = i
				continue
			}
		}

		if mapping.QuantityIdx == -1 {
			if strings.Contains(clean, "qty") || strings.Contains(clean, "quantity") || strings.Contains(clean, "count") {
				mapping.QuantityIdx = i
				continue
			}
		}

		if mapping.DescIdx == -1 {
			if clean == "description" || clean == "desc" || clean == "component" || clean == "name" {
				mapping.DescIdx = i
				continue
			}
		}
	}

	return mapping
}
