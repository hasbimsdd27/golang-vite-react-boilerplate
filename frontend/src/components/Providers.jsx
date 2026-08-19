import { BrowserRouter } from 'react-router'
import { QueryClientProvider } from '@tanstack/react-query'
import queryClient from '@/lib/queryClient'

export default function Providers({ children }) {
  return (
    <BrowserRouter>
      <QueryClientProvider client={queryClient}>
        {children}
      </QueryClientProvider>
    </BrowserRouter>
  )
}
