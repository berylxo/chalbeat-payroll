package engine

import "github.com/berylxo/chalbeat-payroll/models"

func CalculateSHIF(gross float64, rules models.Rules) float64 {
	return gross * rules.ShifRate
}
