package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"time" // Added for connection pooling durations

	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/routes"
	"github.com/berylxo/chalbeat-payroll/services"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// Read configuration from environment (useful for container deployments)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Require Postgres URL (Neon) via `DATABASE_URL`.
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required and must point to your Neon Postgres instance")
	}
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	defer db.Close()

	// Configure connection pooling for serverless database stability
	db.SetMaxOpenConns(10)                 // Limits concurrent connections to Neon
	db.SetMaxIdleConns(2)                  // Keeps a small pool of idle connections warm
	db.SetConnMaxLifetime(5 * time.Minute) // Recycles connections safely before serverless termination

	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	if err := initializeDatabase(db); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 3. Load default calculation rules
	rules := models.DefaultRules()

	// 4. Build engine
	calculator := &engine.Calculator{
		Rules: rules,
	}

	// 5. Build services
	payrollService := &services.PayrollService{
		Calculator: calculator,
	}

	employeeService := services.NewEmployeeService(db)

	// 6. Router Setup
	// SetupRouter registers backend API handlers (e.g., matching /api/...)
	r := routes.SetupRouter(payrollService, employeeService)

	baseURL := getBaseURL(port)
	log.Printf("Server running on %s", baseURL)
	log.Fatal(http.ListenAndServe(":"+port, r))
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
		created_at TIMESTAMPTZ DEFAULT now()
	);
	`
	_, err := db.Exec(schema)
	return err
}

func getBaseURL(port string) string {
	if v := os.Getenv("BASE_URL"); v != "" {
		return v
	}
	if v := os.Getenv("EXTERNAL_URL"); v != "" {
		return v
	}
	if v := os.Getenv("APP_URL"); v != "" {
		return v
	}

	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, a := range addrs {
			switch v := a.(type) {
			case *net.IPNet:
				ip := v.IP
				if ip == nil || ip.IsLoopback() {
					continue
				}
				if ip4 := ip.To4(); ip4 != nil {
					return "http://" + ip4.String() + ":" + port
				}
			case *net.IPAddr:
				ip := v.IP
				if ip == nil || ip.IsLoopback() {
					continue
				}
				if ip4 := ip.To4(); ip4 != nil {
					return "http://" + ip4.String() + ":" + port
				}
			}
		}
	}

	return "http://localhost:" + port
}
