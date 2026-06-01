import { FormEvent, useEffect, useMemo, useState } from 'react'
import { NavLink, Navigate, Route, Routes } from 'react-router-dom'
import { getHealth, getMe, login } from './api'
import Agents from './Agents'
import AgentDetail from './AgentDetail'
import Workflows from './Workflows'
import WorkflowBuilder from './WorkflowBuilder'
import WorkflowDetail from './WorkflowDetail'

function App() {
  const [status, setStatus] = useState('loading')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<{ email: string; full_name: string } | null>(null)
  const [error, setError] = useState<string | null>(null)

  const isAuthenticated = useMemo(() => !!token, [token])

  useEffect(() => {
    const savedToken = localStorage.getItem('velora_token')
    if (savedToken) {
      setToken(savedToken)
    }
  }, [])

  useEffect(() => {
    if (!token) {
      setUser(null)
      return
    }

    getMe(token).then((profile) => setUser(profile)).catch(() => {
      setToken(null)
      localStorage.removeItem('velora_token')
    })
  }, [token])

  useEffect(() => {
    getHealth(token ?? undefined)
      .then((result) => setStatus(result.status))
      .catch(() => setStatus('unavailable'))
  }, [token])

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setError(null)

    try {
      const result = await login(email, password)
      setToken(result.token)
      localStorage.setItem('velora_token', result.token)
      setUser(result.user)
      setEmail('')
      setPassword('')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Login failed')
    }
  }

  const handleLogout = () => {
    setToken(null)
    setUser(null)
    localStorage.removeItem('velora_token')
  }

  return (
    <div className="app-container">
      <header>
        <h1>Velora Dashboard</h1>
        <p>AI automation platform running with frontend and backend services.</p>
      </header>

      <main>
        <section className="status-card">
          <h2>Backend health</h2>
          <p>
            Current status: <strong>{status}</strong>
          </p>
        </section>

        {isAuthenticated && user && (
          <section className="status-card">
            <h2>Signed in as</h2>
            <p>{user.full_name}</p>
            <p>{user.email}</p>
            <button type="button" onClick={handleLogout}>
              Sign Out
            </button>
          </section>
        )}

        {isAuthenticated && (
          <nav style={{ marginTop: 12 }}>
            <NavLink to="/" style={{ marginRight: 8 }}>
              Home
            </NavLink>
            <NavLink to="/agents" style={{ marginRight: 8 }}>
              Agents
            </NavLink>
            <NavLink to="/workflows" style={{ marginRight: 8 }}>
              Workflows
            </NavLink>
            <NavLink to="/workflow-builder">Workflow Builder</NavLink>
          </nav>
        )}

        {!isAuthenticated ? (
          <section className="feature-card">
            <h2>Sign in to Velora</h2>
            <form className="login-form" onSubmit={handleSubmit}>
              <label>
                Email
                <input
                  type="email"
                  value={email}
                  onChange={(e) => setEmail(e.target.value)}
                  required
                />
              </label>
              <label>
                Password
                <input
                  type="password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  required
                />
              </label>
              <button type="submit">Sign In</button>
              {error && <p className="error-message">{error}</p>}
            </form>
          </section>
        ) : (
          <section className="feature-card">
            <Routes>
              <Route
                path="/"
                element={
                  <>
                    <h2>Getting started</h2>
                    <ol>
                      <li>Use the seeded credentials to sign in.</li>
                      <li>Inspect backend status and user profile.</li>
                      <li>Extend the frontend with agents and workflows.</li>
                    </ol>
                  </>
                }
              />
              <Route path="/agents" element={<Agents token={token} />} />
              <Route path="/agents/:id" element={<AgentDetail token={token} />} />
              <Route path="/workflows" element={<Workflows token={token} />} />
              <Route path="/workflow-builder" element={<WorkflowBuilder token={token} />} />
              <Route path="/workflows/:id" element={<WorkflowDetail token={token} />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </Routes>
          </section>
        )}
      </main>
    </div>
  )
}

export default App
