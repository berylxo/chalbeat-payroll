package services

import (
	"database/sql"
	"errors"
	"fmt"
	"sync"

	"github.com/berylxo/chalbeat-payroll/models"
)

var ErrEmployeeNotFound = errors.New("employee not found")
var ErrEmployeeExists = errors.New("employee already exists")

type EmployeeService struct {
	db   *sql.DB
	mu   sync.Mutex
}
func NewEmployeeService(db *sql.DB) *EmployeeService {
	return &EmployeeService{db: db}
}

func (s *EmployeeService) generateEmployeeID() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	row := s.db.QueryRow("SELECT COUNT(*) FROM employees")
	var count int
	_ = row.Scan(&count)

	id := fmt.Sprintf("EMP-%d", count+1)
	return id
}

func (s *EmployeeService) List() []models.Employee {
	rows, err := s.db.Query("SELECT employee_id, name, phone_number, email, national_id, kra_pin, position, basic_pay FROM employees ORDER BY employee_id")
	if err != nil {
		return []models.Employee{}
	}
	defer rows.Close()

	var employees []models.Employee
	for rows.Next() {
		var emp models.Employee
		err := rows.Scan(
			&emp.EmployeeID,
			&emp.Name,
			&emp.PhoneNumber,
			&emp.Email,
			&emp.NationalID,
			&emp.KraPin,
			&emp.Position,
			&emp.BasicPay,
		)
		if err == nil {
			employees = append(employees, emp)
		}
	}
	return employees
}

func (s *EmployeeService) Get(id string) (models.Employee, bool) {
	var emp models.Employee
	row := s.db.QueryRow(
		"SELECT employee_id, name, phone_number, email, national_id, kra_pin, position, basic_pay FROM employees WHERE employee_id = $1",
		id,
	)
	err := row.Scan(
		&emp.EmployeeID,
		&emp.Name,
		&emp.PhoneNumber,
		&emp.Email,
		&emp.NationalID,
		&emp.KraPin,
		&emp.Position,
		&emp.BasicPay,
	)
	if err != nil {
		return models.Employee{}, false
	}
	return emp, true
}

func (s *EmployeeService) Create(emp models.Employee) (models.Employee, error) {
	emp.EmployeeID = s.generateEmployeeID()

	_, err := s.db.Exec(
		"INSERT INTO employees (employee_id, name, phone_number, email, national_id, kra_pin, position, basic_pay) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)",
		emp.EmployeeID,
		emp.Name,
		emp.PhoneNumber,
		emp.Email,
		emp.NationalID,
		emp.KraPin,
		emp.Position,
		emp.BasicPay,
	)
	if err != nil {
		return models.Employee{}, err
	}
	return emp, nil
}

func (s *EmployeeService) Update(id string, emp models.Employee) (models.Employee, error) {
	if id == "" {
		return models.Employee{}, errors.New("employee id is required")
	}

	_, err := s.db.Exec(
		"UPDATE employees SET name = $1, phone_number = $2, email = $3, national_id = $4, kra_pin = $5, position = $6, basic_pay = $7 WHERE employee_id = $8",
		emp.Name,
		emp.PhoneNumber,
		emp.Email,
		emp.NationalID,
		emp.KraPin,
		emp.Position,
		emp.BasicPay,
		id,
	)
	if err != nil {
		return models.Employee{}, err
	}

	emp.EmployeeID = id
	return emp, nil
}
