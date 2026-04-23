import { createContext, useContext, useState, useCallback } from 'react'

const ToastContext = createContext(null)

let toastId = 0

export function ToastProvider({ children }) {
  const [toasts, setToasts] = useState([])

  const addToast = useCallback((message, type = 'success') => {
    const id = ++toastId
    setToasts(prev => [...prev, { id, message, type }])
    setTimeout(() => {
      setToasts(prev => prev.filter(t => t.id !== id))
    }, 3000)
  }, [])

  const toast = useCallback((message, type = 'success') => {
    addToast(message, type)
  }, [addToast])

  return (
    <ToastContext.Provider value={{ toast }}>
      {children}
      {toasts.length > 0 && (
        <div className="fixed top-20 left-1/2 -translate-x-1/2 z-[100] flex flex-col gap-2">
          {toasts.map(t => (
            <div
              key={t.id}
              className={`px-6 py-3 rounded-full shadow-lg text-sm font-medium animate-toast-in ${
                t.type === 'success'
                  ? 'bg-[#1A1A1A] text-white'
                  : t.type === 'error'
                  ? 'bg-[#E85D2B] text-white'
                  : 'bg-gray-700 text-white'
              }`}
            >
              {t.message}
            </div>
          ))}
        </div>
      )}
    </ToastContext.Provider>
  )
}

export function useToast() {
  const context = useContext(ToastContext)
  if (!context) {
    throw new Error('useToast must be used within ToastProvider')
  }
  return context
}

export const toast = () => {
  console.warn('toast must be used within ToastProvider')
}