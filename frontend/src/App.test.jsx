import { describe, it, expect } from 'vitest'
import { render, screen, waitFor } from '@testing-library/react'
import { MemoryRouter } from 'react-router'
import { Suspense } from 'react'
import App from './App'

function renderWithSuspense(ui) {
  return render(
    <Suspense fallback={<div>Loading...</div>}>
      {ui}
    </Suspense>
  )
}

describe('App', () => {
  it('renders Home page at /', async () => {
    renderWithSuspense(
      <MemoryRouter initialEntries={['/']}>
        <App />
      </MemoryRouter>
    )
    await waitFor(() => {
      expect(screen.getByText(/Welcome to the Inventory Management System/i)).toBeInTheDocument()
    })
  })

  it('renders About page at /about', async () => {
    renderWithSuspense(
      <MemoryRouter initialEntries={['/about']}>
        <App />
      </MemoryRouter>
    )
    await waitFor(() => {
      expect(screen.getByText(/Inventory Management System/i)).toBeInTheDocument()
    })
  })
})
