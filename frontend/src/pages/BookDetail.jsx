import { useState, useEffect } from 'react'
import { useParams, useNavigate, Link } from 'react-router-dom'
import { books, rentals, sells, gifts } from '../api'
import { useUser } from '../context/UserContext'
import BackButton from '../components/BackButton'

const modeLabels = {
  1: { name: '租赁', color: '#E85D2B', deposit: true, rent: true },
  2: { name: '出售', color: '#2B9CD8', price: true },
  3: { name: '赠送', color: '#4ADE80', gift: true },
}

export default function BookDetail() {
  const { id } = useParams()
  const navigate = useNavigate()
  const { user } = useUser()
  const [book, setBook] = useState(null)
  const [loading, setLoading] = useState(true)
  const [actionLoading, setActionLoading] = useState(false)
  const [rentDays, setRentDays] = useState(7)
  const [showRentModal, setShowRentModal] = useState(false)

  useEffect(() => {
    async function fetchBook() {
      try {
        const res = await books.get(id)
        setBook(res.data)
      } catch (err) {
        console.error('Failed to fetch book:', err)
        navigate('/')
      } finally {
        setLoading(false)
      }
    }
    fetchBook()
  }, [id, navigate])

  const handleRent = async () => {
    if (!user) {
      navigate('/login')
      return
    }
    setActionLoading(true)
    try {
      const res = await rentals.create({ book_id: book.id, rent_days: rentDays })
      if (res.data?.id) {
        await rentals.pay(res.data.id)
        alert('租赁申请已提交！')
        navigate('/orders')
      }
    } catch (err) {
      alert(err.message || '租赁失败')
    } finally {
      setActionLoading(false)
    }
  }

  const handleBuy = async () => {
    if (!user) {
      navigate('/login')
      return
    }
    setActionLoading(true)
    try {
      const res = await sells.create({ book_id: book.id })
      if (res.data?.id) {
        await sells.pay(res.data.id)
        alert('购买订单已创建！')
        navigate('/orders')
      }
    } catch (err) {
      alert(err.message || '购买失败')
    } finally {
      setActionLoading(false)
    }
  }

  const handleApply = async () => {
    if (!user) {
      navigate('/login')
      return
    }
    setActionLoading(true)
    try {
      await gifts.create({ book_id: book.id })
      alert('申请已提交，等待赠送人确认')
      navigate('/orders')
    } catch (err) {
      alert(err.message || '申请失败')
    } finally {
      setActionLoading(false)
    }
  }

  if (loading) {
    return (
      <div className="min-h-screen bg-white flex items-center justify-center">
        <div className="animate-pulse text-secondary">加载中...</div>
      </div>
    )
  }

  if (!book) return null

  const mode = modeLabels[book.mode] || modeLabels[1]
  const isOwner = user?.id === book.user_id

  return (
    <div className="min-h-screen bg-white">
      <header className="border-b border-gray-100 py-4 px-6">
        <div className="max-w-4xl mx-auto flex items-center justify-between">
          <BackButton />
          {user && !isOwner && (
            <button className="text-sm text-secondary hover:text-primary transition-colors">
              举报
            </button>
          )}
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-4 md:px-6 py-8 md:py-12">
        <div className="grid md:grid-cols-2 gap-8 md:gap-12">
          <div className="aspect-[3/4] bg-gray-100 rounded-2xl overflow-hidden shadow-card">
            {book.cover_image ? (
              <img src={book.cover_image} alt={book.title} className="w-full h-full object-cover" />
            ) : (
              <div className="w-full h-full flex items-center justify-center text-gray-300">
                <svg className="w-20 md:w-24 h-20 md:h-24" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
              </div>
            )}
          </div>

          <div>
            <span
              className="inline-block px-3 py-1 rounded-full text-sm text-white mb-4"
              style={{ backgroundColor: mode.color }}
            >
              {mode.name}
            </span>

            <h1 className="font-serif text-2xl md:text-3xl text-primary mb-2">{book.title}</h1>
            <p className="text-secondary mb-6 text-sm md:text-base">{book.author || '未知作者'}</p>

            {book.daily_rent && (
              <div className="mb-4">
                <span className="text-3xl font-serif text-accent-orange">¥{book.daily_rent}</span>
                <span className="text-secondary">/天</span>
              </div>
            )}

            {book.sell_price && book.mode === 2 && (
              <div className="mb-4">
                <span className="text-3xl font-serif text-accent-blue">¥{book.sell_price}</span>
              </div>
            )}

            {book.mode === 3 && (
              <div className="mb-4">
                <span className="text-3xl font-serif text-accent-green">免费</span>
              </div>
            )}

            <div className="space-y-3 mb-8 text-sm">
              {book.deposit && (
                <div className="flex justify-between">
                  <span className="text-secondary">押金</span>
                  <span className="text-primary">¥{book.deposit}</span>
                </div>
              )}
              {book.min_rent_days && (
                <div className="flex justify-between">
                  <span className="text-secondary">最短租期</span>
                  <span className="text-primary">{book.min_rent_days}天</span>
                </div>
              )}
              <div className="flex justify-between">
                <span className="text-secondary">所在学校</span>
                <span className="text-primary">{book.User?.school || '未填写'}</span>
              </div>
              <div className="flex justify-between">
                <span className="text-secondary">浏览次数</span>
                <span className="text-primary">{book.view_count || 0}</span>
              </div>
            </div>

            {book.description && (
              <div className="mb-8">
                <h3 className="text-sm font-medium text-primary mb-2">图书描述</h3>
                <p className="text-secondary text-sm leading-relaxed">{book.description}</p>
              </div>
            )}

            {book.pickup_location && (
              <div className="mb-8">
                <h3 className="text-sm font-medium text-primary mb-2">取书地点</h3>
                <p className="text-secondary text-sm">{book.pickup_location}</p>
              </div>
            )}

            {book.images?.length > 0 && (
              <div className="mb-8">
                <h3 className="text-sm font-medium text-primary mb-3">实物图片</h3>
                <div className="grid grid-cols-2 md:grid-cols-3 gap-2 md:gap-3">
                  {book.images.map((img, i) => (
                    <div key={i} className="aspect-square bg-gray-100 rounded-lg overflow-hidden">
                      <img src={img} alt="" className="w-full h-full object-cover" />
                    </div>
                  ))}
                </div>
              </div>
            )}

            {!isOwner && (
              <div className="space-y-3 md:space-y-4">
                {book.mode === 1 && (
                  <>
                    <div className="flex items-center gap-3 flex-wrap">
                      <label className="text-sm text-secondary">租赁天数:</label>
                      <input
                        type="number"
                        min={book.min_rent_days || 1}
                        value={rentDays}
                        onChange={(e) => setRentDays(Number(e.target.value))}
                        className="w-20 px-3 py-2 border border-gray-200 rounded-lg text-center text-sm"
                      />
                      <span className="text-sm text-secondary">
                        预计: ¥{(book.daily_rent * rentDays + book.deposit).toFixed(2)}
                      </span>
                    </div>
                    <button
                      onClick={handleRent}
                      disabled={actionLoading}
                      className="w-full py-3 rounded-lg text-white font-medium transition-colors disabled:opacity-50"
                      style={{ backgroundColor: mode.color }}
                    >
                      {actionLoading ? '处理中...' : '立即租赁'}
                    </button>
                  </>
                )}

                {book.mode === 2 && (
                  <button
                    onClick={handleBuy}
                    disabled={actionLoading}
                    className="w-full py-3 rounded-lg text-white font-medium transition-colors disabled:opacity-50"
                    style={{ backgroundColor: mode.color }}
                  >
                    {actionLoading ? '处理中...' : '立即购买'}
                  </button>
                )}

                {book.mode === 3 && (
                  <button
                    onClick={handleApply}
                    disabled={actionLoading}
                    className="w-full py-3 rounded-lg text-white font-medium transition-colors disabled:opacity-50"
                    style={{ backgroundColor: mode.color }}
                  >
                    {actionLoading ? '处理中...' : '申请领取'}
                  </button>
                )}
              </div>
            )}

            {isOwner && (
              <div className="space-y-4">
                <p className="text-center text-secondary">这是您发布的图书</p>
                <Link
                  to="/my-books"
                  className="block w-full py-3 bg-gray-100 rounded-lg text-center text-primary font-medium"
                >
                  管理我的图书
                </Link>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
