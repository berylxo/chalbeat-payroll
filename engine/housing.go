package engine

import "github.com/berylxo/chalbeat-payroll/models"

func CalculateHousing(gross float64, rules models.Rules) float64 {
	return gross * rules.HousingLevyRate
}
