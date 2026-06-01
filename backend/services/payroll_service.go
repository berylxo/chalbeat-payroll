package services

import (
	"math"

	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
)

type PayrollService struct {
	Calculator *engine.Calculator
}

func (s *PayrollService) Run(basicPay float64, optional []models.Deduction) models.PayrollResult {
	cents := int64(math.Round(basicPay * 100.0))
	return s.Calculator.Calculate(cents, optional)
}
