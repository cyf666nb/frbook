import { useState, useEffect, useRef } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { useUser } from '../context/UserContext'
import { books as booksApi } from '../api'

function Logo() {
  return (
    <div className="flex items-center gap-1">
      <div className="w-4 h-4 bg-[#E85D2B] rounded-sm" />
      <div className="w-4 h-4 bg-[#2B9CD8] rounded-sm" />
      <div className="w-4 h-4 bg-[#4ADE80] rounded-sm" />
      <div className="w-4 h-4 bg-[#1A1A1A] rounded-sm" />
    </div>
  )
}

function Navbar() {
  const navigate = useNavigate()
  const { user, logout } = useUser()
  const [menuOpen, setMenuOpen] = useState(false)

  return (
    <nav className="sticky top-0 z-50 bg-white/95 backdrop-blur-md border-b border-gray-100">
      <div className="max-w-6xl mx-auto px-4 md:px-8 py-3 md:py-4 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-2 md:gap-3">
          <Logo />
          <span className="font-serif text-base md:text-xl text-[#1A1A1A]">高校图书流转</span>
        </Link>

        {/* Desktop nav */}
        <div className="hidden md:flex items-center gap-8">
          <Link to="/books" className="text-sm text-[#666] hover:text-[#1A1A1A] transition-colors">
            浏览图书
          </Link>
          <Link to="/guide" className="text-sm text-[#666] hover:text-[#1A1A1A] transition-colors">
            使用指南
          </Link>
          {user ? (
            <>
              <Link to="/orders" className="text-sm text-[#666] hover:text-[#1A1A1A] transition-colors">
                我的订单
              </Link>
              <button
                onClick={logout}
                className="text-sm text-[#666] hover:text-[#1A1A1A] transition-colors"
              >
                退出
              </button>
              <Link
                to="/profile"
                className="w-9 h-9 rounded-full bg-[#1A1A1A] text-white flex items-center justify-center text-sm font-medium hover:bg-[#333] transition-colors"
              >
                {user.nickname?.[0] || user.account?.[0] || 'U'}
              </Link>
            </>
          ) : (
            <>
              <Link
                to="/login"
                className="text-sm text-[#666] hover:text-[#1A1A1A] transition-colors"
              >
                登录
              </Link>
              <Link
                to="/register"
                className="px-5 py-2 bg-[#1A1A1A] text-white rounded-full text-sm font-medium hover:bg-[#333] transition-colors"
              >
                注册
              </Link>
            </>
          )}
        </div>

        {/* Mobile hamburger */}
        <button
          className="md:hidden p-2 -mr-2"
          onClick={() => setMenuOpen(!menuOpen)}
          aria-label="菜单"
        >
          <svg className="w-6 h-6 text-[#1A1A1A]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            {menuOpen ? (
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            ) : (
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            )}
          </svg>
        </button>
      </div>

      {/* Mobile menu drawer */}
      {menuOpen && (
        <div className="md:hidden bg-white border-t border-gray-100 px-4 py-4 space-y-3">
          <Link
            to="/books"
            className="block text-sm text-[#666] hover:text-[#1A1A1A] py-2"
            onClick={() => setMenuOpen(false)}
          >
            浏览图书
          </Link>
          <Link
            to="/guide"
            className="block text-sm text-[#666] hover:text-[#1A1A1A] py-2"
            onClick={() => setMenuOpen(false)}
          >
            使用指南
          </Link>
          {user ? (
            <>
              <Link
                to="/orders"
                className="block text-sm text-[#666] hover:text-[#1A1A1A] py-2"
                onClick={() => setMenuOpen(false)}
              >
                我的订单
              </Link>
              <Link
                to="/profile"
                className="block text-sm text-[#666] hover:text-[#1A1A1A] py-2"
                onClick={() => setMenuOpen(false)}
              >
                个人中心
              </Link>
              <button
                onClick={() => { logout(); setMenuOpen(false) }}
                className="block text-sm text-[#666] hover:text-[#1A1A1A] py-2 w-full text-left"
              >
                退出登录
              </button>
            </>
          ) : (
            <>
              <Link
                to="/login"
                className="block text-sm text-[#666] hover:text-[#1A1A1A] py-2"
                onClick={() => setMenuOpen(false)}
              >
                登录
              </Link>
              <Link
                to="/register"
                className="block text-sm font-medium text-[#1A1A1A] py-2"
                onClick={() => setMenuOpen(false)}
              >
                注册
              </Link>
            </>
          )}
        </div>
      )}
    </nav>
  )
}

