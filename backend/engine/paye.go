package engine

import "github.com/berylxo/chalbeat-payroll/models"

func CalculatePAYE(gross float64, rules models.Rules) float64 {
	var tax float64

	for _, band := range rules.Paye.Bands {
		if gross > band.Min {
			upper := min(gross, band.Max)
			tax += (upper - band.Min) * band.Rate
		}
	}

	tax -= rules.Paye.PersonalRelief

	if tax < 0 {
		return 0
	}

	return tax
}
