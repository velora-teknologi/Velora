import '@testing-library/jest-dom'
import { render, screen } from '@testing-library/react'
import App from './App'

test('renders Velora Dashboard text', () => {
  render(<App />)
  expect(screen.getByText(/Velora Dashboard/i)).toBeInTheDocument()
})