function HeroSection() {
  const features = [
    {
      title: '教材共享',
      subtitle: '课程书籍循环使用',
      description: '汇聚全校教材资源，课程用书低价流转，减少书籍闲置与浪费，让知识在不同学期间延续。',
      color: '#E85D2B',
      rotate: '-12deg',
      y: '20px',
      x: '-120px',
      link: '/books',
    },
    {
      title: '使用指南',
      subtitle: '如何发布与借阅',
      description: '三分钟快速上手。发布图书只需拍照、填信息；借阅图书直接联系学长学姐，校园内当面交接。',
      color: '#F5F0E8',
      textColor: '#1A1A1A',
      rotate: '-6deg',
      y: '0px',
      x: '-60px',
      link: '/guide',
    },
    {
      title: '智能匹配',
      subtitle: 'AI推荐你需要的书',
      description: '输入课程名称或关键词，系统智能推荐相关图书，发现其他同学分享的意外宝藏。',
      color: '#2B9CD8',
      rotate: '0deg',
      y: '-15px',
      x: '0px',
      link: '/books',
    },
    {
      title: '流转方法',
      subtitle: '校园内的交换智慧',
      description: '支持免费赠送、低价出售、积分交换多种方式，灵活选择，环保又经济。',
      color: '#4ADE80',
      textColor: '#1A1A1A',
      rotate: '6deg',
      y: '0px',
      x: '60px',
      link: '/guide',
    },
    {
      title: '图书工具',
      subtitle: '分类与检索系统',
      description: '按学科、课程、学院精准分类，支持 ISBN、书名、作者多维搜索，找书快人一步。',
      color: '#1A1A1A',
      rotate: '12deg',
      y: '20px',
      x: '120px',
      link: '/books',
    },
  ]

  const [expandedIndex, setExpandedIndex] = useState(null)
  const [isAnimating, setIsAnimating] = useState(false)
  const [mousePosition, setMousePosition] = useState({ x: 0.5, y: 0.5 })
  const CARD_W = 288
  const CARD_H = 400

  const cardOffsets = [-384, -264, -144, -24, 96]

  // Cache stable transforms — recompute when not animating
  const cardTransformsRef = useRef([])
  if (!isAnimating) {
    cardTransformsRef.current = cardOffsets.map((offset, index) => {
      const mx = (mousePosition.x - 0.5) * 20
      const my = (mousePosition.y - 0.5) * -15
      return `translateX(calc(${(index - 2) * 120}px + ${mx}px)) translateY(${my}px) rotateY(${(mousePosition.x - 0.5) * 15}deg) rotateX(${(mousePosition.y - 0.5) * -10}deg) rotate(${features[index].rotate})`
    })
  }

  const getCardTransform = (index) => {
    const isExpanded = expandedIndex === index
    if (isExpanded) {
      return `translateX(calc(-144px - ${cardOffsets[index]}px)) translateY(-96px) rotate(0deg) scale(1.08)`
    }
    return cardTransformsRef.current[index]
  }

  const handleCardClick = (index) => {
    if (isAnimating) return
    setIsAnimating(true)
    if (expandedIndex === index) {
      setExpandedIndex(null)
    } else {
      setExpandedIndex(index)
    }
  }

  useEffect(() => {
    const handleMouseMove = (e) => {
      setMousePosition({
        x: e.clientX / window.innerWidth,
        y: e.clientY / window.innerHeight,
      })
    }
    window.addEventListener('mousemove', handleMouseMove)
    return () => window.removeEventListener('mousemove', handleMouseMove)
  }, [])

  // Reset isAnimating after transition completes
  useEffect(() => {
    if (isAnimating) {
      const duration = expandedIndex === null ? 250 : 500
      const timer = setTimeout(() => setIsAnimating(false), duration)
      return () => clearTimeout(timer)
    }
  }, [isAnimating])

  return (
    <>
      <section className="min-h-screen flex flex-col items-center justify-center px-4 md:px-8 pt-20 pb-12 md:pt-24 md:pb-16">
        <div className="text-center mb-20 max-w-2xl mx-auto">
          <p className="text-xs tracking-[0.3em] text-[#999] uppercase mb-6 animate-fade-in">
            University Book Exchange
          </p>
          <h1 className="font-serif text-4xl sm:text-5xl md:text-7xl lg:text-8xl text-[#1A1A1A] mb-6 md:mb-8 tracking-tight leading-none animate-fade-in-up" style={{ animationDelay: '100ms' }}>
            高校图书
            <br />
            <span className="text-[#E85D2B]">流转</span>
          </h1>
          <p className="text-[#666] text-base md:text-lg leading-relaxed animate-fade-in-up px-4 md:px-0" style={{ animationDelay: '200ms' }}>
            让知识在校园里自由流动
          </p>
          <div className="mt-8 md:mt-12 flex items-center justify-center gap-3 md:gap-4 animate-fade-in-up px-4" style={{ animationDelay: '300ms' }}>
            <Link
              to="/books"
              className="px-6 md:px-8 py-2.5 md:py-3 bg-[#1A1A1A] text-white rounded-full text-xs md:text-sm font-medium hover:bg-[#333] transition-all duration-300 hover:scale-105"
            >
              开始探索
            </Link>
            <Link
              to="/register"
              className="px-6 md:px-8 py-2.5 md:py-3 border border-[#1A1A1A] text-[#1A1A1A] rounded-full text-xs md:text-sm font-medium hover:bg-[#1A1A1A] hover:text-white transition-all duration-300"
            >
              加入我们
            </Link>
          </div>
        </div>

        {/* PC: Stacked cards with expand interaction */}
        <div
          className="relative hidden md:flex items-center justify-center"
          style={{ perspective: '1400px', height: `${CARD_H}px`, maxWidth: '900px', width: '100%', margin: '0 auto' }}
        >
          {features.map((item, index) => {
            const isExpanded = expandedIndex === index
            return (
              <div
                key={index}
                className="absolute top-0"
                style={{
                  left: '50%',
                  marginLeft: [-384, -264, -144, -24, 96][index] + 'px',
                  width: `${CARD_W}px`,
                  height: `${CARD_H}px`,
                  zIndex: isExpanded ? 50 : 5 - Math.abs(index - 2),
                  transition: isExpanded
                    ? 'all 0.5s cubic-bezier(0.34, 1.56, 0.64, 1)'
                    : 'all 0.25s cubic-bezier(0.4, 0, 0.2, 1)',
                  transform: getCardTransform(index),
                }}
              >
                {/* Card face */}
                <div
                  className="relative w-full h-full rounded-2xl cursor-pointer overflow-hidden"
                  style={{
                    backgroundColor: item.color,
                    boxShadow: isExpanded
                      ? '0 40px 80px rgba(0,0,0,0.25)'
                      : '0 25px 50px -12px rgba(0, 0, 0, 0.15)',
                  }}
                  onClick={() => handleCardClick(index)}
                >
                  <div className="absolute inset-0 rounded-2xl p-6 flex flex-col justify-between">
                    {/* Top area */}
                    <div />
                    {/* Bottom content */}
                    <div>
                      <h3
                        className="font-serif text-2xl mb-2 leading-tight"
                        style={{ color: item.textColor || '#FFFFFF' }}
                      >
                        {item.title}
                      </h3>
                      <p
                        className="text-sm leading-relaxed"
                        style={{ color: item.textColor ? `${item.textColor}CC` : 'rgba(255,255,255,0.85)' }}
                      >
                        {item.subtitle}
                      </p>
                    </div>
                  </div>

                  {/* Expanded overlay */}
                  <div
                    className="absolute inset-0 rounded-2xl flex flex-col justify-between p-8"
                    style={{
                      backgroundColor: item.color,
                      opacity: isExpanded ? 1 : 0,
                      transition: 'opacity 0.4s ease',
                      pointerEvents: isExpanded ? 'auto' : 'none',
                    }}
                  >
                    <div className="flex-1 flex flex-col justify-center">
                      <h3
                        className="font-serif text-4xl mb-4 leading-tight"
                        style={{ color: item.textColor || '#FFFFFF' }}
                      >
                        {item.title}
                      </h3>
                      <p
                        className="text-base mb-2 font-medium"
                        style={{ color: item.textColor ? `${item.textColor}CC` : 'rgba(255,255,255,0.85)' }}
                      >
                        {item.subtitle}
                      </p>
                      <p
                        className="text-sm leading-relaxed"
                        style={{ color: item.textColor ? `${item.textColor}99` : 'rgba(255,255,255,0.7)' }}
                      >
                        {item.description}
                      </p>
                    </div>
                    <div className="flex items-center justify-between mt-8">
                      <Link
                        to={item.link}
                        className="flex items-center gap-2 text-sm font-medium rounded-full px-5 py-2.5 transition-all hover:gap-3"
                        style={{
                          backgroundColor: item.textColor ? 'rgba(255,255,255,0.15)' : 'rgba(255,255,255,0.2)',
                          color: item.textColor || '#FFFFFF',
                        }}
                        onClick={(e) => e.stopPropagation()}
                      >
                        了解详情
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
                        </svg>
                      </Link>
                      <button
                        onClick={(e) => { e.stopPropagation(); setExpandedIndex(null) }}
                        className="w-9 h-9 rounded-full flex items-center justify-center"
                        style={{
                          backgroundColor: item.textColor ? 'rgba(0,0,0,0.08)' : 'rgba(0,0,0,0.15)',
                          color: item.textColor || '#FFFFFF',
                        }}
                        aria-label="关闭"
                      >
                        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            )
          })}
        </div>

        {/* Mobile cards: 2 rows - 3 top, 2 bottom, centered */}
        <div className="md:hidden flex flex-col items-center gap-3 px-4">
          {/* Row 1: 3 cards */}
          <div className="flex gap-2">
            {features.slice(0, 3).map((item, index) => (
              <Link
                key={index}
                to={item.link}
                className="w-[calc((100vw-40px)/3)] h-28 rounded-xl flex flex-col justify-end p-2.5 snap-start"
                style={{
                  backgroundColor: item.color === '#F5F0E8' ? '#FAF6F0' : item.color === '#1A1A1A' ? '#2A2A2A' : item.color,
                  boxShadow: '0 4px 12px rgba(0,0,0,0.12)',
                }}
              >
                <h3
                  className="font-serif text-xs leading-tight mb-0.5"
                  style={{ color: item.color === '#F5F0E8' ? '#1A1A1A' : item.textColor || '#FFFFFF' }}
                >
                  {item.title}
                </h3>
                <p
                  className="text-[10px] leading-relaxed"
                  style={{ color: item.color === '#F5F0E8' ? '#666' : item.textColor ? `${item.textColor}99` : 'rgba(255,255,255,0.75)' }}
                >
                  {item.subtitle}
                </p>
              </Link>
            ))}
          </div>
          {/* Row 2: 2 cards, centered */}
          <div className="flex gap-2 justify-center">
            {/* Left spacer */}
            <div className="w-[calc((100vw-40px)/6)]" />
            {features.slice(3).map((item, index) => (
              <Link
                key={index + 3}
                to={item.link}
                className="w-[calc((100vw-40px)/3)] h-28 rounded-xl flex flex-col justify-end p-2.5 snap-start"
                style={{
                  backgroundColor: item.color === '#F5F0E8' ? '#FAF6F0' : item.color === '#1A1A1A' ? '#2A2A2A' : item.color,
                  boxShadow: '0 4px 12px rgba(0,0,0,0.12)',
                }}
              >
                <h3
                  className="font-serif text-xs leading-tight mb-0.5"
                  style={{ color: item.color === '#F5F0E8' ? '#1A1A1A' : item.textColor || '#FFFFFF' }}
                >
                  {item.title}
                </h3>
                <p
                  className="text-[10px] leading-relaxed"
                  style={{ color: item.color === '#F5F0E8' ? '#666' : item.textColor ? `${item.textColor}99` : 'rgba(255,255,255,0.75)' }}
                >
                  {item.subtitle}
                </p>
              </Link>
            ))}
            {/* Right spacer */}
            <div className="w-[calc((100vw-40px)/6)]" />
          </div>
        </div>

        <div className="mt-20 animate-fade-in" style={{ animationDelay: '800ms' }}>
          <p className="text-xs text-[#999] tracking-widest">SCROLL TO EXPLORE</p>
          <div className="w-px h-12 bg-[#E5E5E5] mx-auto mt-2" />
        </div>
      </section>
    </>
  )
}

