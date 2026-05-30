import { BrowserRouter, Routes, Route, NavLink } from 'react-router-dom'
import EmployeeList from './pages/EmployeeList'
import EmployeeForm from './pages/EmployeeForm'

export const API_BASE = import.meta.env.VITE_API_BASE || 'http://localhost:8080'

function App() {
  return (
    <BrowserRouter>
      <div className="app-shell">

        <aside className="sidebar">
          <div className="brand">
            <div className="brand-name">CHALBEAT MEDICAL CENTRE</div>
            <div className="brand-sub">Payroll Management</div>
          </div>

          <nav>
            <NavLink to="/" end>
              Employee Directory
            </NavLink>
            <NavLink to="/new">
              Add Employee
            </NavLink>
          </nav>

          <div className="sidebar-footer">
            Kisumu, Kenya
          </div>
        </aside>

        <main className="content">
          <div className="header-panel">
            <div>
              <h1>Payroll Management</h1>
              <p>Manage employee records and generate compliant payslips.</p>
            </div>
          </div>

          <Routes>
            <Route path="/"        element={<EmployeeList />} />
            <Route path="/new"     element={<EmployeeForm />} />
            <Route path="/edit/:id" element={<EmployeeForm editMode />} />
          </Routes>
        </main>

      </div>
    </BrowserRouter>
  )
}

export default App