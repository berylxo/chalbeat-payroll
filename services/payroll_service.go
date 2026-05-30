package services

import (
	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
)

type PayrollService struct {
	Calculator *engine.Calculator
}

func (s *PayrollService) Run(emp models.Employee, optional []models.Deduction) models.PayrollResult {
	return s.Calculator.Calculate(emp, optional)
}
