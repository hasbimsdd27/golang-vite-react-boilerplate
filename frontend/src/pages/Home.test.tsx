import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import Home from './Home'

describe('Home', () => {
  it('renders heading', () => {
    render(<Home />)
    expect(screen.getByRole('heading', { name: /home/i })).toBeInTheDocument()
  })

  it('renders welcome message', () => {
    render(<Home />)
    expect(screen.getByText(/Welcome to the Inventory Management System/i)).toBeInTheDocument()
  })
})
