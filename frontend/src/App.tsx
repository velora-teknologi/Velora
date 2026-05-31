import { FormEvent, useEffect, useMemo, useState } from 'react'
import { getHealth, getMe, login } from './api'
import Agents from './Agents'
import Workflows from './Workflows'

function App() {
  const [status, setStatus] = useState('loading')
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [token, setToken] = useState<string | null>(null)
  const [user, setUser] = useState<{ email: string; full_name: string } | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [view, setView] = useState<'home' | 'agents' | 'workflows'>('home')

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
            <button onClick={() => setView('home')}>Home</button>
            <button onClick={() => setView('agents')} style={{ marginLeft: 8 }}>
              Agents
            </button>
            <button onClick={() => setView('workflows')} style={{ marginLeft: 8 }}>
              Workflows
            </button>
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
            {view === 'home' && (
              <>
                <h2>Getting started</h2>
                <ol>
                  <li>Use the seeded credentials to sign in.</li>
                  <li>Inspect backend status and user profile.</li>
                  <li>Extend the frontend with agents and workflows.</li>
                </ol>
              </>
            )}
            {view === 'agents' && token && <Agents token={token} />}
            {view === 'workflows' && token && <Workflows token={token} />}
          </section>
        )}
      </main>
    </div>
  )
}

export default App
