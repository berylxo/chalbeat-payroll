import { useEffect, useState } from 'react'
import { useNavigate, useParams, Link } from 'react-router-dom'
import { API_BASE } from '../App'

export default function EmployeeForm({ editMode = false }) {
  const navigate = useNavigate()
  const { id } = useParams()
  const [employee, setEmployee] = useState({
    name: '',
    kra_pin: '',
    position: '',
    basic_pay: '',
  })
  const [generatedId, setGeneratedId] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)
  const [message, setMessage] = useState(null)

  useEffect(() => {
    if (editMode && id) {
      fetchEmployee(id)
    }
  }, [editMode, id])

  async function fetchEmployee(employeeId) {
    setLoading(true)
    try {
      const res = await fetch(`${API_BASE}/api/v1/employees/${employeeId}`)
      if (!res.ok) throw new Error('Unable to load employee')
      setEmployee(await res.json())
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  async function handleSubmit(event) {
    event.preventDefault()
    setError(null)
    setMessage(null)
    setLoading(true)

    const payload = {
      name: employee.name,
      kra_pin: employee.kra_pin,
      position: employee.position,
      basic_pay: Number(employee.basic_pay),
    }

    try {
      const url = editMode ? `${API_BASE}/api/v1/employees/${id}` : `${API_BASE}/api/v1/employees`
      const method = editMode ? 'PUT' : 'POST'
      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload),
      })
      if (!res.ok) {
        const body = await res.json()
        throw new Error(body.error || 'Unable to save employee')
      }
      const saved = await res.json()
      if (!editMode) {
        setGeneratedId(saved.employee_id)
        setMessage(`Employee created successfully. ID: ${saved.employee_id}`)
        setTimeout(() => navigate('/'), 2000)
      } else {
        setMessage('Employee updated successfully.')
        setEmployee(saved)
      }
    } catch (err) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  function updateField(name, value) {
    setEmployee((prev) => ({ ...prev, [name]: value }))
  }

  return (
    <div className="page-card">
      <div className="toolbar">
        <div>
          <h2>{editMode ? 'Edit Employee' : 'Add New Employee'}</h2>
          <p>
            {editMode
              ? 'Update the employee record and save changes.'
              : 'Enter the details for a new team member.'}
          </p>
        </div>
        <Link className="button button-outline" to="/">
          ← Back to Directory
        </Link>
      </div>

      {/* Feedback strips sit flush between toolbar and form */}
      <div style={{ display: 'flex', flexDirection: 'column', gap: 10, padding: '16px 28px 0' }}>
        {error   && <p className="error-message"   style={{ margin: 0 }}>{error}</p>}
        {message && <p className="success-message" style={{ margin: 0 }}>{message}</p>}
        {generatedId && !editMode && (
          <div className="info-banner" style={{ margin: 0 }}>
            <strong>Employee ID:</strong> {generatedId}
          </div>
        )}
        {loading && <p style={{ color: 'var(--muted)', fontSize: 13, margin: 0 }}>Please wait…</p>}
      </div>

      <form className="form-grid" onSubmit={handleSubmit}>
        <label>
          <span>Full Name</span>
          <input
            value={employee.name}
            onChange={(e) => updateField('name', e.target.value)}
            required
            placeholder="Jane Doe"
          />
        </label>

        <label>
          <span>KRA PIN</span>
          <input
            value={employee.kra_pin}
            onChange={(e) => updateField('kra_pin', e.target.value)}
            required
            placeholder="A123456789X"
            style={{ fontFamily: 'monospace', letterSpacing: '0.5px' }}
          />
        </label>

        <label>
          <span>Job Position</span>
          <input
            value={employee.position}
            onChange={(e) => updateField('position', e.target.value)}
            required
            placeholder="Clinical Officer"
          />
        </label>

        <label>
          <span>Basic Salary (KES)</span>
          <input
            type="number"
            value={employee.basic_pay}
            onChange={(e) => updateField('basic_pay', e.target.value)}
            required
            min="0"
            step="0.01"
            placeholder="40000"
          />
        </label>

        <div className="form-actions">
          <button className="button button-primary" type="submit" disabled={loading}>
            {loading ? 'Saving…' : editMode ? 'Save Changes' : 'Create Employee'}
          </button>
          <Link className="button button-outline" to="/">
            Cancel
          </Link>
        </div>
      </form>
    </div>
  )
}