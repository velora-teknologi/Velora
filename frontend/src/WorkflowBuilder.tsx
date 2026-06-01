import { FormEvent, useEffect, useState } from 'react'
import { createWorkflow, listAgents } from './api'

export default function WorkflowBuilder({ token }: { token: string | null }) {
  const [agents, setAgents] = useState<any[]>([])
  const [agentId, setAgentId] = useState('')
  const [name, setName] = useState('')
  const [description, setDescription] = useState('')
  const [stepInput, setStepInput] = useState('')
  const [steps, setSteps] = useState<string[]>([])
  const [status, setStatus] = useState<string | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!token) return

    const loadAgents = async () => {
      try {
        const res = await listAgents(token)
        setAgents(res.data || res)
      } catch (err: any) {
        setError(err.message)
      }
    }

    loadAgents()
  }, [token])

  const addStep = () => {
    if (!stepInput.trim()) return
    setSteps((current) => [...current, stepInput.trim()])
    setStepInput('')
  }

  const removeStep = (index: number) => {
    setSteps((current) => current.filter((_, idx) => idx !== index))
  }

  const handleSubmit = async (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault()
    setError(null)
    setStatus(null)

    if (!agentId) {
      setError('Please select an agent before building a workflow.')
      return
    }

    try {
      await createWorkflow(token, {
        name,
        description,
        agent_id: agentId,
        definition: { steps },
      })
      setStatus('Workflow created successfully')
      setName('')
      setDescription('')
      setSteps([])
    } catch (err: any) {
      setError(err.message)
    }
  }

  if (!token) {
    return <p>Authentication is required to access the workflow builder.</p>
  }

  return (
    <div>
      <h2>Workflow Builder</h2>
      <p>Build a reusable workflow for an existing agent with step-based definition data.</p>

      {error && <div style={{ color: 'red' }}>{error}</div>}
      {status && <div style={{ color: 'green' }}>{status}</div>}

      <form onSubmit={handleSubmit} style={{ marginBottom: 16 }}>
        <div style={{ marginBottom: 8 }}>
          <label>
            Workflow name
            <input value={name} onChange={(e) => setName(e.target.value)} required />
          </label>
        </div>

        <div style={{ marginBottom: 8 }}>
          <label>
            Description
            <input value={description} onChange={(e) => setDescription(e.target.value)} />
          </label>
        </div>

        <div style={{ marginBottom: 8 }}>
          <label>
            Agent
            <select value={agentId} onChange={(e) => setAgentId(e.target.value)} required>
              <option value="">Select an agent</option>
              {agents.map((agent) => (
                <option key={agent.id} value={agent.id}>
                  {agent.name}
                </option>
              ))}
            </select>
          </label>
        </div>

        <div style={{ marginBottom: 8 }}>
          <label>
            Step name
            <input value={stepInput} onChange={(e) => setStepInput(e.target.value)} placeholder="Add workflow step" />
          </label>
          <button type="button" onClick={addStep} style={{ marginLeft: 8 }}>
            Add step
          </button>
        </div>

        <div style={{ marginBottom: 12 }}>
          <strong>Workflow steps</strong>
          <ul>
            {steps.map((step, index) => (
              <li key={`${step}-${index}`} style={{ marginBottom: 4 }}>
                {step}
                <button type="button" onClick={() => removeStep(index)} style={{ marginLeft: 8 }}>
                  Remove
                </button>
              </li>
            ))}
          </ul>
        </div>

        <button type="submit">Save workflow</button>
      </form>

      <section>
        <h3>Preview definition</h3>
        <pre style={{ background: '#f7f7f7', padding: 12 }}>{JSON.stringify({ steps }, null, 2)}</pre>
      </section>
    </div>
  )
}
