# 核心 API 接口清单

## 基础信息
- 基础路径: `/api/v1`
- 认证方式: Bearer Token (Header: Authorization)
- 响应格式: JSON

---

## 一、用户模块 `/api/v1/users`

### 1.1 用户注册
```
POST /api/v1/users/register
```
**请求体:**
```json
{
    "account": "13800138000",
    "password": "password123",
    "birth_date": "2000-01-01",
    "code": "123456"
}
```
**功能:** 注册新用户,验证年龄小于30岁

### 1.2 用户登录
```
POST /api/v1/users/login
```
**请求体:**
```json
{
    "account": "13800138000",
    "password": "password123"
}
```
**功能:** 用户登录,返回Token

### 1.3 用户登出
```
POST /api/v1/users/logout
```
**功能:** 登出,将Token加入黑名单

### 1.4 获取当前用户信息
```
GET /api/v1/users/me
```
**功能:** 获取当前登录用户详细信息

### 1.5 更新用户资料
```
PUT /api/v1/users/me
```
**请求体:**
```json
{
    "nickname": "昵称",
    "avatar": "http://...",
    "school": "北京大学",
    "campus": "海淀校区",
    "grade": "大三",
    "interest_tags": ["文学", "科技"]
}
```
**功能:** 更新用户资料(学校、校区、兴趣标签等)

### 1.6 发送验证码
```
POST /api/v1/users/verify-code
```
**请求体:**
```json
{
    "target": "13800138000",
    "type": 1
}
```
**功能:** 发送短信/邮件验证码(type: 1-注册 2-登录 3-重置密码)

### 1.7 获取用户信用记录
```
GET /api/v1/users/me/credit-records
```
**功能:** 获取当前用户信用分变更记录

### 1.8 获取用户钱包信息
```
GET /api/v1/users/me/wallet
```
**功能:** 获取钱包余额、冻结金额等信息

### 1.9 获取其他用户公开信息
```
GET /api/v1/users/{user_id}
```
**功能:** 获取其他用户的公开信息(昵称、学校、信用分等)

---

## 二、图书模块 `/api/v1/books`

### 2.1 发布图书
```
POST /api/v1/books
```
**请求体:**
```json
{
    "isbn": "9787111111111",
    "title": "书名",
    "author": "作者",
    "publisher": "出版社",
    "cover_image": "http://...",
    "description": "描述",
    "category": "文学",
    "mode": 1,
    "daily_rent": 5.00,
    "weekly_rent": 30.00,
    "deposit": 50.00,
    "min_rent_days": 3,
    "images": ["http://...1", "http://...2"],
    "pickup_location": "北京大学东门"
}
```
**功能:** 发布图书(租赁/出售/赠送)

### 2.2 ISBN查询图书信息
```
GET /api/v1/books/isbn/{isbn}
```
**功能:** 通过ISBN调用豆瓣API获取图书信息

### 2.3 获取图书详情
```
GET /api/v1/books/{book_id}
```
**功能:** 获取图书详细信息

### 2.4 更新图书信息
```
PUT /api/v1/books/{book_id}
```
**功能:** 更新图书信息(仅发布者可操作)

### 2.5 删除/下架图书
```
DELETE /api/v1/books/{book_id}
```
**功能:** 下架图书

### 2.6 上架图书
```
PUT /api/v1/books/{book_id}/online
```
**功能:** 重新上架已下架的图书

### 2.7 获取图书列表
```
GET /api/v1/books
```
**查询参数:**
- `mode`: 流转模式(1-租赁 2-出售 3-赠送)
- `category`: 分类
- `keyword`: 搜索关键词
- `school`: 学校筛选
- `sort`: 排序字段(created_at/view_count/price)
- `order`: 排序方式(asc/desc)
- `page`: 页码
- `page_size`: 每页数量

**功能:** 分页获取图书列表

### 2.8 获取首页推荐
```
GET /api/v1/books/recommend
```
**功能:** 获取首页推荐图书(基于用户兴趣标签和学校)

### 2.9 获取热门图书
```
GET /api/v1/books/hot
```
**功能:** 获取热门图书列表

### 2.10 获取最新赠送图书
```
GET /api/v1/books/gift/latest
```
**功能:** 获取最新赠送图书列表

