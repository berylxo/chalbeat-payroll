package main

import (
	"database/sql"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/berylxo/chalbeat-payroll/engine"
	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/routes"
	"github.com/berylxo/chalbeat-payroll/services"
	_ "modernc.org/sqlite"
)

func main() {
	// Read configuration from environment (useful for container deployments)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		// Ensure the default database file lives inside the backend folder
		dbPath = filepath.Join("backend", "payroll.db")
	}

	// 1. Ensure database parent directory exists before initialization.
	// Fly.io mounts override permissions on /data; this guarantees the application
	// has structural access to create the db file.
	dbDir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		log.Fatalf("Failed to create database directory %s: %v", dbDir, err)
	}

	// 2. Initialize database
	db, err := sql.Open("sqlite", dbPath)
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
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(schema)
	return err
}

func getBaseURL(port string) string {
	// Prefer explicit configuration if provided by the environment
	if v := os.Getenv("BASE_URL"); v != "" {
		return v
	}
	if v := os.Getenv("EXTERNAL_URL"); v != "" {
		return v
	}
	if v := os.Getenv("APP_URL"); v != "" {
		return v
	}

	// Attempt to pick a non-loopback IPv4 address from the host
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

	// Final fallback
	return "http://localhost:" + port
}
