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

func (h *PayrollHandler) Calculate(w http.ResponseWriter, r *http.Request) {

	var req models.PayrollRequest
	json.NewDecoder(r.Body).Decode(&req)

	result := h.Service.Run(req.BasicPay, req.OptionalDeductions)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
