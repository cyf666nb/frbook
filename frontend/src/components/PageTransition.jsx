import { useState, useEffect, useRef } from 'react'
import { useLocation } from 'react-router-dom'

export default function PageTransition({ children }) {
  const location = useLocation()
  const [displayChildren, setDisplayChildren] = useState(children)
  const [visible, setVisible] = useState(true)
  const timeoutRef = useRef(null)
  const prevPathnameRef = useRef(location.pathname)

  useEffect(() => {
    if (location.pathname !== prevPathnameRef.current) {
      setVisible(false)

      timeoutRef.current = setTimeout(() => {
        setDisplayChildren(children)
        prevPathnameRef.current = location.pathname
        requestAnimationFrame(() => {
          setVisible(true)
        })
      }, 300)
    }

    return () => {
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current)
      }
    }
  }, [location.pathname, children])

  return (
    <div
      className={`transition-all duration-500 ease-out ${
        visible ? 'opacity-100 translate-y-0 scale-100' : 'opacity-0 translate-y-4 scale-[0.98]'
      }`}
      style={{
        transitionTimingFunction: 'cubic-bezier(0.22, 1, 0.36, 1)',
      }}
    >
      {displayChildren}
    </div>
  )
}