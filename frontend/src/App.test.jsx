import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import App from './App'

describe('App', () => {
  it('renders Home page at /', () => {
    render(
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    )
    expect(screen.getByText(/Welcome to the Inventory Management System/i)).toBeInTheDocument()
  })

  it('renders About page at /about', () => {
    render(
      <MemoryRouter initialEntries={['/about']}>
        <App />
      </MemoryRouter>
    )
    expect(screen.getByText(/Inventory Management System/i)).toBeInTheDocument()
  })
})