### 2.11 收藏图书
```
POST /api/v1/books/{book_id}/favorite
```
**功能:** 收藏图书

### 2.12 取消收藏
```
DELETE /api/v1/books/{book_id}/favorite
```
**功能:** 取消收藏

### 2.13 获取收藏列表
```
GET /api/v1/books/favorites
```
**功能:** 获取当前用户收藏的图书列表

### 2.14 获取我发布的图书
```
GET /api/v1/books/my
```
**功能:** 获取当前用户发布的图书列表

---

## 三、租赁模块 `/api/v1/rentals`

### 3.1 提交租赁申请
```
POST /api/v1/rentals
```
**请求体:**
```json
{
    "book_id": 1,
    "rent_days": 7
}
```
**功能:** 提交租赁申请,创建订单,锁定库存

### 3.2 支付租赁订单
```
POST /api/v1/rentals/{order_id}/pay
```
**请求体:**
```json
{
    "payment_method": "wallet"
}
```
**功能:** 支付押金+首期租金

### 3.3 取消租赁订单
```
POST /api/v1/rentals/{order_id}/cancel
```
**功能:** 取消租赁订单(待支付状态)

### 3.4 出租人确认取书
```
POST /api/v1/rentals/{order_id}/confirm-pickup
```
**功能:** 出租人确认承租人已取书,开始计时

### 3.5 承租人归还图书
```
POST /api/v1/rentals/{order_id}/return
```
**功能:** 承租人发起归还

### 3.6 出租人验收
```
POST /api/v1/rentals/{order_id}/inspect
```
**请求体:**
```json
{
    "passed": true,
    "remark": "图书完好"
}
```
**功能:** 出租人验收归还的图书

### 3.7 获取租赁订单详情
```
GET /api/v1/rentals/{order_id}
```
**功能:** 获取订单详细信息

### 3.8 获取我的租赁订单列表
```
GET /api/v1/rentals/my
```
**查询参数:**
- `role`: 角色(owner-出租方/renter-承租方)
- `status`: 订单状态
- `page`: 页码
- `page_size`: 每页数量

**功能:** 获取当前用户的租赁订单列表

### 3.9 评价租赁订单
```
POST /api/v1/rentals/{order_id}/rate
```
**请求体:**
```json
{
    "rating": 5,
    "comment": "很好的交易体验"
}
```
**功能:** 对交易对方进行评价

### 3.10 获取逾期订单
```
GET /api/v1/rentals/overdue
```
**功能:** 获取当前逾期的租赁订单(管理员或当前用户)

---

## 四、出售模块 `/api/v1/sells`

### 4.1 购买图书
```
POST /api/v1/sells
```
**请求体:**
```json
{
    "book_id": 1
}
```
**功能:** 拍下图书,创建出售订单

### 4.2 支付订单
```
POST /api/v1/sells/{order_id}/pay
```
**请求体:**
```json
{
    "payment_method": "wallet"
}
```
**功能:** 支付货款到担保账户

### 4.3 取消订单
```
POST /api/v1/sells/{order_id}/cancel
```
**功能:** 取消订单(待支付状态)

### 4.4 卖家发货
```
POST /api/v1/sells/{order_id}/ship
```
**请求体:**
```json
{
    "delivery_type": 1,
    "delivery_company": "顺丰",
    "delivery_no": "SF1234567890"
}
```
**功能:** 卖家发货,上传快递信息

### 4.5 线下自提确认
```
POST /api/v1/sells/{order_id}/pickup
```
**功能:** 确认线下自提已交付

### 4.6 买家确认收货
```
POST /api/v1/sells/{order_id}/confirm
```
**功能:** 买家确认收货,触发结算

### 4.7 申请退款
```
POST /api/v1/sells/{order_id}/refund
```
**请求体:**
```json
{
    "reason": "图书与描述不符"
}
```
**功能:** 买家申请退款

### 4.8 获取订单详情
```
GET /api/v1/sells/{order_id}
```
**功能:** 获取出售订单详情

### 4.9 获取我的出售订单
```
GET /api/v1/sells/my
```
**查询参数:**
- `role`: 角色(seller-卖家/buyer-买家)
- `status`: 订单状态
- `page`: 页码
- `page_size`: 每页数量

