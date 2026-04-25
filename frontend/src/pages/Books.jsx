import { useState, useEffect } from 'react'
import { Link, useSearchParams } from 'react-router-dom'
import { books } from '../api'
import { useUser } from '../context/UserContext'
import BackButton from '../components/BackButton'

const modeConfig = {
  1: { name: '租赁', color: '#E85D2B', bg: 'bg-[#E85D2B]' },
  2: { name: '出售', color: '#2B9CD8', bg: 'bg-[#2B9CD8]' },
  3: { name: '赠送', color: '#4ADE80', bg: 'bg-[#4ADE80]' },
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
  const { user } = useUser()

  return (
    <nav className="sticky top-0 z-50 bg-white/95 backdrop-blur-md border-b border-gray-100">
      <div className="max-w-6xl mx-auto px-8 py-4 flex items-center justify-between">
        <Link to="/" className="flex items-center gap-3">
          <Logo />
          <span className="font-serif text-xl text-[#1A1A1A]">高校图书流转</span>
        </Link>
        <div className="flex items-center gap-6">
          <Link to="/books" className="text-sm text-[#1A1A1A] font-medium">浏览</Link>
          {user ? (
            <>
              <Link to="/publish" className="text-sm text-[#666] hover:text-[#1A1A1A]">发布</Link>
              <Link to="/orders" className="text-sm text-[#666] hover:text-[#1A1A1A]">订单</Link>
              <Link to="/profile" className="text-sm text-[#1A1A1A] font-medium">{user.nickname || user.account}</Link>
            </>
          ) : (
            <Link to="/login" className="px-5 py-2 bg-[#F5F5F5] rounded-full text-sm">登录</Link>
          )}
        </div>
      </div>
    </nav>
  )
}

