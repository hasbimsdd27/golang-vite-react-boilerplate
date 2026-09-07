import { lazy, Suspense } from 'react'
import { Routes, Route } from 'react-router'
import MainLayout from '@/components/layout/MainLayout'

const Home = lazy(() => import('@/pages/Home'))
const About = lazy(() => import('@/pages/About'))

function App() {
  return (
    <Suspense fallback={<div className="p-4">Loading...</div>}>
      <Routes>
        <Route element={<MainLayout />}>
          <Route index element={<Home />} />
          <Route path="about" element={<About />} />
        </Route>
      </Routes>
    </Suspense>
  )
}

export default App
