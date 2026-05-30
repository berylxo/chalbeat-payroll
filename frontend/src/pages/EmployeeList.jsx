import { useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { API_BASE } from '../App'

function downloadBlob(blob, fileName) {
  const url = window.URL.createObjectURL(blob)
  const anchor = document.createElement('a')
  anchor.href = url
  anchor.download = fileName
  document.body.appendChild(anchor)
  anchor.click()
  anchor.remove()
  window.URL.revokeObjectURL(url)
}

export default function EmployeeList() {
  const [employees, setEmployees] = useState([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState(null)

  useEffect(() => {
    fetchEmployees()
  }, [])

  async function fetchEmployees() {
    setLoading(true)
    setError(null)
    try {
      const res = await fetch(`${API_BASE}/api/v1/employees`)
      if (!res.ok) throw new Error('Unable to load employees')
      setEmployees(await res.json())
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function downloadPayslip(id) {
    window.open(`${API_BASE}/api/v1/employees/${id}/payslip`, '_blank')
  }

  async function downloadZip() {
    try {
      const res = await fetch(`${API_BASE}/api/v1/employees/payslips/zip`)
      if (!res.ok) throw new Error('Unable to create zip')
      const blob = await res.blob()
      downloadBlob(blob, 'Payslips.zip')
    } catch (err) {
      setError(err.message)
    }
  }

  return (
    <div className="page-card">
      <div className="toolbar">
        <h2>Employee Directory</h2>
        <div className="toolbar-actions">
          <Link className="button button-primary" to="/new">
            Add Employee
          </Link>
          <button className="button button-secondary" onClick={downloadZip}>
            Download All Payslips
          </button>
        </div>
      </div>

      {loading && <p>Loading employees…</p>}
      {error && <p className="error-message">{error}</p>}
      {!loading && employees.length === 0 && <p>No employees yet. Add one to get started.</p>}

      {employees.length > 0 && (
        <div className="table-wrapper">
          <table>
            <thead>
              <tr>
                <th>ID</th>
                <th>Name</th>
                <th>KRA PIN</th>
                <th>Position</th>
                <th>Basic Pay</th>
                <th>Actions</th>
              </tr>
            </thead>
            <tbody>
              {employees.map((employee) => (
                <tr key={employee.id}>
                  <td>{employee.id}</td>
                  <td>{employee.name}</td>
                  <td>{employee.kra_pin}</td>
                  <td>{employee.position}</td>
                  <td>{employee.basic_pay.toFixed(2)}</td>
                  <td className="actions-cell">
                    <Link className="button button-small" to={`/edit/${employee.id}`}>
                      Edit
                    </Link>
                    <button className="button button-small button-outline" onClick={() => downloadPayslip(employee.id)}>
                      Payslip
                    </button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  )
}
