package parser

import (
	"strings"
	"testing"
)

func TestParseCSV(t *testing.T) {
	csvData := `Item,Part Number,Description,Qty
1,STM32F103C8T6,,5
2,,10k 0603 resistor,10
3,MAX232-ESE,RS232 Driver,
`

	reader := strings.NewReader(csvData)
	res, err := ParseCSV(reader, nil) // nil mapping triggers auto-detect

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !res.MappingFound {
		t.Fatal("expected mapping to be found via auto-detect")
	}

	if len(res.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(res.Rows))
	}

	// Row 1: Part Number only
	if res.Rows[0].RawName != "STM32F103C8T6" {
		t.Errorf("expected STM32F103C8T6, got %s", res.Rows[0].RawName)
	}
	if res.Rows[0].Quantity != 5 {
		t.Errorf("expected qty 5, got %d", res.Rows[0].Quantity)
	}

	// Row 2: Description fallback
	if res.Rows[1].RawName != "10k 0603 resistor" {
		t.Errorf("expected '10k 0603 resistor', got '%s'", res.Rows[1].RawName)
	}
	if res.Rows[1].Quantity != 10 {
		t.Errorf("expected qty 10, got %d", res.Rows[1].Quantity)
	}

	// Row 3: Missing quantity defaults to 1
	if res.Rows[2].Quantity != 1 {
		t.Errorf("expected default qty 1, got %d", res.Rows[2].Quantity)
	}
}

func TestParseCSV_Fallback(t *testing.T) {
	csvData := `A,B,C,D
1,STM32,Desc,5
`
	reader := strings.NewReader(csvData)
	res, err := ParseCSV(reader, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if res.MappingFound {
		t.Fatal("expected mapping NOT to be found")
	}

	if len(res.Headers) != 4 {
		t.Fatalf("expected 4 headers returned for fallback mapping")
	}
}
