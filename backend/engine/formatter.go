package engine

import (
	"fmt"
	"math"
)

// FormatCentsToKES converts integer cents back to human readable currency format string
func FormatCentsToKES(cents int64) string {
	shillings := cents / 100
	remainingCents := math.Abs(float64(cents % 100))
	return fmt.Sprintf("KES %d.%02d", shillings, int(remainingCents))
}
