import { createContext, useContext, useState, useEffect } from 'react'
import { auth } from '../api'

const UserContext = createContext(null)

export function UserProvider({ children }) {
  const [user, setUser] = useState(null)
  const [loading, setLoading] = useState(true)
  const [token, setTokenState] = useState(() => localStorage.getItem('token'))

  useEffect(() => {
    if (token) {
      auth.getUser()
        .then(res => setUser(res.data))
        .catch(() => {
          localStorage.removeItem('token')
          setTokenState(null)
        })
        .finally(() => setLoading(false))
    } else {
      setLoading(false)
    }
  }, [token])

  const login = (newToken, userData) => {
    localStorage.setItem('token', newToken)
    setTokenState(newToken)
    setUser(userData)
  }

  const logout = () => {
    localStorage.removeItem('token')
    setTokenState(null)
    setUser(null)
    auth.logout()
  }

  const updateUser = (userData) => {
    setUser(userData)
  }

  return (
    <UserContext.Provider value={{ user, loading, token, login, logout, updateUser }}>
      {children}
    </UserContext.Provider>
  )
}

export function useUser() {
  const context = useContext(UserContext)
  if (!context) {
    throw new Error('useUser must be used within UserProvider')
  }
  return context
}
