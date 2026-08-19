import { Routes, Route } from 'react-router'
import MainLayout from '@/components/layout/MainLayout'
import Home from '@/pages/Home'
import About from '@/pages/About'

function App() {
  return (
    <Routes>
      <Route element={<MainLayout />}>
        <Route index element={<Home />} />
        <Route path="about" element={<About />} />
      </Route>
    </Routes>
  )
}

export default App
