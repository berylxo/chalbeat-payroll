package engine

import (
	"math"

	"github.com/berylxo/chalbeat-payroll/models"
)

// Helper safely handles rounding float products back into precise integer cents
func calculatePercentage(amount int64, rate float64) int64 {
	return int64(math.Round(float64(amount) * rate))
}

type Calculator struct {
	Rules models.Rules
}

func (c *Calculator) Calculate(basicPayCents int64, optional []models.Deduction) models.PayrollResult {
	gross := basicPayCents

	// 1. Calculate Statutory Deductions
	nssf := CalculateNSSF(gross, c.Rules)
	shif := CalculateSHIF(gross, c.Rules)
	housing := CalculateHousing(gross, c.Rules)

	// 2. KRA compliance: Deduct allowable elements to find Taxable Salary
	taxableIncome := gross - nssf - shif - housing
	if taxableIncome < 0 {
		taxableIncome = 0
	}

	// 3. Compute income tax on the correct base
	grossPaye, paye := CalculatePAYE(taxableIncome, c.Rules)

	deductions := []models.Deduction{
		// Statutory Pension & Levies (Deducted first)
		{Code: "NSSF", Description: "Social Security (NSSF)", Amount: nssf, Mandatory: true},
		{Code: "SHIF", Description: "Health Insurance (SHIF)", Amount: shif, Mandatory: true},
		{Code: "HOUSING", Description: "Housing Levy", Amount: housing, Mandatory: true},

		{Code: "G-TAX", Description: "Gross PAYE (Before Relief)", Amount: grossPaye, Mandatory: true},
		{Code: "RELIEF", Description: "Less: Personal Tax Relief", Amount: -c.Rules.PersonalRelief, Mandatory: true},

		{Code: "PAYE", Description: "Net Income Tax (PAYE)", Amount: paye, Mandatory: true},
	}

	totalDeductions := paye + nssf + shif + housing
	for _, d := range optional {
		deductions = append(deductions, d)
		totalDeductions += d.Amount
	}

	return models.PayrollResult{
		GrossPay:   gross,
		NetPay:     gross - totalDeductions,
		Deductions: deductions,
	}
}
