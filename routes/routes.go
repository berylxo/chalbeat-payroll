package routes

import (
	"net/http"

	"github.com/berylxo/chalbeat-payroll/handlers"
	"github.com/berylxo/chalbeat-payroll/services"
)

func SetupRouter(service *services.PayrollService) *http.ServeMux {

	mux := http.NewServeMux()

	handler := &handlers.PayrollHandler{
		Service: service,
	}

	mux.HandleFunc("/api/v1/payroll/calculate", handler.Calculate)

	return mux
}
