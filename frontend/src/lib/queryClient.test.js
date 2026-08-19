import { describe, it, expect } from 'vitest'
import queryClient from './queryClient'

describe('queryClient', () => {
  it('is a QueryClient instance', () => {
    expect(queryClient).toBeDefined()
    expect(typeof queryClient.getQueryData).toBe('function')
    expect(typeof queryClient.setQueryData).toBe('function')
  })

  it('has default options configured', () => {
    expect(queryClient.getDefaultOptions()).toBeDefined()
  })
})
