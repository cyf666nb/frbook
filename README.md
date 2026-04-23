# 全市高校二手图书共享流转平台

基于 Go + Gin + MySQL + Redis 的高校二手图书共享平台后端服务。

## 功能特性

- **用户模块**: 简化注册(年龄<30岁验证)、JWT认证、资料补全、信用评分
- **图书发布**: ISBN扫码填充、三种流转模式(租赁/出售/赠送)
- **租赁业务**: 分布式锁、押金担保、租期计时、逾期扣费
- **出售业务**: 货款担保、7天自动确认、平台抽成
- **赠送业务**: 信息撮合、申请确认流程、双方互评

## 项目结构

```
bookshare/
├── cmd/server/          # 主入口
├── config/              # 配置文件
├── docs/                # 设计文档
│   ├── sql/            # 数据库表结构
│   ├── api-design.md   # API接口设计
│   ├── redis-design.md # Redis数据结构
│   ├── state-machine.md # 状态机设计
│   └── high-concurrency.md # 高并发注意事项
├── internal/           # 内部包
│   ├── config/         # 配置管理
│   ├── database/       # 数据库连接
│   ├── handler/        # HTTP处理器
│   ├── middleware/     # 中间件
│   ├── model/          # 数据模型
│   ├── repository/     # 数据仓储
│   ├── router/         # 路由配置
│   └── service/         # 业务逻辑
└── pkg/response/        # 统一响应
```

## 快速开始

### 方式一: Docker Compose (本地开发)

```bash
# 启动本地MySQL + Redis + API
docker-compose -f docker-compose.local.yml up -d
```

### 方式二: 连接云端数据库

1. 复制环境变量模板并配置:
```bash
cp .env.example .env
# 编辑.env,填入您的云端数据库地址
```

2. 使用云端数据库启动:
```bash
# 设置环境变量后运行
docker-compose up -d

# 或直接指定环境变量
DB_HOST=your-mysql-host \
REDIS_HOST=your-redis-host \
docker-compose up -d
```

### 方式三: 本地运行

```bash
# 1. 确保MySQL和Redis已启动
# 2. 导入数据库结构
mysql -u root -p < docs/sql/schema.sql

# 3. 设置环境变量
export DB_HOST=localhost
export DB_PORT=3306
export DB_USER=root
export DB_PASSWORD=your_password
export DB_NAME=bookshare
export REDIS_HOST=localhost
export REDIS_PORT=6379

# 4. 运行服务
go run cmd/server/main.go
```

## 环境变量说明

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

## API文档

启动服务后访问: `http://localhost:8080/health`

详细API文档请查看 [docs/api-design.md](docs/api-design.md)

## 数据库表结构

请查看 [docs/sql/schema.sql](docs/sql/schema.sql)

## 主要技术栈

- Go 1.21+
- Gin Web Framework
- GORM (ORM)
- MySQL 8.0
- Redis 7
- JWT Authentication
- Docker / Docker Compose
