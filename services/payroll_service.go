package services

import (
	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
)

type PayrollService struct {
	Calculator *engine.Calculator
}

func (s *PayrollService) Run(basicPay float64, optional []models.Deduction) models.PayrollResult {
	return s.Calculator.Calculate(basicPay, optional)
}
