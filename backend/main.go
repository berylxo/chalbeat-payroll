package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
		dbPath = "payroll.db"
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

	// 7. Isolated SPA Frontend Routing Strategy
	distDir := os.Getenv("FRONTEND_DIST")
	if distDir == "" {
		distDir = "../frontend/dist"
	}
	fileServer := http.FileServer(http.Dir(distDir))

	// Register a catch-all route handler.
	// NOTE: If routes.SetupRouter uses gorilla/mux, replace this statement with:
	// r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) { ... })
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Stop the frontend catch-all from intercepting structural API paths
		if strings.HasPrefix(req.URL.Path, "/api") {
			http.NotFound(w, req)
			return
		}

		// Handle explicit file requests
		relPath := strings.TrimPrefix(req.URL.Path, "/")
		if relPath == "" {
			// Serve fallback core entrypoint
			http.ServeFile(w, req, filepath.Join(distDir, "index.html"))
			return
		}

		requested := filepath.Join(distDir, relPath)
		if info, err := os.Stat(requested); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, req)
			return
		}

		// Fallback for client-side routing hydration (React, Vue, Svelte, etc.)
		http.ServeFile(w, req, filepath.Join(distDir, "index.html"))
	})

	log.Printf("Server running on http://localhost:%s", port)
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
