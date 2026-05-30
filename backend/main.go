package main

import (
	"log"
	"net/http"

	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/routes"
	"github.com/berylxo/chalbeat-payroll/services"
)

func main() {

	// 1. Load default calculation rules
	rules := models.DefaultRules()

	// 2. Build engine
	calculator := &engine.Calculator{
		Rules: rules,
	}

	// 3. Build services
	payrollService := &services.PayrollService{
		Calculator: calculator,
	}

	employeeService := services.NewEmployeeService()

	// 4. Router
	r := routes.SetupRouter(payrollService, employeeService)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
