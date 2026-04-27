# frbook - 高校二手图书共享流转平台

一个面向大学生的二手书籍共享/租赁平台，使用 React + localStorage 实现，纯前端版本，无需后端。

> 💡 如需完整版本（包含后端），请访问：[frbook 完整版](https://github.com/cyf666nb/frbook/tree/main)

## 功能特点

- **用户模块** - 注册登录、个人资料管理、信用评分
- **图书浏览** - 分类筛选、搜索、热门图书
- **租赁功能** - 图书租赁、订单管理
- **购买功能** - 图书购买、订单追踪
- **赠送功能** - 图书赠送、信息撮合

## 技术栈

- React 18
- Vite
- Tailwind CSS
- localStorage (数据持久化)

## 快速开始

```bash
cd frontend
npm install
npm run dev
```

访问 http://localhost:5173

## 测试账号

| 账号 | 密码 | 学校 |
|------|------|------|
| 13800138000 | 123456 | 北京大学 |
| testuser | 123456 | 清华大学 |
| 13599887766 | 123456 | 复旦大学 |

## 项目结构

```
frontend/
├── src/
│   ├── api/          # API 接口（包含 mock 数据）
│   ├── components/  # 公共组件
│   ├── context/      # React Context
│   └── pages/       # 页面组件
├── public/           # 静态资源
└── dist/             # 构建产物
```

## 数据存储

使用浏览器 localStorage 存储数据，包括：
- 用户信息
- 图书数据
- 订单记录
- 钱包余额

如需重置数据，在浏览器控制台执行：
```javascript
localStorage.clear()
window.location.reload()
```
