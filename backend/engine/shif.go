package engine

import "github.com/berylxo/chalbeat-payroll/models"

func CalculateSHIF(gross int64, rules models.Rules) int64 {
	return calculatePercentage(gross, rules.ShifRate)
}
