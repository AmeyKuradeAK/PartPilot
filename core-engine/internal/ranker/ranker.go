package ranker

import (
	"sort"

	"github.com/partpilot/core-engine/internal/supplier"
)

// RankResults takes a slice of supplier results for a SPECIFIC part
// and the required quantity, and returns the same slice sorted from best to worst.
func RankResults(results []supplier.Result, requiredQty int) []supplier.Result {
	if len(results) <= 1 {
		return results
	}

	// Create a copy to sort
	ranked := make([]supplier.Result, len(results))
	copy(ranked, results)

	sort.SliceStable(ranked, func(i, j int) bool {
		rA := ranked[i]
		rB := ranked[j]
		return CalculateScore(rA, requiredQty) > CalculateScore(rB, requiredQty)
	})

	return ranked
}

// CalculateScore calculates a composite score for a result. Higher is better.
func CalculateScore(r supplier.Result, requiredQty int) float64 {
	score := 1000.0 // Base score

	// 1. Total Cost Penalty (lower cost = higher score)
	totalCost := r.UnitPrice * float64(requiredQty)
	// Subtract a scaled cost. The scaling factor depends on expected price ranges, 
	// but a simple linear penalty works for relative comparison within the exact same part.
	score -= totalCost

	// 2. Stock Factor
	if r.StockQty >= requiredQty {
		score += 500.0 // Full boost
	} else if r.StockQty > 0 {
		// Partial stock boost
		score += 100.0 * (float64(r.StockQty) / float64(requiredQty))
	} else {
		score -= 500.0 // Out of stock penalty
	}

	// 3. MOQ Factor
	if r.MOQ > requiredQty {
		// e.g. Need 5, MOQ 100 -> penalize heavily
		score -= 300.0
	}

	// 4. Lead Time Factor
	if r.LeadTimeDays > 0 {
		score -= float64(r.LeadTimeDays) * 10.0
	}

	return score
}
