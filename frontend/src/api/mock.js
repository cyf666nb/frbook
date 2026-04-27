const STORAGE_KEY = 'bookshare_data';

function getStoredData() {
  const data = localStorage.getItem(STORAGE_KEY);
  if (data) {
    const parsed = JSON.parse(data);
    if (!parsed.books || parsed.books.length === 0) {
      return getInitialData();
    }
    return parsed;
  }
  return getInitialData();
}

function saveData(data) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(data));
}

function getInitialData() {
  const users = [
    { id: 1, account: '13800138000', password: '123456', birth_date: '2000-01-01', nickname: '测试用户1', school: '北京大学', campus: '海淀校区', grade: '2020级', credit_score: 100, status: 1, created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
    { id: 2, account: 'testuser', password: '123456', birth_date: '1999-05-15', nickname: '测试用户2', school: '清华大学', campus: '五道口校区', grade: '2021级', credit_score: 100, status: 1, created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' },
    { id: 3, account: '13599887766', password: '123456', birth_date: '2001-03-20', nickname: '测试用户3', school: '复旦大学', campus: '邯郸校区', grade: '2022级', credit_score: 100, status: 1, created_at: '2024-01-01T00:00:00Z', updated_at: '2024-01-01T00:00:00Z' }
  ];

  const wallets = [
    { id: 1, user_id: 1, balance: 1000, frozen_balance: 0, total_income: 0, total_expense: 0 },
    { id: 2, user_id: 2, balance: 1000, frozen_balance: 0, total_income: 0, total_expense: 0 },
    { id: 3, user_id: 3, balance: 1000, frozen_balance: 0, total_income: 0, total_expense: 0 }
  ];

  const books = [
    { id: 1, user_id: 1, isbn: '978-7530211538', title: '活着', author: '余华', publisher: '北京十月文艺出版社', cover_image: 'https://picsum.photos/seed/huozhe123/400/600', description: '讲述农民福贵悲惨的一生，是中国当代文学的经典之作。', category: '文学', mode: 2, status: 1, daily_rent: 2, weekly_rent: 3, deposit: 15, sell_price: 35, min_rent_days: 1, images: ['https://picsum.photos/seed/huozhe123/400/600'], pickup_location: '北京大学东门', view_count: 45890, created_at: '2024-01-15T10:00:00Z', updated_at: '2024-01-15T10:00:00Z' },
    { id: 2, user_id: 2, isbn: '978-7229100640', title: '三体', author: '刘慈欣', publisher: '重庆出版社', cover_image: 'https://picsum.photos/seed/santi456/400/600', description: '中国科幻巅峰之作，亚洲首部雨果奖长篇小说。', category: '科技', mode: 2, status: 1, daily_rent: 5, weekly_rent: 35, deposit: 7, sell_price: 68, min_rent_days: 1, images: ['https://picsum.photos/seed/santi456/400/600'], pickup_location: '清华大学北门', view_count: 52340, created_at: '2024-01-16T10:00:00Z', updated_at: '2024-01-16T10:00:00Z' },
    { id: 3, user_id: 1, isbn: '978-7530211533', title: '平凡的世界', author: '路遥', publisher: '北京十月文艺出版社', cover_image: 'https://picsum.photos/seed/pingfan789/400/600', description: '激励无数青年的现实主义巨著。', category: '文学', mode: 1, status: 1, daily_rent: 4, weekly_rent: 25, deposit: 5, sell_price: 45, min_rent_days: 1, images: ['https://picsum.photos/seed/pingfan789/400/600'], pickup_location: '北京大学东门', view_count: 38920, created_at: '2024-01-17T10:00:00Z', updated_at: '2024-01-17T10:00:00Z' },
    { id: 4, user_id: 2, isbn: '978-7544263926', title: '我的阿勒泰', author: '李娟', publisher: '云南人民出版社', cover_image: 'https://picsum.photos/seed/alatai123/400/600', description: '2024年京东图书畅销榜冠军。', category: '文学', mode: 2, status: 1, daily_rent: 3.5, weekly_rent: 20, deposit: 5, sell_price: 38, min_rent_days: 1, images: ['https://picsum.photos/seed/alatai123/400/600'], pickup_location: '清华大学北门', view_count: 41230, created_at: '2024-01-18T10:00:00Z', updated_at: '2024-01-18T10:00:00Z' },
    { id: 5, user_id: 3, isbn: '978-7530211535', title: '我与地坛', author: '史铁生', publisher: '人民文学出版社', cover_image: 'https://picsum.photos/seed/ditan321/400/600', description: '史铁生的代表作，关于生命和母爱的深刻思考。', category: '文学', mode: 1, status: 1, daily_rent: 3, weekly_rent: 18, deposit: 5, sell_price: 32, min_rent_days: 1, images: ['https://picsum.photos/seed/ditan321/400/600'], pickup_location: '复旦大学邯郸校区', view_count: 29870, created_at: '2024-01-19T10:00:00Z', updated_at: '2024-01-19T10:00:00Z' },
    { id: 6, user_id: 3, isbn: '978-7020047540', title: '红楼梦', author: '曹雪芹', publisher: '人民文学出版社', cover_image: 'https://picsum.photos/seed/hongloumeng/400/600', description: '中国古典四大名著之首。', category: '文学', mode: 1, status: 1, daily_rent: 4, weekly_rent: 22, deposit: 5, sell_price: 48, min_rent_days: 1, images: ['https://picsum.photos/seed/hongloumeng/400/600'], pickup_location: '复旦大学邯郸校区', view_count: 35670, created_at: '2024-01-20T10:00:00Z', updated_at: '2024-01-20T10:00:00Z' }
  ];

  return { users, wallets, books, rentals: [], sells: [], gifts: [], transactions: [], favorites: [], credit_records: [], nextIds: { user: 4, book: 7, rental: 1, sell: 1, gift: 1, wallet: 4, transaction: 1, favorite: 1 } };
}

let currentUser = null;

function setToken(token) {
  if (token) {
    localStorage.setItem('token', token);
    try {
      const payload = JSON.parse(atob(token));
      currentUser = getStoredData().users.find(u => u.id === payload.id);
    } catch (e) {
      currentUser = null;
    }
  } else {
    localStorage.removeItem('token');
    localStorage.removeItem('currentUser');
    currentUser = null;
  }
}

function getToken() {
  return localStorage.getItem('token');
}

function delay(ms = 100) {
  return new Promise(resolve => setTimeout(resolve, ms));
}

function getCurrentUserId() {
  const userStr = localStorage.getItem('currentUser');
  if (!userStr) return null;
  try {
    const user = JSON.parse(userStr);
    return user.id;
  } catch {
    return null;
  }
}

async function mockRequest(endpoint, options = {}) {
  await delay(100);
  const method = options.method || 'GET';
  const body = options.body ? JSON.parse(options.body) : null;

  // 分离路径和查询参数
  const [path, queryString] = endpoint.split('?');
  const parts = path.split('/').filter(Boolean);
  const queryParams = {};
  if (queryString) {
    queryString.split('&').forEach(param => {
      const [key, value] = param.split('=');
      if (key) queryParams[decodeURIComponent(key)] = decodeURIComponent(value || '');
    });
  }

  // Users
  if (parts[0] === 'users') {
    if (parts[1] === 'login' && method === 'POST') {
      const data = getStoredData();
      const user = data.users.find(u => u.account === body.account && u.password === body.password);
      if (!user) throw new Error('账号或密码错误');
      const token = btoa(JSON.stringify({ id: user.id, account: user.account, exp: Date.now() + 7 * 24 * 60 * 60 * 1000 }));
      localStorage.setItem('token', token);
      localStorage.setItem('currentUser', JSON.stringify(user));
      return { data: { token, user } };
    }
    if (parts[1] === 'register' && method === 'POST') {
      const data = getStoredData();
      if (data.users.find(u => u.account === body.account)) throw new Error('账号已存在');
      const newUser = { id: data.nextIds.user++, account: body.account, password: body.password, birth_date: body.birth_date, nickname: '新用户', school: '', campus: '', grade: '', credit_score: 100, status: 1, created_at: new Date().toISOString(), updated_at: new Date().toISOString() };
      data.users.push(newUser);
      data.wallets.push({ id: data.nextIds.wallet++, user_id: newUser.id, balance: 1000, frozen_balance: 0, total_income: 0, total_expense: 0 });
      saveData(data);
      return { data: { id: newUser.id } };
    }
    if (parts[1] === 'me' && method === 'GET') {
      const userStr = localStorage.getItem('currentUser');
      if (!userStr) throw new Error('请先登录');
      return { data: JSON.parse(userStr) };
    }
    if (parts[1] === 'me' && method === 'PUT') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const data = getStoredData();
      const idx = data.users.findIndex(u => u.id === userId);
      if (idx >= 0) {
        data.users[idx] = { ...data.users[idx], ...body, updated_at: new Date().toISOString() };
        localStorage.setItem('currentUser', JSON.stringify(data.users[idx]));
        saveData(data);
        return { data: data.users[idx] };
      }
      return { data: null };
    }
    if (parts[1] === 'logout' && method === 'POST') {
      setToken(null);
      return { data: null };
    }
  }

  // Books
  if (parts[0] === 'books') {
    const data = getStoredData();
    if (!parts[1]) {
      if (method === 'GET') {
        let filtered = data.books.filter(b => b.status === 1);
        const page = parseInt(queryParams.page || 1);
        const pageSize = parseInt(queryParams.page_size || 20);
        const mode = queryParams.mode;
        const category = queryParams.category;
        const keyword = queryParams.keyword;
        if (mode) filtered = filtered.filter(b => b.mode === parseInt(mode));
        if (category) filtered = filtered.filter(b => b.category === category);
        if (keyword) { const kw = keyword.toLowerCase(); filtered = filtered.filter(b => b.title.toLowerCase().includes(kw) || b.author.toLowerCase().includes(kw)); }
        const start = (page - 1) * pageSize;
        return { data: { list: filtered.slice(start, start + pageSize).map(b => ({ ...b, user: data.users.find(u => u.id === b.user_id) })), total: filtered.length, page, page_size: pageSize } };
      }
      if (method === 'POST') {
        if (!getToken()) throw new Error('请先登录');
        const userId = getCurrentUserId();
        const newBook = { id: data.nextIds.book++, user_id: userId, ...body, status: 1, view_count: 0, created_at: new Date().toISOString(), updated_at: new Date().toISOString() };
        data.books.push(newBook);
        saveData(data);
        return { data: newBook };
      }
    }
    if (parts[1] === 'hot' && method === 'GET') {
      return { data: data.books.filter(b => b.status === 1).sort((a, b) => b.view_count - a.view_count).slice(0, 10) };
    }
    if (parts[1] === 'gift' && parts[2] === 'latest' && method === 'GET') {
      return { data: data.books.filter(b => b.status === 1 && b.mode === 3).slice(0, 6) };
    }
    if (parts[1] === 'my' && method === 'GET') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const myBooks = data.books.filter(b => b.user_id === userId);
      return { data: { list: myBooks, total: myBooks.length } };
    }
    if (parts[1] === 'favorites' && method === 'GET') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const favs = data.favorites.filter(f => f.user_id === userId);
      const favBooks = favs.map(f => data.books.find(b => b.id === f.book_id)).filter(Boolean).map(b => ({ ...b, user: data.users.find(u => u.id === b.user_id) }));
      return { data: { list: favBooks, total: favBooks.length } };
    }
    const bookId = parseInt(parts[1]);
    if (bookId && method === 'GET') {
      const book = data.books.find(b => b.id === bookId);
      if (!book) throw new Error('图书不存在');
      return { data: { ...book, user: data.users.find(u => u.id === book.user_id) } };
    }
  }

  // Rentals
  if (parts[0] === 'rentals') {
    const data = getStoredData();
    if (parts[1] === 'my' && method === 'GET') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const myRentals = data.rentals.filter(r => r.renter_id === userId || r.owner_id === userId);
      return { data: { list: myRentals, total: myRentals.length } };
    }
    const rentalId = parseInt(parts[1]);
    if (rentalId && method === 'GET') {
      const rental = data.rentals.find(r => r.id === rentalId);
      if (!rental) return { data: null };
      return { data: rental };
    }
    if (method === 'POST') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const newRental = { id: data.nextIds.rental++, renter_id: userId, ...body, status: 1, created_at: new Date().toISOString() };
      data.rentals.push(newRental);
      saveData(data);
      return { data: newRental };
    }
  }

  // Sells
  if (parts[0] === 'sells') {
    const data = getStoredData();
    if (parts[1] === 'my' && method === 'GET') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const mySells = data.sells.filter(s => s.seller_id === userId || s.buyer_id === userId);
      return { data: { list: mySells, total: mySells.length } };
    }
    const sellId = parseInt(parts[1]);
    if (sellId && method === 'GET') {
      const sell = data.sells.find(s => s.id === sellId);
      if (!sell) return { data: null };
      return { data: sell };
    }
    if (method === 'POST') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const newSell = { id: data.nextIds.sell++, seller_id: userId, ...body, status: 1, created_at: new Date().toISOString() };
      data.sells.push(newSell);
      saveData(data);
      return { data: newSell };
    }
  }

  // Gifts
  if (parts[0] === 'gifts') {
    const data = getStoredData();
    if (parts[1] === 'my' && method === 'GET') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const myGifts = data.gifts.filter(g => g.sender_id === userId || g.receiver_id === userId);
      return { data: { list: myGifts, total: myGifts.length } };
    }
    const giftId = parseInt(parts[1]);
    if (giftId && method === 'GET') {
      const gift = data.gifts.find(g => g.id === giftId);
      if (!gift) return { data: null };
      return { data: gift };
    }
    if (method === 'POST') {
      if (!getToken()) throw new Error('请先登录');
      const userId = getCurrentUserId();
      const newGift = { id: data.nextIds.gift++, sender_id: userId, ...body, status: 1, created_at: new Date().toISOString() };
      data.gifts.push(newGift);
      saveData(data);
      return { data: newGift };
    }
  }

  // Wallets
  if (parts[0] === 'wallets' && method === 'GET') {
    if (!getToken()) throw new Error('请先登录');
    const userId = getCurrentUserId();
    const data = getStoredData();
    const wallet = data.wallets.find(w => w.user_id === userId);
    return { data: wallet || { balance: 1000, frozen_balance: 0 } };
  }

  // Guide
  if (parts[0] === 'guide' && method === 'GET') {
    return { data: { content: '<h2>使用指南</h2><p>欢迎使用高校图书流转平台！</p>' } };
  }

  // Fallback: return empty success
  return { data: null };
}

window.mockRequest = mockRequest;
window.setToken = setToken;
window.getToken = getToken;
