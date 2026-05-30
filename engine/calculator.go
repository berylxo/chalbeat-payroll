package engine

import "github.com/berylxo/chalbeat-payroll/models"

type Calculator struct {
	Rules models.Rules
}

func (c *Calculator) Calculate(emp models.Employee, optional []models.Deduction) models.PayrollResult {

	gross := emp.BasicPay

	paye := CalculatePAYE(gross, c.Rules)
	nssf := CalculateNSSF(gross, c.Rules)
	shif := CalculateSHIF(gross, c.Rules)
	housing := CalculateHousing(gross, c.Rules)

	deductions := []models.Deduction{
		{Code: "PAYE", Description: "Income Tax", Amount: paye, Mandatory: true},
		{Code: "NSSF", Description: "Social Security", Amount: nssf, Mandatory: true},
		{Code: "SHIF", Description: "Health Insurance", Amount: shif, Mandatory: true},
		{Code: "HOUSING", Description: "Housing Levy", Amount: housing, Mandatory: true},
	}

	total := paye + nssf + shif + housing

	for _, d := range optional {
		deductions = append(deductions, d)
		total += d.Amount
	}

	return models.PayrollResult{
		EmployeeID: emp.ID,
		GrossPay:   gross,
		NetPay:     gross - total,
		Deductions: deductions,
	}
}
