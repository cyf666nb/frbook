# frbook - 高校二手图书共享流转平台

基于 Go + Gin + Vue + MySQL + Redis 的高校二手图书共享平台。

## 功能特性

- **用户模块**: 年龄<30岁验证、JWT认证、资料补全、信用评分
- **图书发布**: ISBN扫码填充、三种流转模式（租赁/出售/赠送）
- **租赁业务**: 分布式锁、押金担保、租期计时、逾期扣费
- **出售业务**: 货款担保、7天自动确认、平台抽成
- **赠送业务**: 信息撮合、申请确认流程、双方互评

## 项目结构

```
frbook/
├── cmd/
│   ├── server/          # 主入口
│   └── seed/            # 数据初始化
├── config/              # 配置文件
├── docs/                # 设计文档
│   ├── sql/            # 数据库表结构
│   ├── api-design.md   # API接口设计
│   ├── redis-design.md # Redis数据结构
│   ├── state-machine.md # 状态机设计
│   └── high-concurrency.md # 高并发注意事项
├── frontend/           # Vue 前端
│   └── src/
│       ├── api/        # API调用
│       ├── components/  # 公共组件
│       ├── context/    # React Context
│       └── pages/      # 页面组件
├── internal/           # 内部包
│   ├── config/        # 配置管理
│   ├── database/       # 数据库连接
│   ├── handler/       # HTTP处理器
│   ├── middleware/    # 中间件
│   ├── model/         # 数据模型
│   ├── repository/    # 数据仓储
│   ├── router/        # 路由配置
│   └── service/       # 业务逻辑
└── pkg/response/       # 统一响应
```

## 技术栈

- **后端**: Go 1.21+ / Gin / GORM
- **前端**: Vue 3 / Vite / Tailwind CSS
- **数据库**: MySQL 8.0 / Redis 7
- **认证**: JWT
- **部署**: Docker / Docker Compose

## 快速开始

### Docker Compose

```bash
# 启动 MySQL + Redis + API
docker-compose up -d
```

### 本地运行

**后端：**

```bash
# 1. 确保 MySQL 和 Redis 已启动

# 2. 导入数据库结构
mysql -u root -p < docs/sql/schema.sql

# 3. 配置环境变量
cp .env.example .env
# 编辑 .env 填入数据库信息

# 4. 运行
go run cmd/server/main.go
```

**前端：**

```bash
cd frontend
npm install
npm run dev
```

## 环境变量

| 变量名 | 必填 | 默认值 | 说明 |
|--------|------|--------|------|
| DB_HOST | 是 | - | MySQL主机地址 |
| DB_PORT | 否 | 3306 | MySQL端口 |
| DB_USER | 是 | - | 数据库用户名 |
| DB_PASSWORD | 是 | - | 数据库密码 |
| DB_NAME | 否 | bookshare | 数据库名 |
| REDIS_HOST | 是 | - | Redis主机地址 |
| REDIS_PORT | 否 | 6379 | Redis端口 |
| REDIS_PASSWORD | 否 | - | Redis密码 |
| JWT_SECRET | 否 | - | JWT签名密钥 |

## API 文档

启动服务后访问: `http://localhost:8080/health`

详细API文档请查看 [docs/api-design.md](docs/api-design.md)

数据库表结构请查看 [docs/sql/schema.sql](docs/sql/schema.sql)
