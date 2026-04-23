import { useState } from 'react'
import { Link, useNavigate } from 'react-router-dom'
import { auth } from '../../api'
import { useUser } from '../../context/UserContext'
import BackButton from '../../components/BackButton'

export default function Login() {
  const [account, setAccount] = useState('')
  const [password, setPassword] = useState('')
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState('')
  const navigate = useNavigate()
  const { login } = useUser()

  const handleSubmit = async (e) => {
    e.preventDefault()
    setError('')
    setLoading(true)

    try {
      const res = await auth.login(account, password)
      if (res.data?.token) {
        login(res.data.token, res.data.user)
        navigate('/')
      }
    } catch (err) {
      setError(err.message || '登录失败')
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
          <h1 className="font-serif text-4xl text-primary mb-3">登录</h1>
          <p className="text-secondary">欢迎回到高校图书流转</p>
        </div>

        <form onSubmit={handleSubmit} className="space-y-6">
          {error && (
            <div className="bg-red-50 text-red-600 px-4 py-3 rounded-lg text-sm">
              {error}
            </div>
          )}

          <div>
            <label className="block text-sm font-medium text-primary mb-2">
              账号
            </label>
            <input
              type="text"
              value={account}
              onChange={(e) => setAccount(e.target.value)}
              className="w-full px-4 py-3 border border-gray-200 rounded-lg focus:outline-none focus:border-primary transition-colors"
              placeholder="手机号或昵称"
              required
            />
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
              placeholder="请输入密码"
              required
            />
          </div>

          <button
            type="submit"
            disabled={loading}
            className="w-full py-3 bg-primary text-white rounded-lg font-medium hover:bg-gray-800 transition-colors disabled:opacity-50"
          >
            {loading ? '登录中...' : '登录'}
          </button>
        </form>

        <p className="text-center mt-8 text-secondary text-sm">
          还没有账号？{' '}
          <Link to="/register" className="text-primary underline">
            立即注册
          </Link>
        </p>
      </div>
    </div>
  )
}