function FilterBar({ filters, setFilters, total }) {
  const [showMobileFilters, setShowMobileFilters] = useState(false)
  const [showCategoryDropdown, setShowCategoryDropdown] = useState(false)
  const [showSortDropdown, setShowSortDropdown] = useState(false)
  const [searchInput, setSearchInput] = useState(filters.keyword || '')

  const categoryOptions = ['教材', '文学', '科技', '艺术', '历史', '其他']
  const sortOptions = [
    { value: 'created_at-desc', label: '最新发布' },
    { value: 'view_count-desc', label: '热度最高' },
    { value: 'price-asc', label: '价格从低到高' },
    { value: 'price-desc', label: '价格从高到低' },
  ]

  const currentSortLabel = sortOptions.find(o => o.value === `${filters.sort}-${filters.order}`)?.label || '排序'
  const currentCategoryLabel = filters.category || '分类'

  const handleSearch = (e) => {
    e.preventDefault()
    setFilters(f => ({ ...f, keyword: searchInput || null }))
  }

  const clearSearch = () => {
    setSearchInput('')
    setFilters(f => ({ ...f, keyword: null }))
  }

  useEffect(() => {
    if (showMobileFilters) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
    return () => { document.body.style.overflow = '' }
  }, [showMobileFilters])

  const activeFilterCount = [
    filters.mode !== null,
    filters.category !== null,
    filters.sort !== 'created_at' || filters.order !== 'desc',
  ].filter(Boolean).length

  return (
    <>
      {/* Mobile filter overlay — fixed fullscreen, z-index max, body scroll locked */}
      {showMobileFilters && (
        <div className="fixed inset-0 z-[9999] flex flex-col">
          <div className="absolute inset-0 bg-black/50" onClick={() => setShowMobileFilters(false)} />
          <div
            className="relative w-full max-w-lg mx-auto bg-white overflow-y-auto"
            style={{ marginTop: '20vh', maxHeight: '80vh' }}
          >
            {/* Header */}
            <div className="sticky top-0 bg-white z-10 px-6 py-4 flex items-center justify-between border-b border-gray-100">
              <h3 className="font-serif text-xl text-[#1A1A1A]">筛选</h3>
              <button onClick={() => setShowMobileFilters(false)} className="p-2">
                <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            {/* Content */}
            <div className="p-6">
              {/* Mode filters */}
              <div className="mb-6">
                <p className="text-sm font-medium text-[#1A1A1A] mb-3">流转方式</p>
                <div className="flex flex-wrap gap-2">
                  <button
                    onClick={() => setFilters(f => ({ ...f, mode: null }))}
                    className={`px-4 py-2 rounded-full text-sm font-medium transition-all ${
                      filters.mode === null ? 'bg-[#1A1A1A] text-white' : 'bg-[#F5F5F5] text-[#666]'
                    }`}
                  >
                    全部
                  </button>
                  {Object.entries(modeConfig).map(([key, val]) => (
                    <button
                      key={key}
                      onClick={() => setFilters(f => ({ ...f, mode: Number(key) }))}
                      className={`px-4 py-2 rounded-full text-sm font-medium transition-all ${
                        filters.mode === Number(key) ? 'text-white' : 'bg-[#F5F5F5] text-[#666]'
                      }`}
                      style={filters.mode === Number(key) ? { backgroundColor: val.color } : {}}
                    >
                      {val.name}
                    </button>
                  ))}
                </div>
              </div>

              {/* Category filters */}
              <div className="mb-6">
                <p className="text-sm font-medium text-[#1A1A1A] mb-3">分类</p>
                <div className="flex flex-wrap gap-2">
                  <button
                    onClick={() => setFilters(f => ({ ...f, category: null }))}
                    className={`px-4 py-2 rounded-full text-sm transition-all ${
                      !filters.category ? 'bg-[#1A1A1A] text-white' : 'bg-[#F5F5F5] text-[#666]'
                    }`}
                  >
                    全部
                  </button>
                  {categoryOptions.map(cat => (
                    <button
                      key={cat}
                      onClick={() => setFilters(f => ({ ...f, category: cat }))}
                      className={`px-4 py-2 rounded-full text-sm transition-all ${
                        filters.category === cat ? 'bg-[#1A1A1A] text-white' : 'bg-[#F5F5F5] text-[#666]'
                      }`}
                    >
                      {cat}
                    </button>
                  ))}
                </div>
              </div>

              {/* Sort filters */}
              <div className="mb-6">
                <p className="text-sm font-medium text-[#1A1A1A] mb-3">排序</p>
                <div className="space-y-2">
                  {sortOptions.map(option => {
                    const isActive = filters.sort === option.value.split('-')[0] && filters.order === option.value.split('-')[1]
                    return (
                      <button
                        key={option.value}
                        onClick={() => {
                          const [sort, order] = option.value.split('-')
                          setFilters(f => ({ ...f, sort, order }))
                        }}
                        className={`w-full px-4 py-3 rounded-xl text-sm text-left transition-all ${
                          isActive ? 'bg-[#1A1A1A] text-white' : 'bg-[#F5F5F5] text-[#666]'
                        }`}
                      >
                        {option.label}
                      </button>
                    )
                  })}
                </div>
              </div>

              {/* Actions */}
              <div className="flex gap-3 pt-2">
                <button
                  onClick={() => {
                    setFilters({ mode: null, category: null, keyword: null, sort: 'created_at', order: 'desc' })
                    setSearchInput('')
                  }}
                  className="flex-1 py-3 bg-[#F5F5F5] rounded-xl text-sm text-[#666]"
                >
                  清除全部
                </button>
                <button
                  onClick={() => setShowMobileFilters(false)}
                  className="flex-1 py-3 bg-[#1A1A1A] text-white rounded-xl text-sm font-medium"
                >
                  应用 ({total} 本)
                </button>
              </div>
            </div>
          </div>
        </div>
      )}

      {/* Filter bar — sticky, z-index below modal */}
      <div className="sticky top-[48px] z-40 bg-white/60 backdrop-blur-md border-b border-gray-100 py-2 md:py-4 shadow-sm">
        <div className="max-w-6xl mx-auto px-4 md:px-8">
          {/* Mobile: search + filter button in one row */}
          <div className="flex items-center gap-2 md:hidden">
            <BackButton />
            <form onSubmit={handleSearch} className="flex-1">
              <div className="relative">
                <input
                  type="text"
                  value={searchInput}
                  onChange={(e) => setSearchInput(e.target.value)}
                  placeholder="搜索..."
                  className="w-full pl-8 pr-7 py-1.5 bg-white border border-gray-200 rounded-xl text-sm outline-none placeholder:text-[#999] transition-all shadow-sm"
                />
                <svg className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#999]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
                {searchInput && (
                  <button type="button" onClick={clearSearch} className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#999]">
                    <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                )}
              </div>
            </form>
            <button
              onClick={() => setShowMobileFilters(true)}
              className="relative flex items-center gap-1 px-2.5 py-1.5 bg-white border border-gray-200 rounded-xl text-sm text-[#1A1A1A] transition-all shadow-sm"
            >
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 4a1 1 0 011-1h16a1 1 0 011 1v2.586a1 1 0 01-.293.707l-6.414 6.414a1 1 0 00-.293.707V17l-4 4v-6.586a1 1 0 00-.293-.707L3.293 7.293A1 1 0 013 6.586V4z" />
              </svg>
              {activeFilterCount > 0 && (
                <span className="absolute -top-1 -right-1 w-4 h-4 bg-[#E85D2B] text-white text-[10px] rounded-full flex items-center justify-center">
                  {activeFilterCount}
                </span>
              )}
            </button>
          </div>

          {/* Active filter chips (mobile) */}
          {(filters.mode !== null || filters.category) && (
            <div className="flex items-center gap-2 mt-2 md:hidden overflow-x-auto">
              {filters.mode !== null && (
                <span className="flex items-center gap-1 px-3 py-1 bg-[#F5F5F5] rounded-full text-xs text-[#666] whitespace-nowrap">
                  {modeConfig[filters.mode].name}
                  <button onClick={() => setFilters(f => ({ ...f, mode: null }))} className="ml-1">×</button>
                </span>
              )}
              {filters.category && (
                <span className="flex items-center gap-1 px-3 py-1 bg-[#F5F5F5] rounded-full text-xs text-[#666] whitespace-nowrap">
                  {filters.category}
                  <button onClick={() => setFilters(f => ({ ...f, category: null }))} className="ml-1">×</button>
                </span>
              )}
            </div>
          )}

          {/* Desktop filter bar */}
          <div className="hidden md:flex items-center justify-between gap-6">
            <div className="flex items-center gap-6 flex-1">
              <BackButton />

              <form onSubmit={handleSearch} className="flex-1 max-w-md">
                <div className="relative">
                  <input
                    type="text"
                    value={searchInput}
                    onChange={(e) => setSearchInput(e.target.value)}
                    placeholder="搜索书名或作者..."
                    className="w-full pl-10 pr-10 py-2 bg-[#F5F5F5] rounded-full text-sm outline-none focus:ring-2 focus:ring-[#1A1A1A]/10 transition-all"
                  />
                  <svg className="absolute left-3.5 top-1/2 -translate-y-1/2 w-4 h-4 text-[#999]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                  {searchInput && (
                    <button
                      type="button"
                      onClick={clearSearch}
                      className="absolute right-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[#999] hover:text-[#666]"
                    >
                      <svg fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  )}
                </div>
              </form>

              <div className="flex items-center gap-3">
                <div className="flex gap-2">
                  <button
                    onClick={() => setFilters(f => ({ ...f, mode: null }))}
                    className={`px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${
                      filters.mode === null
                        ? 'bg-[#1A1A1A] text-white shadow-sm'
                        : 'bg-[#F5F5F5] text-[#666] hover:bg-[#EFEFEF]'
                    }`}
                  >
                    全部
                  </button>
                  {Object.entries(modeConfig).map(([key, val]) => (
                    <button
                      key={key}
                      onClick={() => setFilters(f => ({ ...f, mode: Number(key) }))}
                      className={`px-4 py-2 rounded-full text-sm font-medium transition-all duration-200 ${
                        filters.mode === Number(key)
                          ? 'text-white shadow-sm'
                          : 'bg-[#F5F5F5] text-[#666] hover:bg-[#EFEFEF]'
                      }`}
                      style={filters.mode === Number(key) ? { backgroundColor: val.color } : {}}
                    >
                      {val.name}
                    </button>
                  ))}
                </div>
              </div>

              <div className="relative">
                <button
                  onClick={() => {
                    setShowCategoryDropdown(!showCategoryDropdown)
                    setShowSortDropdown(false)
                  }}
                  className={`flex items-center gap-2 px-4 py-2 rounded-full text-sm transition-all duration-200 ${
                    filters.category
                      ? 'bg-[#1A1A1A] text-white'
                      : 'bg-[#F5F5F5] text-[#666] hover:bg-[#EFEFEF]'
                  }`}
                >
                  <span>{currentCategoryLabel}</span>
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                {showCategoryDropdown && (
                  <div className="absolute top-full left-0 mt-2 bg-white rounded-xl shadow-lg border border-gray-100 py-2 min-w-[120px] z-50">
                    <button
                      onClick={() => {
                        setFilters(f => ({ ...f, category: null }))
                        setShowCategoryDropdown(false)
                      }}
                      className="w-full px-4 py-2 text-left text-sm text-[#666] hover:bg-gray-50 transition-colors"
                    >
                      全部分类
                    </button>
                    {categoryOptions.map(cat => (
                      <button
                        key={cat}
                        onClick={() => {
                          setFilters(f => ({ ...f, category: cat }))
                          setShowCategoryDropdown(false)
                        }}
                        className="w-full px-4 py-2 text-left text-sm text-[#666] hover:bg-gray-50 transition-colors"
                      >
                        {cat}
                      </button>
                    ))}
                  </div>
                )}
              </div>

              <div className="relative">
                <button
                  onClick={() => {
                    setShowSortDropdown(!showSortDropdown)
                    setShowCategoryDropdown(false)
                  }}
                  className="flex items-center gap-2 px-4 py-2 rounded-full text-sm bg-[#F5F5F5] text-[#666] hover:bg-[#EFEFEF] transition-all duration-200"
                >
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
                  </svg>
                  <span>{currentSortLabel}</span>
                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                {showSortDropdown && (
                  <div className="absolute top-full left-0 mt-2 bg-white rounded-xl shadow-lg border border-gray-100 py-2 min-w-[140px] z-50">
                    {sortOptions.map(option => (
                      <button
                        key={option.value}
                        onClick={() => {
                          const [sort, order] = option.value.split('-')
                          setFilters(f => ({ ...f, sort, order }))
                          setShowSortDropdown(false)
                        }}
                        className="w-full px-4 py-2 text-left text-sm text-[#666] hover:bg-gray-50 transition-colors"
                      >
                        {option.label}
                      </button>
                    ))}
                  </div>
                )}
              </div>
            </div>

            <span className="text-sm text-[#999] whitespace-nowrap">{total} 本图书</span>
          </div>
        </div>
      </div>
    </>
  )
}

