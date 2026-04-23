import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { auth } from '../../api'
import { useUser } from '../../context/UserContext'
import BackButton from '../../components/BackButton'

export default function Register() {
  const [account, setAccount] = useState('')
  const [password, setPassword] = useState('')
  const [confirmPassword, setConfirmPassword] = useState('')
  const [birthDate, setBirthDate] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login } = useUser()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')

    if (password !== confirmPassword) {
      setError('两次密码输入不一致')
      return
    }

    const age = new Date().getFullYear() - new Date(birthDate).getFullYear()
    if (age >= 30) {
      setError('年龄需小于30岁')
      return
    }

    setLoading(true)

    try {
      await auth.register(account, password, birthDate)
      const res = await auth.login(account, password)
      if (res.data?.token) {
        login(res.data.token, res.data.user)
        navigate('/')
      }
    } catch (err) {
      setError(err.message || '注册失败')
    } finally {
      setLoading(false)
    }
  }

  return (
    <div className="min-h-screen bg-white flex items-center justify-center px-6">
      <div className="absolute top-8 left-8">
        <BackButton />
      </div>
      <div className="w-full max-w-md">
        <div className="text-center mb-12">
          <h1 className="font-serif text-4xl text-primary mb-3">注册</h1>
          <p className="text-secondary">加入高校图书流转大家庭</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="bg-red-50 text-red-600 px-4 py-3 rounded-lg text-sm">
              {error}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              账号（手机号或昵称）
            </label>
            <input
              type="text"
              value={account}
              onChange={(e) => setAccount(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:border-primary transition-colors"
              placeholder="请输入账号"
              required
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              出生日期
            </label>
            <input
              type="date"
              value={birthDate}
              onChange={(e) => setBirthDate(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:border-primary transition-colors"
              required
            />
            <p className="text-xs text-secondary mt-1">需小于30岁</p>
          </div>

          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              密码
            </label>
            <input
              type="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:border-primary transition-colors"
              placeholder="6位以上密码"
              required
              minLength={6}
            />
          </div>

          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              确认密码
            </label>
            <input
              type="password"
              value={confirmPassword}
              onChange={(e) => setConfirmPassword(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:border-primary transition-colors"
              placeholder="再次输入密码"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary text-white rounded-lg font-medium hover:bg-gray-800 transition-colors disabled:opacity-50"
          >
            {loading ? '注册中...' : '注册'}
          </button>
        </form>

        <p className="text-center mt-8 text-secondary text-sm">
          已有账号？{' '}
          <Link to="/login" className="text-primary underline">
            立即登录
          </Link>
        </p>
      </div>
    </div>
  )
}
