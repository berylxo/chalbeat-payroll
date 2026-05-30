package models

type Employee struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	KraPin   string  `json:"kra_pin"`
	Position string  `json:"position"`
	BasicPay float64 `json:"basic_pay"`
}

type Deduction struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Mandatory   bool    `json:"mandatory"`
}

type PayrollRequest struct {
	BasicPay          float64     `json:"basic_pay"`
	OptionalDeductions []Deduction `json:"optional_deductions"`
}

type PayrollResult struct {
	EmployeeID string

	GrossPay float64
	NetPay   float64

	Deductions []Deduction
}
