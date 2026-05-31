const API_URL = import.meta.env.VITE_API_URL ?? 'http://localhost:8080/api/v1'
const BACKEND_URL = API_URL.replace(/\/api\/v1\/?$/, '')

const authHeaders = (token?: string) => {
  const headers: Record<string, string> = {
    'Content-Type': 'application/json',
  }
  if (token) {
    headers.Authorization = `Bearer ${token}`
  }
  return headers
}

export async function getHealth(token?: string) {
  const response = await fetch(`${BACKEND_URL}/health`, {
    headers: authHeaders(token),
  })
  if (!response.ok) {
    throw new Error('Health check failed')
  }
  return response.json()
}

export async function login(email: string, password: string) {
  const response = await fetch(`${API_URL}/auth/login`, {
    method: 'POST',
    headers: authHeaders(),
    body: JSON.stringify({ email, password }),
  })

  if (!response.ok) {
    const error = await response.json().catch(() => null)
    throw new Error(error?.error ?? 'Login failed')
  }

  return response.json()
}

export async function getMe(token: string) {
  const response = await fetch(`${API_URL}/users/me`, {
    headers: authHeaders(token),
  })

  if (!response.ok) {
    const error = await response.json().catch(() => null)
    throw new Error(error?.error ?? 'Unable to fetch user')
  }

  return response.json()
}
