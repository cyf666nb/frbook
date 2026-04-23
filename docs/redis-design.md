# Redis 数据结构设计

## 一、Key 命名规范

```
{业务模块}:{具体功能}:{唯一标识}[:子标识]
```

### 命名规范说明
- 使用冒号 `:` 分隔层级
- 使用小写字母和下划线
- 避免特殊字符
- 保持简洁明了

---

## 二、用户模块

### 2.1 Session/Token 存储

```
Key: user:session:{user_id}
Type: String (JSON)
Value: {
    "token": "xxx",
    "device": "web/ios/android",
    "ip": "192.168.1.1",
    "login_time": 1700000000,
    "expire_time": 1700086400
}
TTL: 7天 (604800秒)

Key: user:token:{token}
Type: String
Value: {user_id}
TTL: 7天
```

### 2.2 Token 黑名单

```
Key: user:token:blacklist:{token}
Type: String
Value: 1
TTL: 7天 (与Token有效期一致)
```

### 2.3 用户信息缓存

```
Key: user:info:{user_id}
Type: Hash
Fields:
    - nickname: "用户昵称"
    - avatar: "头像URL"
    - school: "学校"
    - credit_score: "100"
TTL: 1小时
```

### 2.4 验证码存储

```
Key: user:verify:{type}:{target}
Type: String
Value: {code}
TTL: 5分钟 (300秒)

示例:
user:verify:register:13800138000 -> "123456"
user:verify:login:13800138000 -> "654321"
```

---

## 三、图书模块

### 3.1 图书详情缓存

```
Key: book:detail:{book_id}
Type: Hash
Fields:
    - id: "1"
    - title: "书名"
    - author: "作者"
    - cover_image: "封面URL"
    - mode: "1"
    - status: "1"
    - daily_rent: "5.00"
    - sell_price: "50.00"
TTL: 30分钟
```

### 3.2 首页热门图书缓存

```
Key: book:hot:list
Type: List
Value: [book_id1, book_id2, ...]
TTL: 10分钟
```

### 3.3 首页最新赠送图书列表

```
Key: book:gift:latest
Type: List
Value: [book_id1, book_id2, ...]
TTL: 5分钟
```

### 3.4 高校排行缓存

```
Key: book:rank:school:{school_name}
Type: Sorted Set
Score: 交易数量
Member: book_id
TTL: 1小时

Key: book:rank:global
Type: Sorted Set
Score: 浏览量
Member: book_id
TTL: 1小时
```

### 3.5 图书浏览计数

```
Key: book:view:{book_id}
Type: String (Integer)
Value: 浏览次数
TTL: 永久

Key: book:view:daily:{date}
Type: Sorted Set
Score: 浏览次数
Member: book_id
TTL: 7天
```

---

## 四、租赁模块

### 4.1 库存锁定分布式锁

```
Key: lock:rental:book:{book_id}
Type: String
Value: {order_id}
TTL: 10分钟

使用 SETNX 实现:
SET lock:rental:book:123 "order_456" NX EX 600
```

### 4.2 租期倒计时任务

```
Key: rental:countdown:{order_id}
Type: String
Value: {end_timestamp}
TTL: 租赁剩余时间

使用 Redis Keyspace Notifications 监听过期事件
或使用 Redisson RDelayedQueue
```

### 4.3 租赁订单缓存

```
Key: rental:order:{order_id}
Type: Hash
Fields:
    - order_no: "RN202311010001"
    - book_id: "123"
    - status: "2"
    - end_time: "1700086400"
TTL: 1小时
```

### 4.4 用户租赁中的订单列表

```
Key: rental:user:active:{user_id}
Type: Set
Value: {order_id1, order_id2, ...}
TTL: 30分钟
```

---

## 五、出售模块

### 5.1 库存锁定分布式锁

```
Key: lock:sell:book:{book_id}
Type: String
Value: {order_id}
TTL: 15分钟
```

### 5.2 订单超时取消延迟队列

```
Key: sell:delay:queue
Type: Sorted Set
Score: 过期时间戳
Member: {order_id}

使用 Redisson RDelayedQueue 或 Keyspace Notifications
```

### 5.3 自动结算延迟任务

```
Key: sell:settle:delay:{order_id}
Type: String
Value: {auto_settle_timestamp}
TTL: 7天
```

---

## 六、赠送模块

### 6.1 最新赠送列表缓存

```
Key: gift:latest:list
Type: List
Value: [book_id1, book_id2, ...]
Length: 最多100条
TTL: 5分钟
```

### 6.2 赠送请求锁定

```
Key: lock:gift:book:{book_id}
Type: String
Value: {record_id}
TTL: 5分钟
```

---

## 七、统计与计数

### 7.1 用户发布计数

```
Key: stats:user:publish:{user_id}
Type: Hash
Fields:
    - total: "10"
    - rental: "5"
    - sell: "3"
    - gift: "2"
TTL: 1天
```

### 7.2 用户交易计数

```
Key: stats:user:trade:{user_id}
Type: Hash
Fields:
    - rental_as_owner: "5"
    - rental_as_renter: "3"
    - sell_as_seller: "2"
    - sell_as_buyer: "1"
    - gift_as_giver: "3"
    - gift_as_receiver: "2"
TTL: 1天
```

### 7.3 平台总览统计

```
Key: stats:platform:daily:{date}
Type: Hash
Fields:
    - new_users: "100"
    - new_books: "50"
    - rental_orders: "20"
    - sell_orders: "15"
    - gift_records: "10"
TTL: 30天
```

---

## 八、限流与防刷

### 8.1 API 请求限流

```
Key: limit:api:{user_id}:{endpoint}
Type: String (计数器)
Value: 请求次数
TTL: 1分钟

示例: limit:api:123:/api/v1/books -> "10"
```

### 8.2 短信发送限流

```
Key: limit:sms:{phone}
Type: String
Value: 发送次数
TTL: 1小时

Key: limit:sms:daily:{phone}
Type: String
Value: 当日发送次数
TTL: 1天
```

---

## 九、搜索与索引

### 9.1 图书搜索缓存

```
Key: search:book:{keyword_hash}
Type: List
Value: [book_id1, book_id2, ...]
TTL: 10分钟
```

### 9.2 用户搜索缓存

```
Key: search:user:{keyword_hash}
Type: List
Value: [user_id1, user_id2, ...]
TTL: 10分钟
```

---

## 十、消息队列 (替代专业MQ的简单实现)

### 10.1 订单状态变更队列

```
Key: mq:order:status
Type: List (LPUSH/RPOP)
Value: JSON {
    "order_type": "rental",
    "order_id": 123,
    "old_status": 1,
    "new_status": 2,
    "timestamp": 1700000000
}
```

### 10.2 通知消息队列

```
Key: mq:notification:{user_id}
Type: List (LPUSH/RPOP)
Value: JSON {
    "type": "order_status",
    "title": "订单状态变更",
    "content": "您的租赁订单已确认",
    "data": {...}
}
```

---

## 十一、Redis 配置建议

### 11.1 内存配置
```
maxmemory 2gb
maxmemory-policy allkeys-lru
```

### 11.2 持久化配置
```
save 900 1
save 300 10
save 60 10000
appendonly yes
appendfsync everysec
```

### 11.3 Keyspace Notifications (用于过期监听)
```
notify-keyspace-events Ex
```
