package routes

import (
	"net/http"

	"github.com/berylxo/chalbeat-payroll/handlers"
)

func SetupRouter() *http.ServeMux {

	mux := http.NewServeMux()

	h := &handlers.PayrollHandler{}

	mux.HandleFunc("/api/v1/payroll/calculate", h.Calculate)

	return mux
}
