# Gin-Vue-Chat 即时聊天应用
## 技术栈
### 后端
- Golang
- Gin框架
- Gorm
- MySQL
- WebSocket

## 项目结构

```
Gin-Vue-Chat/
├── frontend/           # 前端Vue项目
└── backend/            # 后端Golang项目
```

该应用已经实现了所有基本功能：
- 用户注册和登录
- 添加好友和私聊
- 创建群组和群聊
- 实时消息推送

## 开发指南

### 前端开发

```bash
cd frontend
pnpm install
pnpm run dev
```

### 后端开发

```bash
cd backend
go mod tidy
go run main.go
```

## 后端代码学习查看顺序

为了更好地理解后端代码结构和实现逻辑，建议按照以下顺序进行学习：

1. **入口文件**
   - [main.go](backend/main.go) - 程序入口点，包含路由配置和服务器启动逻辑

2. **配置模块**
   - [config/config.go](backend/config/config.go) - 应用配置管理，包括数据库、JWT等配置

3. **数据模型层**
   - [models/db.go](backend/models/db.go) - 数据库连接和迁移
   - [models/user.go](backend/models/user.go) - 用户和好友关系模型及操作方法
   - [models/group.go](backend/models/group.go) - 群组和群组成员模型及操作方法
   - [models/message.go](backend/models/message.go) - 消息模型及操作方法

4. **中间件**
   - [middlewares/jwt.go](backend/middlewares/jwt.go) - JWT身份验证中间件

5. **控制器层**
   - [controllers/auth.go](backend/controllers/auth.go) - 用户注册和登录接口
   - [controllers/user.go](backend/controllers/user.go) - 用户资料管理和密码修改接口
   - [controllers/friend.go](backend/controllers/friend.go) - 好友关系管理接口
   - [controllers/group.go](backend/controllers/group.go) - 群组管理接口
   - [controllers/message.go](backend/controllers/message.go) - 消息发送和获取接口

6. **WebSocket实时通信**
   - [websocket/hub.go](backend/websocket/hub.go) - WebSocket连接管理和消息广播
   - [websocket/connection.go](backend/websocket/connection.go) - WebSocket连接处理

## 从头到尾编写Go项目的顺序

如果您想从零开始构建这个项目，建议按照以下步骤进行开发：

1. **项目初始化**
   - 创建项目目录结构
   - 初始化Go模块 (`go mod init`)
   - 安装必要的依赖包

2. **配置模块开发**
   - 创建[config/config.go](backend/config/config.go)
   - 定义应用配置结构体
   - 实现配置初始化函数

3. **数据模型层开发**
   - 创建[models/db.go](backend/models/db.go)
   - 实现数据库连接和自动迁移功能
   - 创建[model/user.go](backend/models/user.go)
   - 实现用户和好友关系模型及操作方法
   - 创建[model/group.go](backend/models/group.go)
   - 实现群组和群组成员模型及操作方法
   - 创建[model/message.go](backend/models/message.go)
   - 实现消息模型及操作方法

4. **中间件开发**
   - 创建[middlewares/jwt.go](backend/middlewares/jwt.go)
   - 实现JWT身份验证中间件

5. **控制器层开发**
   - 创建[controllers/auth.go](backend/controllers/auth.go)
   - 实现用户注册和登录接口
   - 创建[controllers/user.go](backend/controllers/user.go)
   - 实现用户资料管理和密码修改接口
   - 创建[controllers/friend.go](backend/controllers/friend.go)
   - 实现好友关系管理接口
   - 创建[controllers/group.go](backend/controllers/group.go)
   - 实现群组管理接口
   - 创建[controllers/message.go](backend/controllers/message.go)
   - 实现消息发送和获取接口

6. **WebSocket实时通信开发**
   - 创建[websocket/connection.go](backend/websocket/connection.go)
   - 实现WebSocket连接处理逻辑
   - 创建[websocket/hub.go](backend/websocket/hub.go)
   - 实现WebSocket连接管理和消息广播功能

7. **主程序开发**
   - 创建[main.go](backend/main.go)
   - 配置路由和中间件
   - 启动HTTP服务器

## 部署指南

待补充
