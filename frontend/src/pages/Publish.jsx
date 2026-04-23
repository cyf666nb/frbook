import { useState, useRef } from 'react'
import { useNavigate } from 'react-router-dom'
import { books, common } from '../api'
import { useUser } from '../context/UserContext'
import BackButton from '../components/BackButton'
import { useToast } from '../components/Toast'

const modeOptions = [
  { value: 1, label: '租赁', description: '收取日租金和押金', color: '#E85D2B' },
  { value: 2, label: '出售', description: '一次性出售图书', color: '#2B9CD8' },
  { value: 3, label: '赠送', description: '免费赠送给需要的人', color: '#4ADE80' },
]

const categoryOptions = ['教材', '文学', '科技', '艺术', '历史', '其他']

export default function Publish() {
  const navigate = useNavigate()
  const { user } = useUser()
  const { toast } = useToast()
  const fileInputRef = useRef(null)
  const [loading, setLoading] = useState(false)
  const [uploading, setUploading] = useState(false)
  const [showCategoryDropdown, setShowCategoryDropdown] = useState(false)
  const [form, setForm] = useState({
    mode: 1,
    title: '',
    author: '',
    isbn: '',
    description: '',
    category: '',
    pickup_location: '',
    daily_rent: '',
    weekly_rent: '',
    deposit: '',
    min_rent_days: 1,
    sell_price: '',
    images: [],
  })

  if (!user) {
    navigate('/login')
    return null
  }

  const handleSubmit = async (e) => {
    e.preventDefault()
    setLoading(true)

    try {
      const data = {
        ...form,
        mode: Number(form.mode),
        daily_rent: form.daily_rent ? Number(form.daily_rent) : null,
        weekly_rent: form.weekly_rent ? Number(form.weekly_rent) : null,
        deposit: form.deposit ? Number(form.deposit) : null,
        min_rent_days: Number(form.min_rent_days),
        sell_price: form.sell_price ? Number(form.sell_price) : null,
      }
      await books.create(data)
      toast('发布成功！')
      navigate('/books')
    } catch (err) {
      toast(err.message || '发布失败', 'error')
    } finally {
      setLoading(false)
    }
  }

  const handleISBNLookup = async () => {
    if (!form.isbn) return
    try {
      const res = await books.queryISBN(form.isbn)
      if (res.data) {
        setForm(prev => ({
          ...prev,
          title: res.data.title || prev.title,
          author: res.data.author || prev.author,
          cover_image: res.data.cover,
        }))
      }
    } catch (err) {
      console.error('ISBN lookup failed:', err)
    }
  }

  const handleFileSelect = async (e) => {
    const files = Array.from(e.target.files)
    if (files.length === 0) return

    setUploading(true)
    try {
      const uploadedUrls = []
      for (const file of files) {
        const res = await common.upload(file)
        if (res.data?.url) {
          uploadedUrls.push(res.data.url)
        }
      }
      setForm(prev => ({
        ...prev,
        images: [...prev.images, ...uploadedUrls].slice(0, 5)
      }))
    } catch (err) {
      toast('图片上传失败', 'error')
    } finally {
      setUploading(false)
    }
  }

  const handleCapture = () => {
    if (fileInputRef.current) {
      fileInputRef.current.click()
    }
  }

  const removeImage = (index) => {
    setForm(prev => ({
      ...prev,
      images: prev.images.filter((_, i) => i !== index)
    }))
  }

  return (
    <div className="min-h-screen bg-[#FAFAFA]">
      <header className="bg-white border-b border-gray-100 py-4 px-6 sticky top-0 z-40">
        <div className="max-w-2xl mx-auto">
          <BackButton />
        </div>
      </header>

      <main className="max-w-2xl mx-auto px-6 py-12">
        <h1 className="font-serif text-3xl text-[#1A1A1A] mb-8">发布图书</h1>

        <form onSubmit={handleSubmit} className="space-y-8">
          <div className="bg-white rounded-2xl p-6 shadow-sm">
            <label className="block text-sm font-medium text-[#1A1A1A] mb-3">流转方式</label>
            <div className="grid grid-cols-3 gap-3">
              {modeOptions.map(option => (
                <button
                  key={option.value}
                  type="button"
                  onClick={() => setForm(prev => ({ ...prev, mode: option.value }))}
                  className={`p-4 rounded-xl border-2 transition-all text-left ${
                    form.mode === option.value
                      ? 'border-[#1A1A1A] bg-[#FAFAFA]'
                      : 'border-[#E5E5E5] hover:border-[#CCC]'
                  }`}
                >
                  <div className="font-medium text-[#1A1A1A] flex items-center gap-2">
                    <span 
                      className="w-2 h-2 rounded-full" 
                      style={{ backgroundColor: option.color }}
                    />
                    {option.label}
                  </div>
                  <div className="text-xs text-[#666] mt-1">{option.description}</div>
                </button>
              ))}
            </div>
          </div>

          <div className="bg-white rounded-2xl p-6 shadow-sm space-y-6">
            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-[#1A1A1A] mb-2">ISBN（可扫码）</label>
                <div className="flex gap-2">
                  <input
                    type="text"
                    value={form.isbn}
                    onChange={(e) => setForm(prev => ({ ...prev, isbn: e.target.value }))}
                    className="flex-1 px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                    placeholder="扫描或输入ISBN"
                  />
                  <button
                    type="button"
                    onClick={handleISBNLookup}
                    className="px-4 py-3 bg-[#1A1A1A] text-white rounded-xl text-sm hover:bg-[#333] transition-colors"
                  >
                    查询
                  </button>
                </div>
              </div>
              <div>
                <label className="block text-sm font-medium text-[#1A1A1A] mb-2">书名 *</label>
                <input
                  type="text"
                  value={form.title}
                  onChange={(e) => setForm(prev => ({ ...prev, title: e.target.value }))}
                  className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                  placeholder="书名"
                  required
                />
              </div>
            </div>

            <div className="grid md:grid-cols-2 gap-6">
              <div>
                <label className="block text-sm font-medium text-[#1A1A1A] mb-2">作者</label>
                <input
                  type="text"
                  value={form.author}
                  onChange={(e) => setForm(prev => ({ ...prev, author: e.target.value }))}
                  className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                  placeholder="作者"
                />
              </div>
              <div className="relative">
                <label className="block text-sm font-medium text-[#1A1A1A] mb-2">分类</label>
                <button
                  type="button"
                  onClick={() => setShowCategoryDropdown(!showCategoryDropdown)}
                  className={`w-full px-4 py-3 rounded-xl text-sm text-left flex items-center justify-between transition-all ${
                    form.category
                      ? 'bg-[#1A1A1A] text-white'
                      : 'bg-[#F5F5F5] text-[#666]'
                  }`}
                >
                  <span>{form.category || '选择分类'}</span>
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                  </svg>
                </button>
                {showCategoryDropdown && (
                  <div className="absolute top-full left-0 right-0 mt-2 bg-white rounded-xl shadow-lg border border-gray-100 py-2 z-50">
                    {categoryOptions.map(cat => (
                      <button
                        key={cat}
                        type="button"
                        onClick={() => {
                          setForm(prev => ({ ...prev, category: cat }))
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
            </div>
          </div>

          {form.mode === 1 && (
            <div className="bg-white rounded-2xl p-6 shadow-sm">
              <label className="block text-sm font-medium text-[#1A1A1A] mb-3">租赁设置</label>
              <div className="grid grid-cols-3 gap-4">
                <div>
                  <label className="block text-xs text-[#666] mb-2">日租金(元) *</label>
                  <input
                    type="number"
                    value={form.daily_rent}
                    onChange={(e) => setForm(prev => ({ ...prev, daily_rent: e.target.value }))}
                    className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                    placeholder="0.00"
                    required={form.mode === 1}
                  />
                </div>
                <div>
                  <label className="block text-xs text-[#666] mb-2">押金(元) *</label>
                  <input
                    type="number"
                    value={form.deposit}
                    onChange={(e) => setForm(prev => ({ ...prev, deposit: e.target.value }))}
                    className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                    placeholder="0.00"
                    required={form.mode === 1}
                  />
                </div>
                <div>
                  <label className="block text-xs text-[#666] mb-2">最短租期(天)</label>
                  <input
                    type="number"
                    value={form.min_rent_days}
                    onChange={(e) => setForm(prev => ({ ...prev, min_rent_days: e.target.value }))}
                    className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                    min={1}
                  />
                </div>
              </div>
            </div>
          )}

          {form.mode === 2 && (
            <div className="bg-white rounded-2xl p-6 shadow-sm">
              <label className="block text-sm font-medium text-[#1A1A1A] mb-2">出售价格(元) *</label>
              <input
                type="number"
                value={form.sell_price}
                onChange={(e) => setForm(prev => ({ ...prev, sell_price: e.target.value }))}
                className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                placeholder="0.00"
                required={form.mode === 2}
              />
            </div>
          )}

          {form.mode === 3 && (
            <div className="bg-[#4ADE80]/10 rounded-2xl p-4">
              <p className="text-[#4ADE80] text-sm text-center">
                免费赠送模式，无需填写价格
              </p>
            </div>
          )}

          <div className="bg-white rounded-2xl p-6 shadow-sm space-y-6">
            <div>
              <label className="block text-sm font-medium text-[#1A1A1A] mb-2">图书描述</label>
              <textarea
                value={form.description}
                onChange={(e) => setForm(prev => ({ ...prev, description: e.target.value }))}
                className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10 resize-none"
                rows={4}
                placeholder="描述图书的新旧程度、亮点等"
              />
            </div>

            <div>
              <label className="block text-sm font-medium text-[#1A1A1A] mb-2">取书地点</label>
              <input
                type="text"
                value={form.pickup_location}
                onChange={(e) => setForm(prev => ({ ...prev, pickup_location: e.target.value }))}
                className="w-full px-4 py-3 bg-[#F5F5F5] rounded-xl text-sm focus:outline-none focus:ring-2 focus:ring-[#1A1A1A]/10"
                placeholder="如：北京大学东门图书馆门口"
              />
            </div>
          </div>

          <div className="bg-white rounded-2xl p-6 shadow-sm">
            <label className="block text-sm font-medium text-[#1A1A1A] mb-3">实物图片</label>
            
            <input
              ref={fileInputRef}
              type="file"
              accept="image/*"
              multiple
              capture="environment"
              onChange={handleFileSelect}
              className="hidden"
            />

            <div className="grid grid-cols-4 gap-3 mb-4">
              {form.images.map((img, index) => (
                <div key={index} className="relative aspect-square rounded-xl overflow-hidden bg-gray-100">
                  <img src={img} alt="" className="w-full h-full object-cover" />
                  <button
                    type="button"
                    onClick={() => removeImage(index)}
                    className="absolute top-1 right-1 w-6 h-6 bg-black/50 rounded-full flex items-center justify-center text-white hover:bg-black/70 transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </button>
                </div>
              ))}
              
              {form.images.length < 5 && (
                <button
                  type="button"
                  onClick={handleCapture}
                  disabled={uploading}
                  className="aspect-square rounded-xl border-2 border-dashed border-[#E5E5E5] flex flex-col items-center justify-center gap-1 hover:border-[#1A1A1A] transition-colors disabled:opacity-50"
                >
                  {uploading ? (
                    <div className="w-5 h-5 border-2 border-[#1A1A1A] border-t-transparent rounded-full animate-spin" />
                  ) : (
                    <>
                      <svg className="w-6 h-6 text-[#999]" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M12 4v16m8-8H4" />
                      </svg>
                      <span className="text-xs text-[#999]">添加</span>
                    </>
                  )}
                </button>
              )}
            </div>

            <div className="flex gap-3">
              <button
                type="button"
                onClick={handleCapture}
                disabled={uploading}
                className="flex-1 py-3 bg-[#F5F5F5] rounded-xl text-sm text-[#666] hover:bg-[#EFEFEF] transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z" />
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M15 13a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                拍照
              </button>
              <button
                type="button"
                onClick={() => {
                  if (fileInputRef.current) {
                    fileInputRef.current.removeAttribute('capture')
                    fileInputRef.current.click()
                  }
                }}
                disabled={uploading}
                className="flex-1 py-3 bg-[#F5F5F5] rounded-xl text-sm text-[#666] hover:bg-[#EFEFEF] transition-colors flex items-center justify-center gap-2 disabled:opacity-50"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1.5} d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                相册
              </button>
            </div>
            
            <p className="text-xs text-[#999] mt-3 text-center">最多上传5张图片</p>
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-4 bg-[#1A1A1A] text-white rounded-xl font-medium hover:bg-[#333] transition-colors disabled:opacity-50"
          >
            {loading ? '发布中...' : '发布图书'}
          </button>
        </form>
      </main>
    </div>
  )
}