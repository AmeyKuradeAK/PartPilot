package ai

import (
	"regexp"
	"strings"
)

var (
	// Matches common structured part number patterns:
	// - Mix of uppercase letters and numbers, length > 4 (e.g. STM32F103C8T6)
	// - Words with hyphens and numbers (e.g. MAX232-ESE)
	// This is a heuristic to save API calls and avoid normalizing already clean parts.
	structuredPattern = regexp.MustCompile(`^[A-Z0-9\-\.]{4,}$`)
	
	// Exclude common generic words that might look like acronyms but aren't specific parts
	genericWords = []string{"RESISTOR", "CAPACITOR", "LED", "DIODE", "TRANSISTOR", "MICROCONTROLLER", "MCU", "IC"}
)

// IsStructuredPartNumber uses heuristics to determine if the name is likely
// a specific manufacturer part number (MPN) rather than a descriptive string.
func IsStructuredPartNumber(name string) bool {
	clean := strings.TrimSpace(strings.ToUpper(name))
	
	// If it's just a generic word, it's not a structured MPN
	for _, gw := range genericWords {
		if clean == gw {
			return false
		}
	}
	
	// If it contains spaces, it's likely descriptive, not a single MPN
	if strings.Contains(clean, " ") {
		return false
	}
	
	return structuredPattern.MatchString(clean)
}