function BookCard({ book }) {
  const mode = modeConfig[book.mode] || { name: '未知', color: '#666' }

  return (
    <Link 
      to={`/book/${book.id}`}
      className="group block"
    >
      <div className="aspect-[3/4] bg-[#F8F8F8] rounded-xl overflow-hidden relative mb-4 transition-all duration-300 group-hover:shadow-lg">
        {book.cover_image ? (
          <img 
            src={book.cover_image} 
            alt={book.title}
            className="w-full h-full object-cover transition-transform duration-500 group-hover:scale-105"
          />
        ) : (
          <div className="w-full h-full flex items-center justify-center">
            <svg className="w-12 h-12 text-[#DDD]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
          </div>
        )}
        <div 
          className="absolute top-3 left-3 px-2.5 py-1 rounded-full text-xs text-white font-medium"
          style={{ backgroundColor: mode.color }}
        >
          {mode.name}
        </div>
      </div>
      <h3 className="font-serif text-base text-[#1A1A1A] mb-1 line-clamp-2 group-hover:text-[#E85D2B] transition-colors duration-200">
        {book.title}
      </h3>
      <p className="text-xs text-[#999] mb-2 line-clamp-1">{book.author || '未知作者'}</p>
      <div className="flex items-center justify-between">
        <span className="text-xs text-[#BBB]">{book.school || '未填写学校'}</span>
        {book.mode === 1 && book.daily_rent && (
          <span className="text-sm font-medium text-[#E85D2B]">¥{book.daily_rent}/天</span>
        )}
        {book.mode === 2 && book.sell_price && (
          <span className="text-sm font-medium text-[#2B9CD8]">¥{book.sell_price}</span>
        )}
        {book.mode === 3 && (
          <span className="text-sm font-medium text-[#4ADE80]">免费</span>
        )}
      </div>
    </Link>
  )
}

