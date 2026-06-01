import { useEffect, useState } from 'react'
import { Link, useParams } from 'react-router-dom'
import { executeWorkflow, getWorkflow } from './api'

export default function WorkflowDetail({ token }: { token: string | null }) {
  const { id } = useParams<{ id: string }>()
  const [workflow, setWorkflow] = useState<any | null>(null)
  const [loading, setLoading] = useState(true)
  const [error, setError] = useState<string | null>(null)
  const [executionStatus, setExecutionStatus] = useState<string | null>(null)
  const [executionError, setExecutionError] = useState<string | null>(null)

  useEffect(() => {
    if (!id || !token) return

    const loadWorkflow = async () => {
      setLoading(true)
      setError(null)
      try {
        const response = await getWorkflow(token, id)
        setWorkflow(response.data || response)
      } catch (err: any) {
        setError(err.message)
      } finally {
        setLoading(false)
      }
    }

    loadWorkflow()
  }, [id, token])

  if (!token) {
    return (
      <div>
        <p>Authentication is required to view workflow details.</p>
        <Link to="/workflows">Back to workflows</Link>
      </div>
    )
  }

  if (!id) {
    return (
      <div>
        <p>Workflow ID is missing from the route.</p>
        <Link to="/workflows">Back to workflows</Link>
      </div>
    )
  }

  return (
    <div>
      <h2>Workflow details</h2>
      <Link to="/workflows" style={{ display: 'inline-block', marginBottom: 12 }}>
        Back to workflows
      </Link>
      {loading ? (
        <p>Loading workflow details...</p>
      ) : error ? (
        <p style={{ color: 'red' }}>{error}</p>
      ) : workflow ? (
        <div>
          <section style={{ marginBottom: 16 }}>
            <h3>{workflow.name}</h3>
            <p>{workflow.description}</p>
            <p>
              <strong>Active:</strong> {workflow.is_active ? 'Yes' : 'No'}
            </p>
            <div>
              <strong>Definition</strong>
              <pre style={{ background: '#f7f7f7', padding: 8 }}>{JSON.stringify(workflow.definition, null, 2)}</pre>
            </div>
          </section>
          {executionError && <p style={{ color: 'red' }}>{executionError}</p>}
          {executionStatus && <p style={{ color: 'green' }}>{executionStatus}</p>}
          <button
            type="button"
            onClick={async () => {
              setExecutionStatus(null)
              setExecutionError(null)
              try {
                const result = await executeWorkflow(token, id)
                setExecutionStatus(result.data?.status || result.status || 'Workflow queued')
              } catch (err: any) {
                setExecutionError(err.message)
              }
            }}
            style={{ marginBottom: 16 }}
          >
            Execute workflow
          </button>
        </div>
      ) : (
        <p>Workflow not found.</p>
      )}
    </div>
  )
}
