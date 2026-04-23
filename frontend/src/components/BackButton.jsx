import { Link, useNavigate } from 'react-router-dom'

export default function BackButton({ to, fallback = '/' }) {
  const navigate = useNavigate()

  const handleClick = (e) => {
    e.preventDefault()
    if (window.history.length > 1) {
      navigate(-1)
    } else {
      navigate(fallback)
    }
  }

  return (
    <button
      onClick={handleClick}
      className="group flex items-center gap-2 text-sm text-[#666] hover:text-[#1A1A1A] transition-colors duration-200"
    >
      <span className="w-8 h-8 rounded-full border border-[#E5E5E5] flex items-center justify-center group-hover:border-[#1A1A1A] group-hover:bg-[#1A1A1A] group-hover:text-white transition-all duration-200">
        ←
      </span>
      <span>返回</span>
    </button>
  )
}
