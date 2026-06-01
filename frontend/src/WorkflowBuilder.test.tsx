import '@testing-library/jest-dom'
import { fireEvent, render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { vi } from 'vitest'
import WorkflowBuilder from './WorkflowBuilder'

describe('WorkflowBuilder', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn((url) => {
      if (String(url).includes('/agents')) {
        return Promise.resolve({
          ok: true,
          json: () => Promise.resolve({ success: true, message: 'success', data: [{ id: 'agent-1', name: 'Agent One' }] }),
        })
      }
      return Promise.resolve({
        ok: true,
        json: () => Promise.resolve({ success: true, message: 'success', data: {} }),
      })
    }))
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders the workflow builder form and allows adding steps', async () => {
    render(
      <MemoryRouter initialEntries={['/workflow-builder']}>
        <Routes>
          <Route path="/workflow-builder" element={<WorkflowBuilder token="test-token" />} />
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByRole('heading', { name: /workflow builder/i })).toBeInTheDocument()
    await waitFor(() => expect(screen.getByRole('combobox')).toBeInTheDocument())

    fireEvent.change(screen.getByPlaceholderText(/add workflow step/i), { target: { value: 'Step 1' } })
    fireEvent.click(screen.getByRole('button', { name: /add step/i }))

    expect(screen.getAllByText(/Step 1/i).length).toBeGreaterThan(0)
  })
})
