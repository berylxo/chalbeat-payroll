package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/berylxo/chalbeat-payroll/models"
	"github.com/berylxo/chalbeat-payroll/services"
)

type EmployeeHandler struct {
	Service *services.EmployeeService
}

func (h *EmployeeHandler) HandleListCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		h.list(w, r)
		return
	}

	if r.Method == http.MethodPost {
		h.create(w, r)
		return
	}

	h.methodNotAllowed(w, r)
}

func (h *EmployeeHandler) HandleDetail(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/api/v1/employees/")
	id = strings.TrimSuffix(id, "/")
	if id == "" {
		h.methodNotAllowed(w, r)
		return
	}

	if r.Method == http.MethodGet {
		h.get(w, r, id)
		return
	}

	if r.Method == http.MethodPut {
		h.update(w, r, id)
		return
	}

	h.methodNotAllowed(w, r)
}

func (h *EmployeeHandler) list(w http.ResponseWriter, r *http.Request) {
	employees := h.Service.List()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(employees)
}

func (h *EmployeeHandler) create(w http.ResponseWriter, r *http.Request) {
	var emp models.Employee
	json.NewDecoder(r.Body).Decode(&emp)

	if err := h.Service.Create(emp); err != nil {
		h.writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(emp)
}

func (h *EmployeeHandler) get(w http.ResponseWriter, r *http.Request, id string) {
	emp, ok := h.Service.Get(id)
	if !ok {
		h.writeError(w, services.ErrEmployeeNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(emp)
}

func (h *EmployeeHandler) update(w http.ResponseWriter, r *http.Request, id string) {
	var emp models.Employee
	json.NewDecoder(r.Body).Decode(&emp)

	updated, err := h.Service.Update(id, emp)
	if err != nil {
		h.writeError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updated)
}

func (h *EmployeeHandler) writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	if err == services.ErrEmployeeNotFound {
		status = http.StatusNotFound
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func (h *EmployeeHandler) methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusMethodNotAllowed)
}