**功能:** 获取当前用户的出售订单列表

### 4.10 评价订单
```
POST /api/v1/sells/{order_id}/rate
```
**请求体:**
```json
{
    "rating": 5,
    "comment": "图书质量很好"
}
```
**功能:** 对交易对方进行评价

---

## 五、赠送模块 `/api/v1/gifts`

### 5.1 申请领取
```
POST /api/v1/gifts
```
**请求体:**
```json
{
    "book_id": 1,
    "message": "我很需要这本书..."
}
```
**功能:** 申请领取赠送图书

### 5.2 赠送人确认
```
POST /api/v1/gifts/{record_id}/confirm
```
**功能:** 赠送人确认将书送给该申请人

### 5.3 赠送人拒绝
```
POST /api/v1/gifts/{record_id}/reject
```
**功能:** 赠送人拒绝申请

### 5.4 取消申请
```
POST /api/v1/gifts/{record_id}/cancel
```
**功能:** 申请人取消申请

### 5.5 确认交付
```
POST /api/v1/gifts/{record_id}/deliver
```
**功能:** 确认图书已交付

### 5.6 获取赠送记录详情
```
GET /api/v1/gifts/{record_id}
```
**功能:** 获取赠送记录详情

### 5.7 获取我的赠送记录
```
GET /api/v1/gifts/my
```
**查询参数:**
- `role`: 角色(giver-赠送方/receiver-接收方)
- `status`: 状态
- `page`: 页码
- `page_size`: 每页数量

**功能:** 获取当前用户的赠送记录

### 5.8 获取图书的申请列表
```
GET /api/v1/gifts/book/{book_id}/applications
```
**功能:** 获取某本赠送图书的申请列表(仅赠送人可见)

### 5.9 评价
```
POST /api/v1/gifts/{record_id}/rate
```
**请求体:**
```json
{
    "rating": 5,
    "comment": "感谢赠送!"
}
```
**功能:** 双方互评

---

## 六、统计模块 `/api/v1/stats`

### 6.1 获取平台统计
```
GET /api/v1/stats/platform
```
**功能:** 获取平台整体统计数据(管理员)

### 6.2 获取用户统计
```
GET /api/v1/stats/user/{user_id}
```
**功能:** 获取用户交易统计数据

### 6.3 获取学校排行
```
GET /api/v1/stats/rank/school
```
**功能:** 获取各学校交易量排行

### 6.4 获取图书排行
```
GET /api/v1/stats/rank/books
```
**功能:** 获取热门图书排行

---

## 七、公共模块 `/api/v1/common`

### 7.1 获取学校列表
```
GET /api/v1/common/schools
```
**功能:** 获取城市内高校列表

### 7.2 获取分类列表
```
GET /api/v1/common/categories
```
**功能:** 获取图书分类列表

### 7.3 上传图片
```
POST /api/v1/common/upload
```
**请求体:** multipart/form-data
**功能:** 上传图片,返回图片URL

---

## 响应格式

### 成功响应
```json
{
    "code": 0,
    "message": "success",
    "data": { ... }
}
```

### 错误响应
```json
{
    "code": 1001,
    "message": "参数错误",
    "data": null
}
```

### 分页响应
```json
{
    "code": 0,
    "message": "success",
    "data": {
        "list": [...],
        "total": 100,
        "page": 1,
        "page_size": 20
    }
}
```

---

## 错误码定义

| 错误码 | 说明 |
|--------|------|
| 0 | 成功 |
| 1001 | 参数错误 |
| 1002 | 未授权 |
| 1003 | 禁止访问 |
| 1004 | 资源不存在 |
| 2001 | 用户已存在 |
| 2002 | 用户不存在 |
| 2003 | 密码错误 |
| 2004 | 验证码错误 |
| 2005 | 年龄不符合要求 |
| 3001 | 图书不存在 |
| 3002 | 图书已下架 |
| 3003 | 图书已被锁定 |
| 4001 | 订单不存在 |
| 4002 | 订单状态错误 |
| 4003 | 余额不足 |
| 4004 | 库存锁定失败 |
| 5001 | 支付失败 |
| 5002 | 退款失败 |
