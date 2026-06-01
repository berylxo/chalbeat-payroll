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
	_ "github.com/mattn/go-sqlite3"
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

	// Serve frontend static files (production build) with SPA fallback.
	// Routes under /api/ are registered on the mux first; this handler
	// catches other paths and serves files from frontend/dist, falling
	// back to index.html so React Router can handle client-side routes.
	distDir := "../frontend/dist"
	fileServer := http.FileServer(http.Dir(distDir))
	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// If the requested path corresponds to an existing file in distDir, serve it.
		// Note: req.URL.Path begins with "/", so trim it before joining to avoid
		// filepath.Join treating it as an absolute path and ignoring distDir.
		relPath := strings.TrimPrefix(req.URL.Path, "/")
		if relPath == "" {
			// Root path -> serve index
			http.ServeFile(w, req, filepath.Join(distDir, "index.html"))
			return
		}
		requested := filepath.Join(distDir, relPath)
		if info, err := os.Stat(requested); err == nil && !info.IsDir() {
			fileServer.ServeHTTP(w, req)
			return
		}
		// Fallback to index.html for client-side routes
		http.ServeFile(w, req, filepath.Join(distDir, "index.html"))
	})

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
