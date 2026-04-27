import { useState, useEffect } from 'react'
import { Link } from 'react-router-dom'
import { guide } from '../api'
import BackButton from '../components/BackButton'

const iconMap = {
  book: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
    </svg>
  ),
  rent: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
    </svg>
  ),
  buy: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 3h2l.4 2M7 13h10l4-8H5.4M7 13L5.4 5M7 13l-2.293 2.293c-.63.63-.184 1.707.707 1.707H17m0 0a2 2 0 100 4 2 2 0 000-4zm-8 2a2 2 0 11-4 0 2 2 0 014 0z" />
    </svg>
  ),
  gift: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7" />
    </svg>
  ),
  tips: (
    <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z" />
    </svg>
  ),
}

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
  return (
    <nav className="sticky top-0 z-50 bg-white/95 backdrop-blur-md border-b border-gray-100">
      <div className="max-w-6xl mx-auto px-4 md:px-8 py-4 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-3">
          <Logo />
          <span className="font-serif text-xl text-[#1A1A1A]">高校图书流转</span>
        </Link>
        <Link to="/books" className="text-sm text-[#666] hover:text-[#1A1A1A] transition-colors">
          浏览图书
        </Link>
      </div>
    </nav>
  )
}

function SectionCard({ section, index }) {
  const [expanded, setExpanded] = useState(false)

  return (
    <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full px-6 py-5 flex items-center gap-4 text-left hover:bg-gray-50 transition-colors"
      >
        <div className="w-12 h-12 rounded-xl bg-[#F5F5F5] flex items-center justify-center text-[#1A1A1A]">
          {iconMap[section.icon] || iconMap.book}
        </div>
        <div className="flex-1">
          <h3 className="font-serif text-lg text-[#1A1A1A] mb-1">{section.title}</h3>
          <p className="text-sm text-[#666]">{section.description}</p>
        </div>
        <svg
          className={`w-5 h-5 text-[#999] transition-transform duration-200 ${expanded ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>

      {expanded && (
        <div className="px-6 pb-6 border-t border-gray-100">
          {section.steps && (
            <div className="pt-6 space-y-4">
              {section.steps.map((step, i) => (
                <div key={i} className="flex gap-4">
                  <div className="w-8 h-8 rounded-full bg-[#1A1A1A] text-white flex items-center justify-center text-sm font-medium flex-shrink-0">
                    {step.step}
                  </div>
                  <div>
                    <h4 className="font-medium text-[#1A1A1A] mb-1">{step.title}</h4>
                    <p className="text-sm text-[#666]">{step.content}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
          {section.items && (
            <div className="pt-6 grid gap-4">
              {section.items.map((item, i) => (
                <div key={i} className="flex gap-3 items-start">
                  <div className="w-2 h-2 rounded-full bg-[#E85D2B] mt-2 flex-shrink-0" />
                  <div>
                    <h4 className="font-medium text-[#1A1A1A] mb-1">{item.title}</h4>
                    <p className="text-sm text-[#666]">{item.content}</p>
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>
      )}
    </div>
  )
}

function FAQItem({ faq, index }) {
  const [expanded, setExpanded] = useState(false)

  return (
    <div className="border-b border-gray-100 last:border-0">
      <button
        onClick={() => setExpanded(!expanded)}
        className="w-full py-5 flex items-center justify-between text-left hover:bg-gray-50 transition-colors px-2"
      >
        <span className="font-medium text-[#1A1A1A]">{faq.question}</span>
        <svg
          className={`w-5 h-5 text-[#999] transition-transform duration-200 ${expanded ? 'rotate-180' : ''}`}
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
        </svg>
      </button>
      {expanded && (
        <div className="pb-5 px-2">
          <p className="text-[#666] leading-relaxed">{faq.answer}</p>
        </div>
      )}
    </div>
  )
}

export default function Guide() {
  const [guideData, setGuideData] = useState(null)
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    async function fetchGuide() {
      try {
        const res = await guide.get()
        setGuideData(res.data)
      } catch (err) {
        console.error('Failed to fetch guide:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchGuide()
  }, [])

  return (
    <div className="min-h-screen bg-[#FAFAFA]">
      <Navbar />

      <div className="max-w-6xl mx-auto px-4 md:px-8 py-8">
        <BackButton />

        <div className="mt-8 mb-12">
          <h1 className="font-serif text-4xl text-[#1A1A1A] mb-3">使用指南</h1>
          <p className="text-[#666]">了解如何使用高校图书流转平台</p>
        </div>

        {loading ? (
          <div className="space-y-4">
            {[...Array(5)].map((_, i) => (
              <div key={i} className="bg-white rounded-2xl h-24 animate-pulse" />
            ))}
          </div>
        ) : guideData ? (
          <>
            <div className="space-y-4 mb-12">
              {guideData.sections?.map((section, i) => (
                <SectionCard key={section.id || i} section={section} index={i} />
              ))}
            </div>

            {guideData.faq && guideData.faq.length > 0 && (
              <div className="bg-white rounded-2xl shadow-sm border border-gray-100 overflow-hidden">
                <div className="px-6 py-5 border-b border-gray-100">
                  <h2 className="font-serif text-xl text-[#1A1A1A]">常见问题</h2>
                </div>
                <div className="px-4">
                  {guideData.faq.map((faq, i) => (
                    <FAQItem key={i} faq={faq} index={i} />
                  ))}
                </div>
              </div>
            )}

            <div className="mt-12 text-center">
              <p className="text-[#666] mb-4">还有其他问题？</p>
              <Link
                to="/books"
                className="inline-flex items-center gap-2 px-6 py-3 bg-[#1A1A1A] text-white rounded-full text-sm font-medium hover:bg-[#333] transition-colors"
              >
                开始探索图书
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
                </svg>
              </Link>
            </div>
          </>
        ) : (
          <div className="text-center py-12">
            <p className="text-[#666]">暂无使用指南内容</p>
          </div>
        )}
      </div>

      <footer className="py-8 md:py-12 px-4 md:px-8 border-t border-[#F0F0F0] bg-white mt-12">
        <div className="max-w-6xl mx-auto flex items-center justify-between">
          <div className="flex items-center gap-3">
            <Logo />
            <span className="font-serif text-[#1A1A1A]">高校图书流转</span>
          </div>
          <p className="text-xs text-[#999]">© 2026 UNIVERSITY BOOK EXCHANGE</p>
        </div>
      </footer>
    </div>
  )
}
