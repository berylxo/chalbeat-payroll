package models

type Employee struct {
	EmployeeID  string  `json:"employee_id"`
	Name        string  `json:"name"`
	PhoneNumber string  `json:"phone_number"`
	Email       string  `json:"email"`
	NationalID  string  `json:"national_id"`
	KraPin      string  `json:"kra_pin"`
	Position    string  `json:"position"`
	BasicPay    float64 `json:"basic_pay"`
}

type Deduction struct {
	Code        string `json:"code"`
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
	Mandatory   bool   `json:"mandatory"`
}

type PayrollResult struct {
	EmployeeID string
	GrossPay   int64
	NetPay     int64
	Deductions []Deduction
}

type PayrollRequest struct {
	EmployeeID         string      `json:"employee_id"`
	BasicPay           float64     `json:"basic_pay"`
	OptionalDeductions []Deduction `json:"optional_deductions"`
}
