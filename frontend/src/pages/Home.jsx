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

  return (
    <nav className="sticky top-0 z-50 bg-white/95 backdrop-blur-md border-b border-gray-100">
      <div className="max-w-6xl mx-auto px-8 py-4 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-3">
          <Logo />
          <span className="font-serif text-xl text-[#1A1A1A]">高校图书流转</span>
        </Link>
        <div className="flex items-center gap-8">
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
      </div>
    </nav>
  )
}

function HeroSection({ mousePosition }) {
  const features = [
    { title: '教材共享', subtitle: '课程书籍循环使用', color: '#E85D2B', rotate: '-12deg', y: '20px', x: '-120px', link: '/books' },
    { title: '使用指南', subtitle: '如何发布与借阅', color: '#F5F0E8', textColor: '#1A1A1A', rotate: '-6deg', y: '0px', x: '-60px', link: '/guide' },
    { title: '智能匹配', subtitle: 'AI推荐你需要的书', color: '#2B9CD8', rotate: '0deg', y: '-15px', x: '0px', link: '/books' },
    { title: '流转方法', subtitle: '校园内的交换智慧', color: '#4ADE80', textColor: '#1A1A1A', rotate: '6deg', y: '0px', x: '60px', link: '/guide' },
    { title: '图书工具', subtitle: '分类与检索系统', color: '#1A1A1A', rotate: '12deg', y: '20px', x: '120px', link: '/books' },
  ]

  return (
    <section className="min-h-screen flex flex-col items-center justify-center px-8 pt-24 pb-16">
      <div className="text-center mb-20 max-w-2xl mx-auto">
        <p className="text-xs tracking-[0.3em] text-[#999] uppercase mb-6 animate-fade-in">
          University Book Exchange
        </p>
        <h1 className="font-serif text-7xl md:text-8xl text-[#1A1A1A] mb-8 tracking-tight leading-none animate-fade-in-up" style={{ animationDelay: '100ms' }}>
          高校图书
          <br />
          <span className="text-[#E85D2B]">流转</span>
        </h1>
        <p className="text-[#666] text-lg leading-relaxed animate-fade-in-up" style={{ animationDelay: '200ms' }}>
          让知识在校园里自由流动
        </p>
        <div className="mt-12 flex items-center justify-center gap-4 animate-fade-in-up" style={{ animationDelay: '300ms' }}>
          <Link
            to="/books"
            className="px-8 py-3 bg-[#1A1A1A] text-white rounded-full text-sm font-medium hover:bg-[#333] transition-all duration-300 hover:scale-105"
          >
            开始探索
          </Link>
          <Link
            to="/register"
            className="px-8 py-3 border border-[#1A1A1A] text-[#1A1A1A] rounded-full text-sm font-medium hover:bg-[#1A1A1A] hover:text-white transition-all duration-300"
          >
            加入我们
          </Link>
        </div>
      </div>

      <div
        className="relative w-full max-w-5xl h-[400px] flex items-center justify-center overflow-visible"
        style={{ perspective: '1400px' }}
      >
        {features.map((item, index) => (
          <Link
            key={index}
            to={item.link}
            className="absolute w-48 h-64 rounded-2xl cursor-pointer transition-all duration-500 ease-out hover:scale-110 hover:z-10"
            style={{
              backgroundColor: item.color,
              transform: `
                translateX(calc(${item.x} + ${(mousePosition.x - 0.5) * 20}px))
                translateY(calc(${item.y} + ${(mousePosition.y - 0.5) * -15}px))
                rotateY(${(mousePosition.x - 0.5) * 15}deg)
                rotateX(${(mousePosition.y - 0.5) * -10}deg)
                rotate(${item.rotate})
              `,
              zIndex: index === 2 ? 10 : 5 - Math.abs(index - 2),
              boxShadow: '0 25px 50px -12px rgba(0, 0, 0, 0.15)',
              left: '50%',
              marginLeft: `${-288 + index * 120}px`,
            }}
          >
            <div className="absolute inset-0 rounded-2xl p-5 flex flex-col justify-end">
              <h3
                className="font-serif text-xl mb-1 leading-tight"
                style={{ color: item.textColor || '#FFFFFF' }}
              >
                {item.title}
              </h3>
              <p
                className="text-xs leading-relaxed"
                style={{ color: item.textColor ? `${item.textColor}CC` : 'rgba(255,255,255,0.9)' }}
              >
                {item.subtitle}
              </p>
            </div>
          </Link>
        ))}
      </div>

      <div className="mt-20 animate-fade-in" style={{ animationDelay: '800ms' }}>
        <p className="text-xs text-[#999] tracking-widest">SCROLL TO EXPLORE</p>
        <div className="w-px h-12 bg-[#E5E5E5] mx-auto mt-2"></div>
      </div>
    </section>
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
    <section className="py-24 px-8 bg-[#FAFAFA]">
      <div className="max-w-4xl mx-auto">
        <div className="grid grid-cols-3 gap-12 text-center">
          {stats.map((stat, i) => (
            <div key={i} className="scroll-reveal" style={{ transitionDelay: `${i * 150}ms` }}>
              <div className="font-serif text-4xl text-[#1A1A1A] mb-3">{stat.number}</div>
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
    <section className="py-32 px-8">
      <div className="max-w-2xl mx-auto text-center scroll-reveal">
        <h2 className="font-serif text-4xl text-[#1A1A1A] mb-6">
          准备好分享你的图书了吗？
        </h2>
        <p className="text-[#666] mb-10 leading-relaxed">
          无论是闲置的教材还是看完的书籍，都可以在这里找到新的主人。
          加入我们，让知识流动起来。
        </p>
        <Link
          to="/publish"
          className="inline-block px-10 py-4 bg-[#E85D2B] text-white rounded-full text-sm font-medium hover:bg-[#D14E20] transition-all duration-300 hover:scale-105"
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
  const [mousePosition, setMousePosition] = useState({ x: 0.5, y: 0.5 })
  const [isLoaded, setIsLoaded] = useState(false)

  useEffect(() => {
    const handleMouseMove = (e) => {
      setMousePosition({
        x: e.clientX / window.innerWidth,
        y: e.clientY / window.innerHeight,
      })
    }

    window.addEventListener('mousemove', handleMouseMove)

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
      window.removeEventListener('mousemove', handleMouseMove)
      observer.disconnect()
    }
  }, [])

  return (
    <div className={`min-h-screen transition-all duration-700 ease-out ${isLoaded ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
      <Navbar />
      <HeroSection mousePosition={mousePosition} />
      <StatsSection />
      <CTASection />
      <Footer />
    </div>
  )
}