import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App'
import Providers from './components/Providers'

const rootElement = document.getElementById('root')
if (!rootElement) {
  throw new Error('Failed to find root element')
}

createRoot(rootElement).render(
  <StrictMode>
    <Providers>
      <App />
    </Providers>
  </StrictMode>,
)
