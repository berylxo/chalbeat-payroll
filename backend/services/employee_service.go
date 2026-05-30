package services

import (
	"errors"
	"sort"
	"sync"

	"github.com/berylxo/chalbeat-payroll/models"
)

var ErrEmployeeNotFound = errors.New("employee not found")
var ErrEmployeeExists = errors.New("employee already exists")
var ErrEmployeeIDRequired = errors.New("employee id is required")

type EmployeeService struct {
	mu    sync.RWMutex
	store map[string]models.Employee
}

func NewEmployeeService() *EmployeeService {
	return &EmployeeService{store: make(map[string]models.Employee)}
}

func (s *EmployeeService) List() []models.Employee {
	s.mu.RLock()
	defer s.mu.RUnlock()

	employees := make([]models.Employee, 0, len(s.store))
	for _, emp := range s.store {
		employees = append(employees, emp)
	}
	sort.SliceStable(employees, func(i, j int) bool {
		return employees[i].ID < employees[j].ID
	})
	return employees
}

func (s *EmployeeService) Get(id string) (models.Employee, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	emp, ok := s.store[id]
	return emp, ok
}

func (s *EmployeeService) Create(emp models.Employee) error {
	if emp.ID == "" {
		return ErrEmployeeIDRequired
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[emp.ID]; ok {
		return ErrEmployeeExists
	}
	s.store[emp.ID] = emp
	return nil
}

func (s *EmployeeService) Update(id string, emp models.Employee) (models.Employee, error) {
	if id == "" {
		return models.Employee{}, ErrEmployeeIDRequired
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.store[id]; !ok {
		return models.Employee{}, ErrEmployeeNotFound
	}

	emp.ID = id
	s.store[id] = emp
	return emp, nil
}
