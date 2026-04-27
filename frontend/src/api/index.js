import './mock.js'

const USE_MOCK = true

async function request(endpoint, options = {}) {
  if (USE_MOCK) {
    return window.mockRequest(endpoint, options)
  }

  const token = window.getToken()
  const headers = {
    'Content-Type': 'application/json',
    ...(token && { 'Authorization': `Bearer ${token}` }),
    ...options.headers,
  }

  const response = await fetch(`/api/v1${endpoint}`, {
    ...options,
    headers,
  })

  const data = await response.json()

  if (!response.ok) {
    throw new Error(data.message || '请求失败')
  }

  return data
}

export const auth = {
  async register(account, password, birth_date) {
    const data = await request('/users/register', {
      method: 'POST',
      body: JSON.stringify({ account, password, birth_date, code: '123456' }),
    })
    return data
  },

  async login(account, password) {
    const data = await request('/users/login', {
      method: 'POST',
      body: JSON.stringify({ account, password }),
    })
    if (data.data?.token) {
      window.setToken(data.data.token)
    }
    return data
  },

  logout() {
    window.setToken(null)
  },

  getUser() {
    return request('/users/me')
  },

  updateProfile(data) {
    return request('/users/me', {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },
}

export const books = {
  list(params = {}) {
    const query = new URLSearchParams(params).toString()
    return request(`/books${query ? '?' + query : ''}`)
  },

  get(id) {
    return request(`/books/${id}`)
  },

  create(data) {
    return request('/books', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  update(id, data) {
    return request(`/books/${id}`, {
      method: 'PUT',
      body: JSON.stringify(data),
    })
  },

  delete(id) {
    return request(`/books/${id}`, {
      method: 'DELETE',
    })
  },

  getHot() {
    return request('/books/hot')
  },

  getLatestGift() {
    return request('/books/gift/latest')
  },

  getMyBooks(params = {}) {
    const query = new URLSearchParams(params).toString()
    return request(`/books/my${query ? '?' + query : ''}`)
  },

  getFavorites(params = {}) {
    const query = new URLSearchParams(params).toString()
    return request(`/books/favorites${query ? '?' + query : ''}`)
  },

  addFavorite(id) {
    return request(`/books/${id}/favorite`, {
      method: 'POST',
    })
  },

  removeFavorite(id) {
    return request(`/books/${id}/favorite`, {
      method: 'DELETE',
    })
  },

  queryISBN(isbn) {
    return request(`/books/isbn/${isbn}`)
  },
}

export const rentals = {
  create(data) {
    return request('/rentals', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  pay(id) {
    return request(`/rentals/${id}/pay`, {
      method: 'POST',
      body: JSON.stringify({ payment_method: 'wallet' }),
    })
  },

  cancel(id) {
    return request(`/rentals/${id}/cancel`, {
      method: 'POST',
    })
  },

  confirmPickup(id) {
    return request(`/rentals/${id}/confirm-pickup`, {
      method: 'POST',
    })
  },

  returnBook(id) {
    return request(`/rentals/${id}/return`, {
      method: 'POST',
    })
  },

  inspect(id, passed, remark = '') {
    return request(`/rentals/${id}/inspect`, {
      method: 'POST',
      body: JSON.stringify({ passed, remark }),
    })
  },

  get(id) {
    return request(`/rentals/${id}`)
  },

  myList(params = {}) {
    const query = new URLSearchParams(params).toString()
    return request(`/rentals/my${query ? '?' + query : ''}`)
  },

  rate(id, rating, comment = '') {
    return request(`/rentals/${id}/rate`, {
      method: 'POST',
      body: JSON.stringify({ rating, comment }),
    })
  },
}

export const sells = {
  create(data) {
    return request('/sells', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  pay(id) {
    return request(`/sells/${id}/pay`, {
      method: 'POST',
      body: JSON.stringify({ payment_method: 'wallet' }),
    })
  },

  cancel(id) {
    return request(`/sells/${id}/cancel`, {
      method: 'POST',
    })
  },

  ship(id, delivery_type, delivery_company = '', delivery_no = '') {
    return request(`/sells/${id}/ship`, {
      method: 'POST',
      body: JSON.stringify({ delivery_type, delivery_company, delivery_no }),
    })
  },

  confirmPickup(id) {
    return request(`/sells/${id}/pickup`, {
      method: 'POST',
    })
  },

  confirmReceive(id) {
    return request(`/sells/${id}/confirm`, {
      method: 'POST',
    })
  },

  get(id) {
    return request(`/sells/${id}`)
  },

  myList(params = {}) {
    const query = new URLSearchParams(params).toString()
    return request(`/sells/my${query ? '?' + query : ''}`)
  },

  rate(id, rating, comment = '') {
    return request(`/sells/${id}/rate`, {
      method: 'POST',
      body: JSON.stringify({ rating, comment }),
    })
  },
}

export const gifts = {
  create(data) {
    return request('/gifts', {
      method: 'POST',
      body: JSON.stringify(data),
    })
  },

  confirm(id) {
    return request(`/gifts/${id}/confirm`, {
      method: 'POST',
    })
  },

  reject(id) {
    return request(`/gifts/${id}/reject`, {
      method: 'POST',
    })
  },

  cancel(id) {
    return request(`/gifts/${id}/cancel`, {
      method: 'POST',
    })
  },

  deliver(id) {
    return request(`/gifts/${id}/deliver`, {
      method: 'POST',
    })
  },

  get(id) {
    return request(`/gifts/${id}`)
  },

  myList(params = {}) {
    const query = new URLSearchParams(params).toString()
    return request(`/gifts/my${query ? '?' + query : ''}`)
  },

  getApplications(bookId) {
    return request(`/gifts/book/${bookId}/applications`)
  },

  rate(id, rating, comment = '') {
    return request(`/gifts/${id}/rate`, {
      method: 'POST',
      body: JSON.stringify({ rating, comment }),
    })
  },
}

export const common = {
  upload(file) {
    return Promise.resolve({ data: { url: URL.createObjectURL(file) } })
  },
}

export const guide = {
  get() {
    return request('/guide')
  },
}