function StatsSection() {
  const [bookCount, setBookCount] = useState(0)

  useEffect(() => {
    async function fetchStats() {
      try {
        const res = await booksApi.list({ page: 1, page_size: 1 })
        if (res.data?.total !== undefined) {
          setBookCount(res.data.total)
        }
      } catch (err) {
        console.error('Failed to fetch book count:', err)
      }
    }
    fetchStats()
  }, [])

  const formatNumber = (num) => {
    if (num >= 10000) {
      return (num / 10000).toFixed(1) + 'W+'
    }
    if (num >= 1000) {
      return (num / 1000).toFixed(1) + 'K+'
    }
    return num.toLocaleString()
  }

  const stats = [
    { number: formatNumber(bookCount) || '10,000+', label: '注册用户' },
    { number: formatNumber(bookCount) || '50,000+', label: '流转图书' },
    { number: '101%', label: '用户满意度' },
  ]

  return (
    <section className="py-12 md:py-24 px-4 md:px-8 bg-[#FAFAFA]">
      <div className="max-w-4xl mx-auto">
        <div className="grid grid-cols-3 gap-4 md:gap-12 text-center">
          {stats.map((stat, i) => (
            <div key={i} className="scroll-reveal" style={{ transitionDelay: `${i * 150}ms` }}>
              <div className="font-serif text-2xl md:text-4xl text-[#1A1A1A] mb-1 md:mb-3">{stat.number}</div>
              <div className="text-xs text-[#666] tracking-widest uppercase">{stat.label}</div>
            </div>
          ))}
        </div>
      </div>
    </section>
  )
}

