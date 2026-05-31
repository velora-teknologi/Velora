import { FormEvent, useEffect, useState } from 'react'
import { createAgent, deleteAgent, listAgents } from './api'

export default function Agents({ token }: { token: string }) {
  const [agents, setAgents] = useState<any[]>([])
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
    fetchAgents()
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

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
      fetchAgents()
    } catch (err: any) {
      setError(err.message)
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
        <ul>
          {agents.map((a) => (
            <li key={a.id} style={{ marginBottom: 8 }}>
              <strong>{a.name}</strong> — {a.description}
              <button style={{ marginLeft: 8 }} onClick={() => handleDelete(a.id)}>
                Delete
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  )
}
