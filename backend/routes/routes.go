package routes

import (
	"net/http"
	"os"
	"path/filepath"
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

	// Serve built frontend for all non-API routes (SPA support).
	distDir := os.Getenv("FRONTEND_DIST")
	if distDir == "" {
		distDir = "../frontend/dist"
	}
	fileServer := http.FileServer(http.Dir(distDir))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Ensure API paths are not intercepted
		if strings.HasPrefix(r.URL.Path, "/api") {
			http.NotFound(w, r)
			return
		}

		// Serve explicit files when present
		relPath := strings.TrimPrefix(r.URL.Path, "/")
		if relPath == "" {
			http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
			return
		}

		requested := filepath.Join(distDir, relPath)
		if info, err := os.Stat(requested); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for client-side routing
		http.ServeFile(w, r, filepath.Join(distDir, "index.html"))
	})

	return mux
}
