import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import About from './About'

describe('About', () => {
  it('renders heading', () => {
    render(<About />)
    expect(screen.getByRole('heading', { name: /about/i })).toBeInTheDocument()
  })

  it('renders description', () => {
    render(<About />)
    expect(screen.getByText(/Inventory Management System/i)).toBeInTheDocument()
  })
})
