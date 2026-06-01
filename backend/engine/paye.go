package engine

import "github.com/berylxo/chalbeat-payroll/models"

// CalculatePAYE returns the gross PAYE (before relief) and the net PAYE (after applying personal relief).
func CalculatePAYE(taxableIncome int64, rules models.Rules) (int64, int64) {
	var grossTax int64
	remainingIncome := taxableIncome

	for _, band := range rules.PayeBands {
		if remainingIncome <= 0 {
			break
		}

		// If width is -1, tax the entire remaining balance in this final bracket
		if band.Width == -1 {
			grossTax += calculatePercentage(remainingIncome, band.Rate)
			break
		}

		taxableInThisBand := remainingIncome
		if taxableInThisBand > band.Width {
			taxableInThisBand = band.Width
		}

		grossTax += calculatePercentage(taxableInThisBand, band.Rate)
		remainingIncome -= taxableInThisBand
	}

	netTax := grossTax - rules.PersonalRelief
	if netTax < 0 {
		netTax = 0
	}
	return grossTax, netTax
}
