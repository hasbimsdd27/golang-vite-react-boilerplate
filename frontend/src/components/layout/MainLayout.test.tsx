import { describe, it, expect } from 'vitest'
import { render, screen } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import MainLayout from './MainLayout'

describe('MainLayout', () => {
  it('renders navigation links', () => {
    render(
      <MemoryRouter>
        <MainLayout />
      </MemoryRouter>
    )
    expect(screen.getByRole('link', { name: /home/i })).toBeInTheDocument()
    expect(screen.getByRole('link', { name: /about/i })).toBeInTheDocument()
  })

  it('renders children via Outlet', () => {
    render(
      <MemoryRouter>
        <MainLayout />
      </MemoryRouter>
    )
    expect(screen.getByRole('main')).toBeInTheDocument()
  })
})
