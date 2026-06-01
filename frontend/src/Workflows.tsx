import { FormEvent, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { createWorkflow, deleteWorkflow, listAgents, listWorkflows } from './api'

export default function Workflows({ token }: { token: string | null }) {
  const [agentId, setAgentId] = useState('')
  const [workflows, setWorkflows] = useState<any[]>([])
  const [name, setName] = useState('')
  const [definition, setDefinition] = useState('{}')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const [agents, setAgents] = useState<any[]>([])

  const fetchAgents = async () => {
    setError(null)
    try {
      const res = await listAgents(token)
      setAgents(res.data || res)
    } catch (err: any) {
      setError(err.message)
    }
  }

  const fetchWorkflows = async () => {
    if (!agentId) return
    setLoading(true)
    setError(null)
    try {
      const res = await listWorkflows(token, agentId)
      setWorkflows(res.data || res)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (token) {
      fetchAgents()
    }
  }, [token])

  const handleCreate = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError(null)
    try {
      const parsed = JSON.parse(definition)
      await createWorkflow(token, { name, agent_id: agentId, definition: parsed })
      setName('')
      setDefinition('{}')
      fetchWorkflows()
    } catch (err: any) {
      setError(err.message)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this workflow?')) return
    try {
      await deleteWorkflow(token, id)
      fetchWorkflows()
    } catch (err: any) {
      setError(err.message)
    }
  }

  return (
    <div>
      <h2>Workflows</h2>
      <div style={{ marginBottom: 12 }}>
        <select value={agentId} onChange={(e) => setAgentId(e.target.value)}>
          <option value="">Select an agent</option>
          {agents.map((agent) => (
            <option key={agent.id} value={agent.id}>
              {agent.name}
            </option>
          ))}
        </select>
        <button onClick={fetchWorkflows} style={{ marginLeft: 8 }} disabled={!agentId}>
          Load
        </button>
      </div>

      {agentId && (
        <form onSubmit={handleCreate} style={{ marginBottom: 12 }}>
          <input placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} required />
          <textarea placeholder='Definition JSON' value={definition} onChange={(e) => setDefinition(e.target.value)} />
          <button type="submit">Create</button>
        </form>
      )}

      {error && <div style={{ color: 'red' }}>{error}</div>}

      {loading ? (
        <div>Loading...</div>
      ) : (
        <ul>
          {workflows.map((w) => (
            <li key={w.id} style={{ marginBottom: 8 }}>
              <strong>{w.name}</strong> — {JSON.stringify(w.definition)}
              <Link to={`/workflows/${w.id}`} style={{ marginLeft: 8 }}>
                View Details
              </Link>
              <button style={{ marginLeft: 8 }} onClick={() => handleDelete(w.id)}>
                Delete
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