function CTASection() {
  return (
    <section className="py-16 md:py-32 px-4 md:px-8">
      <div className="max-w-2xl mx-auto text-center scroll-reveal">
        <h2 className="font-serif text-2xl md:text-4xl text-[#1A1A1A] mb-4 md:mb-6">
          准备好分享你的图书了吗？
        </h2>
        <p className="text-[#666] mb-8 md:mb-10 leading-relaxed text-sm md:text-base px-4 md:px-0">
          无论是闲置的教材还是看完的书籍，都可以在这里找到新的主人。
          加入我们，让知识流动起来。
        </p>
        <Link
          to="/publish"
          className="inline-block px-8 md:px-10 py-3 md:py-4 bg-[#E85D2B] text-white rounded-full text-xs md:text-sm font-medium hover:bg-[#D14E20] transition-all duration-300 hover:scale-105"
        >
          发布我的图书
        </Link>
      </div>
    </section>
  )
}

function Footer() {
  return (
    <footer className="py-16 px-8 border-t border-[#F0F0F0]">
      <div className="max-w-6xl mx-auto flex flex-col md:flex-row items-center justify-between gap-6">
        <div className="flex items-center gap-3">
          <Logo />
          <span className="font-serif text-lg text-[#1A1A1A]">高校图书流转</span>
        </div>
        <p className="text-xs text-[#999] tracking-widest">
          © 2026 UNIVERSITY BOOK EXCHANGE
        </p>
      </div>
    </footer>
  )
}

export default function Home() {
  const [isLoaded, setIsLoaded] = useState(false)

  useEffect(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        entries.forEach((entry) => {
          if (entry.isIntersecting) {
            entry.target.classList.add('revealed')
          }
        })
      },
      { threshold: 0.1 }
    )

    document.querySelectorAll('.scroll-reveal').forEach((el) => {
      observer.observe(el)
    })

    setIsLoaded(true)

    return () => {
      observer.disconnect()
    }
  }, [])

  return (
    <div className={`min-h-screen transition-all duration-700 ease-out ${isLoaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
      <Navbar />
      <HeroSection />
      <StatsSection />
      <CTASection />
      <Footer />
    </div>
  )
}