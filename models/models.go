package models

type Employee struct {
	ID       string
	Name     string
	KraPin   string
	BasicPay float64
}

type Deduction struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Mandatory   bool    `json:"mandatory"`
}

type PayrollResult struct {
	EmployeeID string

	GrossPay float64
	NetPay   float64

	Deductions []Deduction
}
