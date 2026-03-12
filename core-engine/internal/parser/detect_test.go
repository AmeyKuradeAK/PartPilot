package parser

import (
	"testing"
)

func TestDetectColumns(t *testing.T) {
	tests := []struct {
		name     string
		headers  []string
		expected ColumnMapping
	}{
		{
			name:    "Standard headers",
			headers: []string{"Item", "Part Number", "Description", "Qty"},
			expected: ColumnMapping{
				PartNumberIdx: 1,
				DescIdx:       2,
				QuantityIdx:   3,
			},
		},
		{
			name:    "Alternative headers",
			headers: []string{"MPN", "Component", "Quantity Needed", "Notes"},
			expected: ColumnMapping{
				PartNumberIdx: 0,
				DescIdx:       1,
				QuantityIdx:   2,
			},
		},
		{
			name:    "No matching headers",
			headers: []string{"Foo", "Bar", "Baz"},
			expected: ColumnMapping{
				PartNumberIdx: -1,
				DescIdx:       -1,
				QuantityIdx:   -1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectColumns(tt.headers)
			if got.PartNumberIdx != tt.expected.PartNumberIdx {
				t.Errorf("PartNumberIdx got %v, want %v", got.PartNumberIdx, tt.expected.PartNumberIdx)
			}
			if got.DescIdx != tt.expected.DescIdx {
				t.Errorf("DescIdx got %v, want %v", got.DescIdx, tt.expected.DescIdx)
			}
			if got.QuantityIdx != tt.expected.QuantityIdx {
				t.Errorf("QuantityIdx got %v, want %v", got.QuantityIdx, tt.expected.QuantityIdx)
			}
		})
	}
}
