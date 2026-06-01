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

export async function getMe(token: string | null) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/users/me`, {
    headers: authHeaders(token),
  })

  if (!response.ok) {
    const error = await response.json().catch(() => null)
    throw new Error(error?.error ?? 'Unable to fetch user')
  }

  return response.json()
}

// Agents
export async function listAgents(token: string | null, skip = 0, limit = 100) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/agents?skip=${skip}&limit=${limit}`, {
    headers: authHeaders(token),
  })
  if (!response.ok) throw new Error('Failed to list agents')
  return response.json()
}

export async function createAgent(token: string | null, payload: { name: string; description?: string; config?: any }) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/agents`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const err = await response.json().catch(() => null)
    throw new Error(err?.error ?? 'Failed to create agent')
  }
  return response.json()
}

export async function getAgent(token: string, id: string) {
  const response = await fetch(`${API_URL}/agents/${id}`, {
    headers: authHeaders(token),
  })
  if (!response.ok) {
    const err = await response.json().catch(() => null)
    throw new Error(err?.error ?? 'Failed to fetch agent')
  }
  return response.json()
}

export async function deleteAgent(token: string | null, id: string) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/agents/${id}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  })
  if (!response.ok) throw new Error('Failed to delete agent')
  return true
}

// Workflows
export async function listWorkflows(token: string | null, agentId: string, skip = 0, limit = 100) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/workflows?agent_id=${agentId}&skip=${skip}&limit=${limit}`, {
    headers: authHeaders(token),
  })
  if (!response.ok) throw new Error('Failed to list workflows')
  return response.json()
}

export async function listAgentWorkflows(token: string | null, agentId: string, skip = 0, limit = 100) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/agents/${agentId}/workflows?skip=${skip}&limit=${limit}`, {
    headers: authHeaders(token),
  })
  if (!response.ok) throw new Error('Failed to list agent workflows')
  return response.json()
}

export async function getWorkflow(token: string | null, id: string) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/workflows/${id}`, {
    headers: authHeaders(token),
  })
  if (!response.ok) {
    const err = await response.json().catch(() => null)
    throw new Error(err?.error ?? 'Failed to fetch workflow')
  }
  return response.json()
}

export async function executeWorkflow(token: string | null, id: string) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/workflows/${id}/execute`, {
    method: 'POST',
    headers: authHeaders(token),
  })

  if (!response.ok) {
    const err = await response.json().catch(() => null)
    throw new Error(err?.error ?? 'Failed to execute workflow')
  }
  return response.json()
}

export async function createWorkflow(token: string | null, payload: { name: string; description?: string; agent_id: string; definition: any }) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/workflows`, {
    method: 'POST',
    headers: authHeaders(token),
    body: JSON.stringify(payload),
  })
  if (!response.ok) {
    const err = await response.json().catch(() => null)
    throw new Error(err?.error ?? 'Failed to create workflow')
  }
  return response.json()
}

export async function deleteWorkflow(token: string | null, id: string) {
  if (!token) throw new Error('Authentication token is required')
  const response = await fetch(`${API_URL}/workflows/${id}`, {
    method: 'DELETE',
    headers: authHeaders(token),
  })
  if (!response.ok) throw new Error('Failed to delete workflow')
  return true
}