function Pagination({ page, totalPages, onPageChange }) {
  return (
    <div className="flex items-center justify-center gap-2 py-12">
      <button
        onClick={() => onPageChange(page - 1)}
        disabled={page <= 1}
        className="w-10 h-10 rounded-full border border-[#E5E5E5] text-sm text-[#666] hover:bg-[#F5F5F5] disabled:opacity-30 disabled:cursor-not-allowed transition-all"
      >
        ←
      </button>
      
      <div className="flex items-center gap-1">
        {Array.from({ length: Math.min(5, totalPages) }, (_, i) => {
          let pageNum
          if (totalPages <= 5) {
            pageNum = i + 1
          } else if (page <= 3) {
            pageNum = i + 1
          } else if (page >= totalPages - 2) {
            pageNum = totalPages - 4 + i
          } else {
            pageNum = page - 2 + i
          }
          return (
            <button
              key={pageNum}
              onClick={() => onPageChange(pageNum)}
              className={`w-10 h-10 rounded-full text-sm transition-all ${
                page === pageNum
                  ? 'bg-[#1A1A1A] text-white'
                  : 'text-[#666] hover:bg-[#F5F5F5]'
              }`}
            >
              {pageNum}
            </button>
          )
        })}
      </div>

      <button
        onClick={() => onPageChange(page + 1)}
        disabled={page >= totalPages}
        className="w-10 h-10 rounded-full border border-[#E5E5E5] text-sm text-[#666] hover:bg-[#F5F5F5] disabled:opacity-30 disabled:cursor-not-allowed transition-all"
      >
        →
      </button>
    </div>
  )
}

