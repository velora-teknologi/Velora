import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { getAgent, listAgentWorkflows } from './api'

export default function AgentDetail({ token }: { token: string | null }) {
  const { id } = useParams<{ id: string }>()
  const [agent, setAgent] = useState<any | null>(null)
  const [workflows, setWorkflows] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!id || !token) return

    const fetchDetails = async () => {
      setLoading(true)
      setError(null)
      try {
        const agentResponse = await getAgent(token, id)
        setAgent(agentResponse.data || agentResponse)
        const workflowResponse = await listAgentWorkflows(token, id)
        setWorkflows(workflowResponse.data || workflowResponse)
      } catch (err: any) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }

    fetchDetails()
  }, [id, token])

  if (!token) {
    return (
      <div>
        <p>Authentication is required to view agent details.</p>
        <Link to="/agents">Back to agents</Link>
      </div>
    )
  }

  if (!id) {
    return (
      <div>
        <p>Agent ID is missing from the route.</p>
        <Link to="/agents">Back to agents</Link>
      </div>
    )
  }

  return (
    <div>
      <h2>Agent details</h2>
      <Link to="/agents" style={{ display: 'inline-block', marginBottom: 12 }}>
        Back to agents
      </Link>
      {loading ? (
        <p>Loading agent details...</p>
      ) : error ? (
        <p style={{ color: 'red' }}>{error}</p>
      ) : agent ? (
        <div>
          <section style={{ marginBottom: 16 }}>
            <h3>{agent.name}</h3>
            <p>{agent.description}</p>
            <p>
              <strong>Status:</strong> {agent.status}
            </p>
            <div>
              <strong>Config:</strong>
              <pre style={{ background: '#f7f7f7', padding: 8 }}>{JSON.stringify(agent.config, null, 2)}</pre>
            </div>
          </section>

          <section>
            <h3>Workflows</h3>
            {workflows.length === 0 ? (
              <p>No workflows found for this agent.</p>
            ) : (
              <ul>
                {workflows.map((workflow) => (
                  <li key={workflow.id} style={{ marginBottom: 8 }}>
                    <strong>{workflow.name}</strong> — {workflow.description}
                  </li>
                ))}
              </ul>
            )}
          </section>
        </div>
      ) : (
        <p>Agent not found.</p>
      )}
    </div>
  )
}
