package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/services"
)

type PayrollHandler struct {
	Service *services.PayrollService
}

// PayrollResponseDTO converts integer cents into standard floating-point figures for API consumers
type PayrollResponseDTO struct {
	EmployeeID string                 `json:"employee_id"`
	GrossPay   float64                `json:"gross_pay"`
	NetPay     float64                `json:"net_pay"`
	Deductions []DeductionResponseDTO `json:"deductions"`
}

type DeductionResponseDTO struct {
	Code        string  `json:"code"`
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Mandatory   bool    `json:"mandatory"`
}

func (h *PayrollHandler) Calculate(w http.ResponseWriter, r *http.Request) {
	var req models.PayrollRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	result := h.Service.Run(req.BasicPay, req.OptionalDeductions)

	// Map response from internal cents back into clean floats for the API response
	response := PayrollResponseDTO{
		EmployeeID: result.EmployeeID,
		GrossPay:   float64(result.GrossPay) / 100.0,
		NetPay:     float64(result.NetPay) / 100.0,
	}

	for _, d := range result.Deductions {
		response.Deductions = append(response.Deductions, DeductionResponseDTO{
			Code:        d.Code,
			Description: d.Description,
			Amount:      float64(d.Amount) / 100.0,
			Mandatory:   d.Mandatory,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
