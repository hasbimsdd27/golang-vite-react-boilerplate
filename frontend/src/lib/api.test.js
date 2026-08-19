import { describe, it, expect } from 'vitest'
import api from './api'

describe('api', () => {
  it('is an axios instance', () => {
    expect(api).toBeDefined()
    expect(api.defaults).toBeDefined()
    expect(typeof api.get).toBe('function')
    expect(typeof api.post).toBe('function')
  })

  it('has baseURL configured', () => {
    expect(api.defaults.baseURL).toBe('http://localhost:3000/api')
  })

  it('includes credentials in requests', () => {
    expect(api.defaults.withCredentials).toBe(true)
  })
})
