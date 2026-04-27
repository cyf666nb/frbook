import { Routes, Route, useLocation } from 'react-router-dom'
import Home from './pages/Home'
import Books from './pages/Books'
import Login from './pages/auth/Login'
import Register from './pages/auth/Register'
import BookDetail from './pages/BookDetail'
import Publish from './pages/Publish'
import Profile from './pages/Profile'
import Orders from './pages/Orders'
import Guide from './pages/Guide'
import PageTransition from './components/PageTransition'
import { ToastProvider } from './components/Toast'

function AnimatedRoutes() {
  const location = useLocation()

  return (
    <Routes location={location} key={location.pathname}>
      <Route path="/" element={<PageTransition><Home /></PageTransition>} />
      <Route path="/books" element={<PageTransition><Books /></PageTransition>} />
      <Route path="/login" element={<PageTransition><Login /></PageTransition>} />
      <Route path="/register" element={<PageTransition><Register /></PageTransition>} />
      <Route path="/book/:id" element={<PageTransition><BookDetail /></PageTransition>} />
      <Route path="/publish" element={<PageTransition><Publish /></PageTransition>} />
      <Route path="/profile" element={<PageTransition><Profile /></PageTransition>} />
      <Route path="/orders" element={<PageTransition><Orders /></PageTransition>} />
      <Route path="/guide" element={<PageTransition><Guide /></PageTransition>} />
      <Route path="/my-books" element={<PageTransition><Profile /></PageTransition>} />
    </Routes>
  )
}

export default function App() {
  return (
    <ToastProvider>
      <AnimatedRoutes />
    </ToastProvider>
  )
}