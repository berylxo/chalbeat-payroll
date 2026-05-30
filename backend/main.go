package main

import (
	"database/sql"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/routes"
	"github.com/berylxo/chalbeat-payroll/services"
)

func main() {
	// 1. Initialize database
	db, err := sql.Open("sqlite3", "payroll.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if err := initializeDatabase(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 2. Load default calculation rules
	rules := models.DefaultRules()

	// 3. Build engine
	calculator := &engine.Calculator{
		Rules: rules,
	}

	// 4. Build services
	payrollService := &services.PayrollService{
		Calculator: calculator,
	}

	employeeService := services.NewEmployeeService(db)

	// 5. Router
	r := routes.SetupRouter(payrollService, employeeService)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

func initializeDatabase(db *sql.DB) error {
	schema := `
	CREATE TABLE IF NOT EXISTS employees (
		employee_id TEXT PRIMARY KEY,
		name TEXT NOT NULL,
		phone_number TEXT,
		email TEXT,
		national_id TEXT,
		kra_pin TEXT NOT NULL,
		position TEXT NOT NULL,
		basic_pay REAL NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	return err
}
