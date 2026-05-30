package routes

import (
	"net/http"
	"strings"

	"github.com/berylxo/chalbeat-payroll/handlers"
	"github.com/berylxo/chalbeat-payroll/services"
)

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		h(w, r)
	}
}

func SetupRouter(payroll *services.PayrollService, employees *services.EmployeeService) *http.ServeMux {
	mux := http.NewServeMux()

	payrollHandler := &handlers.PayrollHandler{Service: payroll}
	employeeHandler := &handlers.EmployeeHandler{Service: employees}
	payslipHandler := &handlers.PayslipHandler{EmployeeService: employees, PayrollService: payroll}

	mux.HandleFunc("/api/v1/payroll/calculate", withCORS(payrollHandler.Calculate))
	mux.HandleFunc("/api/v1/employees/payslips/zip", withCORS(payslipHandler.GenerateBulkZip))
	mux.HandleFunc("/api/v1/employees/", withCORS(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasSuffix(r.URL.Path, "/payslip") {
			payslipHandler.GeneratePDF(w, r)
			return
		}
		employeeHandler.HandleDetail(w, r)
	}))
	mux.HandleFunc("/api/v1/employees", withCORS(employeeHandler.HandleListCreate))

	return mux
}
