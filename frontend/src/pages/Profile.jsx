import { useState, useEffect } from 'react'
import { useNavigate, Link } from 'react-router-dom'
import { books, auth } from '../api'
import { useUser } from '../context/UserContext'
import BackButton from '../components/BackButton'
import { useToast } from '../components/Toast'

export default function Profile() {
  const navigate = useNavigate()
  const { user, updateUser } = useUser()
  const { toast } = useToast()
  const [myBooks, setMyBooks] = useState([])
  const [loading, setLoading] = useState(true)
  const [editing, setEditing] = useState(false)
  const [form, setForm] = useState({
    nickname: '',
    school: '',
    campus: '',
    grade: '',
    interest_tags: '',
  })

  useEffect(() => {
    if (!user) {
      navigate('/login')
      return
    }
    setForm({
      nickname: user.nickname || '',
      school: user.school || '',
      campus: user.campus || '',
      grade: user.grade || '',
      interest_tags: user.interest_tags?.join(', ') || '',
    })
    fetchMyBooks()
  }, [user, navigate])

  async function fetchMyBooks() {
    try {
      const res = await books.getMyBooks()
      setMyBooks(res.data?.list || [])
    } catch (err) {
      console.error('Failed to fetch books:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleSave = async () => {
    try {
      const tags = form.interest_tags.split(',').map(t => t.trim()).filter(Boolean)
      await auth.updateProfile({ ...form, interest_tags: tags })
      const res = await auth.getUser()
      updateUser(res.data)
      setEditing(false)
      toast('保存成功')
    } catch (err) {
      toast(err.message || '保存失败', 'error')
    }
  }

  if (!user) return null

  return (
    <div className="min-h-screen bg-gray-50">
      <header className="bg-white border-b border-gray-100 py-4 px-6">
        <div className="max-w-4xl mx-auto flex items-center justify-between">
          <BackButton />
          <h1 className="font-serif text-xl text-primary">个人中心</h1>
          <div className="w-16" />
        </div>
      </header>

      <main className="max-w-4xl mx-auto px-6 py-12">
        <div className="bg-white rounded-2xl shadow-card p-8 mb-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="font-serif text-2xl text-primary">基本信息</h2>
            <button
              onClick={() => editing ? handleSave() : setEditing(true)}
              className="px-4 py-2 bg-primary text-white rounded-lg text-sm font-medium hover:bg-gray-800 transition-colors"
            >
              {editing ? '保存' : '编辑'}
            </button>
          </div>

          <div className="grid md:grid-cols-2 gap-6">
            <div>
              <label className="block text-sm font-medium text-secondary mb-2">账号</label>
              <p className="text-primary">{user.account}</p>
            </div>
            <div>
              <label className="block text-sm font-medium text-secondary mb-2">昵称</label>
              {editing ? (
                <input
                  type="text"
                  value={form.nickname}
                  onChange={(e) => setForm(prev => ({ ...prev, nickname: e.target.value }))}
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg"
                />
              ) : (
                <p className="text-primary">{user.nickname || '-'}</p>
              )}
            </div>
            <div>
              <label className="block text-sm font-medium text-secondary mb-2">学校</label>
              {editing ? (
                <input
                  type="text"
                  value={form.school}
                  onChange={(e) => setForm(prev => ({ ...prev, school: e.target.value }))}
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg"
                  placeholder="所在学校"
                />
              ) : (
                <p className="text-primary">{user.school || '-'}</p>
              )}
            </div>
            <div>
              <label className="block text-sm font-medium text-secondary mb-2">校区</label>
              {editing ? (
                <input
                  type="text"
                  value={form.campus}
                  onChange={(e) => setForm(prev => ({ ...prev, campus: e.target.value }))}
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg"
                  placeholder="校区"
                />
              ) : (
                <p className="text-primary">{user.campus || '-'}</p>
              )}
            </div>
            <div>
              <label className="block text-sm font-medium text-secondary mb-2">年级</label>
              {editing ? (
                <input
                  type="text"
                  value={form.grade}
                  onChange={(e) => setForm(prev => ({ ...prev, grade: e.target.value }))}
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg"
                  placeholder="年级"
                />
              ) : (
                <p className="text-primary">{user.grade || '-'}</p>
              )}
            </div>
            <div>
              <label className="block text-sm font-medium text-secondary mb-2">兴趣标签</label>
              {editing ? (
                <input
                  type="text"
                  value={form.interest_tags}
                  onChange={(e) => setForm(prev => ({ ...prev, interest_tags: e.target.value }))}
                  className="w-full px-4 py-2 border border-gray-200 rounded-lg"
                  placeholder="用逗号分隔"
                />
              ) : (
                <p className="text-primary">{user.interest_tags?.join(', ') || '-'}</p>
              )}
            </div>
            <div>
              <label className="block text-sm font-medium text-secondary mb-2">信用评分</label>
              <p className="text-primary text-xl font-serif">{user.credit_score || 100}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-2xl shadow-card p-8">
          <div className="flex items-center justify-between mb-6">
            <h2 className="font-serif text-2xl text-primary">我发布的图书</h2>
            <Link
              to="/publish"
              className="px-4 py-2 bg-accent-blue text-white rounded-lg text-sm font-medium hover:bg-blue-600 transition-colors"
            >
              发布新书
            </Link>
          </div>

          {loading ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {[...Array(4)].map((_, i) => (
                <div key={i} className="bg-gray-100 rounded-xl animate-pulse h-48" />
              ))}
            </div>
          ) : myBooks.length > 0 ? (
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4">
              {myBooks.map(book => (
                <Link
                  key={book.id}
                  to={`/book/${book.id}`}
                  className="block bg-gray-50 rounded-xl overflow-hidden hover:shadow-md transition-shadow"
                >
                  <div className="aspect-[3/4] bg-gray-100">
                    {book.cover_image && (
                      <img src={book.cover_image} alt={book.title} className="w-full h-full object-cover" />
                    )}
                  </div>
                  <div className="p-3">
                    <h3 className="text-sm font-medium text-primary line-clamp-1">{book.title}</h3>
                    <span 
                      className="inline-block mt-2 px-2 py-0.5 rounded text-xs text-white"
                      style={{ 
                        backgroundColor: book.mode === 1 ? '#E85D2B' : book.mode === 2 ? '#2B9CD8' : '#4ADE80' 
                      }}
                    >
                      {book.mode === 1 ? '租赁' : book.mode === 2 ? '出售' : '赠送'}
                    </span>
                  </div>
                </Link>
              ))}
            </div>
          ) : (
            <p className="text-center text-secondary py-12">暂无发布的图书</p>
          )}
        </div>
      </main>
    </div>
  )
}
