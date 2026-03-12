package ranker

import (
	"testing"

	"github.com/partpilot/core-engine/internal/supplier"
)

func TestCalculateScore(t *testing.T) {
	reqQty := 10

	goodStock := supplier.Result{UnitPrice: 1.0, StockQty: 100, LeadTimeDays: 0, MOQ: 1}
	noStock := supplier.Result{UnitPrice: 1.0, StockQty: 0, LeadTimeDays: 14, MOQ: 1}
	highMOQ := supplier.Result{UnitPrice: 1.0, StockQty: 100, LeadTimeDays: 0, MOQ: 50}
	expensive := supplier.Result{UnitPrice: 5.0, StockQty: 100, LeadTimeDays: 0, MOQ: 1}

	scoreGood := CalculateScore(goodStock, reqQty)
	scoreNo := CalculateScore(noStock, reqQty)
	scoreMOQ := CalculateScore(highMOQ, reqQty)
	scoreExp := CalculateScore(expensive, reqQty)

	// In-stock should rank higher than out-of-stock
	if scoreGood <= scoreNo {
		t.Errorf("Good stock score (%v) should be > no stock (%v)", scoreGood, scoreNo)
	}

	// Good should rank higher than high MOQ
	if scoreGood <= scoreMOQ {
		t.Errorf("Good stock score (%v) should be > high MOQ (%v)", scoreGood, scoreMOQ)
	}

	// Good should rank higher than expensive
	if scoreGood <= scoreExp {
		t.Errorf("Cheap score (%v) should be > expensive score (%v)", scoreGood, scoreExp)
	}
}

func TestRankResults(t *testing.T) {
	reqQty := 5

	results := []supplier.Result{
		{Supplier: "A", UnitPrice: 1.0, StockQty: 0, LeadTimeDays: 14, MOQ: 1}, // out of stock
		{Supplier: "B", UnitPrice: 10.0, StockQty: 100, LeadTimeDays: 0, MOQ: 1}, // expensive but stock
		{Supplier: "C", UnitPrice: 1.2, StockQty: 100, LeadTimeDays: 0, MOQ: 1}, // perfect
		{Supplier: "D", UnitPrice: 1.0, StockQty: 100, LeadTimeDays: 0, MOQ: 50}, // bad MOQ
	}

	ranked := RankResults(results, reqQty)

	// Expected order: C (perfect), B (expensive but stock), D (bad MOQ but stock), A (no stock)
	if ranked[0].Supplier != "C" {
		t.Errorf("Expected C to be #1, got %s", ranked[0].Supplier)
	}
	if ranked[3].Supplier != "A" {
		t.Errorf("Expected A to be last, got %s", ranked[3].Supplier)
	}
}