export default function Books() {
  const [searchParams, setSearchParams] = useSearchParams()
  const [booksList, setBooksList] = useState([])
  const [total, setTotal] = useState(0)
  const [loading, setLoading] = useState(true)
  const [page, setPage] = useState(1)
  const pageSize = 12

  const [filters, setFilters] = useState({
    mode: searchParams.get('mode') ? Number(searchParams.get('mode')) : null,
    category: searchParams.get('category') || null,
    keyword: searchParams.get('keyword') || null,
    sort: searchParams.get('sort') || 'created_at',
    order: searchParams.get('order') || 'desc',
  })

  useEffect(() => {
    async function fetchBooks() {
      setLoading(true)
      try {
        const params = {
          page,
          page_size: pageSize,
          ...(filters.mode && { mode: filters.mode }),
          ...(filters.category && { category: filters.category }),
          ...(filters.keyword && { keyword: filters.keyword }),
          sort: filters.sort,
          order: filters.order,
        }
        const res = await books.list(params)
        setBooksList(res.data?.list || [])
        setTotal(res.data?.total || 0)
      } catch (err) {
        console.error('Failed to fetch books:', err)
      } finally {
        setLoading(false)
      }
    }
    fetchBooks()
  }, [page, filters])

  const totalPages = Math.ceil(total / pageSize)

  return (
    <div className="min-h-screen bg-white">
      <Navbar />
      <FilterBar filters={filters} setFilters={setFilters} total={total} />

      <main className="max-w-6xl mx-auto px-8 py-8">
        {loading ? (
          <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
            {[...Array(8)].map((_, i) => (
              <div key={i}>
                <div className="aspect-[3/4] bg-[#F5F5F5] rounded-xl animate-pulse mb-4" />
                <div className="h-4 bg-[#F5F5F5] rounded animate-pulse w-3/4 mb-2" />
                <div className="h-3 bg-[#F5F5F5] rounded animate-pulse w-1/2" />
              </div>
            ))}
          </div>
        ) : booksList.length > 0 ? (
          <>
            <div className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-8">
              {booksList.map(book => (
                <BookCard key={book.id} book={book} />
              ))}
            </div>
            {totalPages > 1 && (
              <Pagination page={page} totalPages={totalPages} onPageChange={setPage} />
            )}
          </>
        ) : (
          <div className="text-center py-24">
            <svg className="w-16 h-16 mx-auto text-[#DDD] mb-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
            </svg>
            <p className="text-[#999] text-lg mb-4">暂无图书</p>
            <Link to="/publish" className="text-[#E85D2B] text-sm hover:underline">
              成为第一个发布者
            </Link>
          </div>
        )}
      </main>

      <footer className="py-12 px-8 border-t border-[#F0F0F0]">
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
