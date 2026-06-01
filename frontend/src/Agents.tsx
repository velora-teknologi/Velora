import { FormEvent, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import { createAgent, deleteAgent, listAgentWorkflows, listAgents } from './api'

export default function Agents({ token }: { token: string | null }) {
  const [agents, setAgents] = useState<any[]>([])
  const [selectedAgentId, setSelectedAgentId] = useState<string | null>(null)
  const [selectedAgentWorkflows, setSelectedAgentWorkflows] = useState<any[]>([])
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  const fetchAgents = async () => {
    setLoading(true)
    setError(null)
    try {
      const res = await listAgents(token)
      // API returns { data: [...] } in handlers
      setAgents(res.data || res)
    } catch (err: any) {
      setError(err.message)
    } finally {
      setLoading(false)
    }
  }

  useEffect(() => {
    if (!token) return
    fetchAgents()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [token])

  const handleCreate = async (e: FormEvent<HTMLFormElement>) => {
    e.preventDefault()
    setError(null)
    try {
      await createAgent(token, { name, description })
      setName('')
      setDescription('')
      fetchAgents()
    } catch (err: any) {
      setError(err.message)
    }
  }

  const handleDelete = async (id: string) => {
    if (!confirm('Delete this agent?')) return
    try {
      await deleteAgent(token, id)
      if (selectedAgentId === id) {
        setSelectedAgentId(null)
        setSelectedAgentWorkflows([])
      }
      fetchAgents()
    } catch (err: any) {
      setError(err.message)
    }
  }

  const handleViewWorkflows = async (id: string) => {
    setSelectedAgentId(id)
    setError(null)
    setLoading(true)
    try {
      const res = await listAgentWorkflows(token, id)
      setSelectedAgentWorkflows(res.data || res)
    } catch (err: any) {
      setError(err.message)
      setSelectedAgentWorkflows([])
    } finally {
      setLoading(false)
    }
  }

  return (
    <div>
      <h2>Agents</h2>
      <form onSubmit={handleCreate} style={{ marginBottom: 12 }}>
        <input placeholder="Name" value={name} onChange={(e) => setName(e.target.value)} required />
        <input placeholder="Description" value={description} onChange={(e) => setDescription(e.target.value)} />
        <button type="submit">Create</button>
      </form>
      {error && <div style={{ color: 'red' }}>{error}</div>}
      {loading ? (
        <div>Loading...</div>
      ) : (
        <>
          <ul>
            {agents.map((a) => (
              <li key={a.id} style={{ marginBottom: 8 }}>
                <strong>{a.name}</strong> — {a.description}
                <button style={{ marginLeft: 8 }} onClick={() => handleDelete(a.id)}>
                  Delete
                </button>
                <button style={{ marginLeft: 8 }} onClick={() => handleViewWorkflows(a.id)}>
                  View Workflows
                </button>
                <Link to={`/agents/${a.id}`} style={{ marginLeft: 8 }}>
                  View Details
                </Link>
              </li>
            ))}
          </ul>
          {selectedAgentId && (
            <section style={{ marginTop: 16 }}>
              <h3>Workflows for selected agent</h3>
              {selectedAgentWorkflows.length === 0 ? (
                <p>No workflows found.</p>
              ) : (
                <ul>
                  {selectedAgentWorkflows.map((workflow) => (
                    <li key={workflow.id} style={{ marginBottom: 8 }}>
                      <strong>{workflow.name}</strong> — {workflow.description}
                    </li>
                  ))}
                </ul>
              )}
            </section>
          )}
        </>
      )}
    </div>
  )
}
