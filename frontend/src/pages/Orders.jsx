import { useState, useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { rentals, sells, gifts } from '../api'
import { useUser } from '../context/UserContext'
import BackButton from '../components/BackButton'

const statusLabels = {
  0: { label: '待处理', color: '#666' },
  1: { label: '进行中', color: '#2B9CD8' },
  2: { label: '已完成', color: '#4ADE80' },
  3: { label: '已取消', color: '#999' },
  4: { label: '已取消', color: '#999' },
  5: { label: '已取消', color: '#999' },
  6: { label: '逾期', color: '#E85D2B' },
}

const rentalStatusLabels = {
  0: { label: '待支付', color: '#E85D2B' },
  1: { label: '待确认取书', color: '#2B9CD8' },
  2: { label: '租赁中', color: '#4ADE80' },
  3: { label: '待验收', color: '#9B59B6' },
  4: { label: '已完成', color: '#27AE60' },
  5: { label: '已取消', color: '#999' },
  6: { label: '逾期中', color: '#E74C3C' },
}

const sellStatusLabels = {
  0: { label: '待支付', color: '#E85D2B' },
  1: { label: '待发货', color: '#2B9CD8' },
  2: { label: '已发货', color: '#9B59B6' },
  3: { label: '已完成', color: '#27AE60' },
  4: { label: '已取消', color: '#999' },
}

const giftStatusLabels = {
  0: { label: '待确认', color: '#E85D2B' },
  1: { label: '待交付', color: '#2B9CD8' },
  2: { label: '已完成', color: '#27AE60' },
  3: { label: '已取消', color: '#999' },
}

export default function Orders() {
  const navigate = useNavigate()
  const { user } = useUser()
  const [activeTab, setActiveTab] = useState('rental')
  const [rentalsList, setRentalsList] = useState([])
  const [sellsList, setSellsList] = useState([])
  const [giftsList, setGiftsList] = useState([])
  const [loading, setLoading] = useState(true)

  useEffect(() => {
    if (!user) {
      navigate('/login')
      return
    }
    fetchOrders()
  }, [user, navigate])

  async function fetchOrders() {
    try {
      const [rentalRes, sellRes, giftRes] = await Promise.all([
        rentals.myList(),
        sells.myList(),
        gifts.myList(),
      ])
      setRentalsList(rentalRes.data?.list || [])
      setSellsList(sellRes.data?.list || [])
      setGiftsList(giftRes.data?.list || [])
    } catch (err) {
      console.error('Failed to fetch orders:', err)
    } finally {
      setLoading(false)
    }
  }

  if (!user) return null

  const tabs = [
    { key: 'rental', label: '租赁订单' },
    { key: 'sell', label: '购买订单' },
    { key: 'gift', label: '赠送记录' },
  ]

  const statusLabelsMap = {
    rental: rentalStatusLabels,
    sell: sellStatusLabels,
    gift: giftStatusLabels,
  }

  const currentList = {
    rental: rentalsList,
    sell: sellsList,
    gift: giftsList,
  }[activeTab] || []

  const currentStatusLabels = statusLabelsMap[activeTab]

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-100 py-4 px-6">
        <div className="max-w-4xl mx-auto flex items-center justify-between">
          <BackButton />
          <h1 className="font-serif text-xl text-primary">我的订单</h1>
          <div className="w-16" />
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-8">
        <div className="bg-white rounded-2xl shadow-card overflow-hidden">
          <div className="flex border-b border-gray-100">
            {tabs.map(tab => (
              <button
                key={tab.key}
                onClick={() => setActiveTab(tab.key)}
                className={`flex-1 py-4 text-center font-medium transition-colors ${
                  activeTab === tab.key
                    ? 'text-primary border-b-2 border-primary'
                    : 'text-secondary hover:text-primary'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </div>

          <div className="p-6">
            {loading ? (
              <div className="space-y-4">
                {[...Array(3)].map((_, i) => (
                  <div key={i} className="bg-gray-50 rounded-xl p-4 animate-pulse">
                    <div className="h-4 bg-gray-200 rounded w-1/3 mb-2" />
                    <div className="h-3 bg-gray-200 rounded w-1/2" />
                  </div>
                ))}
              </div>
            ) : currentList.length > 0 ? (
              <div className="space-y-4">
                {currentList.map(order => {
                  const status = currentStatusLabels[order.status] || { label: '未知', color: '#666' }
                  const isOwner = user.id === order.owner_id || user.id === order.seller_id || user.id === order.giver_id
                  const otherName = isOwner 
                    ? (order.renter_id === user.id ? order.owner_id : order.renter_id)
                    : (order.owner_id === user.id ? order.renter_id : order.owner_id)

                  return (
                    <Link
                      key={order.id}
                      to={`/order/${activeTab}/${order.id}`}
                      className="block bg-gray-50 rounded-xl p-4 hover:bg-gray-100 transition-colors"
                    >
                      <div className="flex items-start justify-between mb-2">
                        <div>
                          <h3 className="font-medium text-primary line-clamp-1">{order.book_title}</h3>
                          <p className="text-xs text-secondary mt-1">
                            订单号: {order.order_no || order.record_no}
                          </p>
                        </div>
                        <span 
                          className="px-2 py-1 rounded text-xs text-white"
                          style={{ backgroundColor: status.color }}
                        >
                          {status.label}
                        </span>
                      </div>
                      <div className="flex items-center justify-between text-sm">
                        <span className="text-secondary">
                          {isOwner ? '作为出租方' : '作为承租方'}
                        </span>
                        <span className="text-primary">
                          {order.total_amount || order.price ? `¥${order.total_amount || order.price}` : '免费'}
                        </span>
                      </div>
                    </Link>
                  )
                })}
              </div>
            ) : (
              <div className="text-center py-12">
                <p className="text-secondary mb-4">暂无订单</p>
                <Link
                  to="/"
                  className="text-accent-blue hover:underline"
                >
                  去逛逛
                </Link>
              </div>
            )}
          </div>
        </div>
      </main>
    </div>
  )
}
