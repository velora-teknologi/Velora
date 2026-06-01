import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import { BrowserRouter } from 'react-router-dom'
import App from './App'

test('renders Velora Dashboard text', () => {
  render(
    <BrowserRouter>
      <App />
    </BrowserRouter>,
  )
  expect(screen.getByText(/Velora Dashboard/i)).toBeInTheDocument()
})
