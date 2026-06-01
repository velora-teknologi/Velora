import '@testing-library/jest-dom'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter, Route, Routes } from 'react-router-dom'
import { vi } from 'vitest'
import WorkflowDetail from './WorkflowDetail'

describe('WorkflowDetail', () => {
  beforeEach(() => {
    vi.stubGlobal('fetch', vi.fn(() =>
      Promise.resolve({
        ok: true,
        json: () => Promise.resolve({
          success: true,
          message: 'success',
          data: {
            id: 'workflow-1',
            name: 'Workflow One',
            description: 'Test workflow',
            definition: { steps: ['start'] },
            is_active: true,
          },
        }),
      }),
    ))
  })

  afterEach(() => {
    vi.restoreAllMocks()
  })

  it('renders workflow details when authenticated', async () => {
    render(
      <MemoryRouter initialEntries={['/workflows/workflow-1']}>
        <Routes>
          <Route path="/workflows/:id" element={<WorkflowDetail token="test-token" />} />
        </Routes>
      </MemoryRouter>,
    )

    expect(screen.getByRole('heading', { name: /workflow details/i })).toBeInTheDocument()
    await waitFor(() => expect(screen.getByText(/Workflow One/i)).toBeInTheDocument())
    expect(screen.getByText(/Test workflow/i)).toBeInTheDocument()
    expect(screen.getByText(/Active:/i)).toBeInTheDocument()
  })
})
