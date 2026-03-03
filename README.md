# Ingress Portal

Kubernetes Ingress 自动发现与导航门户 - 统一的 Web 入口导航服务

## 项目简介

Ingress Portal 是一个 Kubernetes 集群 Ingress 资源自动发现与可视化管理工具。它能够自动发现集群中的 Ingress 资源，提供统一的 Web 入口导航服务，并通过 Super Mode 支持管理可见性。

### 主要特性

- 🚀 **自动发现**: 自动发现 K8s 集群中的所有 Ingress 资源
- 🎯 **统一导航**: 提供统一的 Web 入口导航服务
- 🔐 **Super Mode**: 支持通过 Token 或密码认证管理可见性
- 📊 **健康检查**: 实时检查后端服务的健康状态
- 🏷️ **元数据支持**: 通过 Annotations 支持分组、描述、团队、优先级等元数据
- 🎨 **现代 UI**: 基于 Vue 3 的现代化前端界面

## 技术栈

### 后端
- Go 1.21+
- Gin - Web 框架
- Kubernetes client-go - K8s API 客户端
- Cobra - CLI 框架
- JWT - Token 认证

### 前端
- Vue 3 + TypeScript
- Vite - 构建工具
- Pinia - 状态管理
- Vue Router - 路由管理
- Axios - HTTP 客户端

## 快速开始

### 前置要求

- Go 1.21+
- Node.js 18+
- Kubernetes 集群 (或 kubeconfig 配置)

### 本地开发

1. **克隆项目**
```bash
git clone https://github.com/example/ingress-portal.git
cd ingress-portal
```

2. **安装依赖**
```bash
# 后端依赖
go mod download

# 前端依赖
cd frontend
npm install
cd ..
```

3. **构建前端**
```bash
make build-frontend
```

4. **运行服务**
```bash
# 直接运行
make run

# 或使用 go run
go run cmd/server/main.go

# 或构建后运行
make build
./bin/ingress-portal
```

服务将在 `http://localhost:8080` 启动

### 使用 Docker

```bash
# 构建镜像
docker build -t ingress-portal:latest .

# 运行容器
docker run -p 8080:8080 \
  -v ~/.kube:/root/.kube \
  ingress-portal:latest
```

### 部署到 Kubernetes

```bash
# 应用部署清单
kubectl apply -f deploy/

# 或使用 Makefile
make deploy
```

## 使用说明

### 命令行参数

```bash
ingress-portal [flags]

Flags:
  -p, --port string           服务端口 (default "8080")
      --super-mode string     Super Mode 认证类型: token 或 password (default "token")
      --super-password string Super Mode 密码 (当 type=password 时)
      --token-ttl int        Token 有效期 (小时) (default 24)
```

### Token 管理

```bash
# 生成新的 Token
ingress-portal token generate --ttl 24

# 查看 Token 状态
ingress-portal token status

# 吊销指定 Token
ingress-portal token revoke <token>

# 吊销所有 Token
ingress-portal token revoke --all
```

## Ingress Annotations

通过 Annotations 配置 Ingress 在门户中的显示方式：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-app
  annotations:
    # 可见性控制 (必需)
    portal.example.com/visible: "true"

    # 可选元数据
    portal.example.com/group: "开发环境"
    portal.example.com/description: "我的应用服务"
    portal.example.com/team: "后端团队"
    portal.example.com/priority: "100"
spec:
  # ... Ingress 配置
```

### Annotations 说明

| Annotation | 说明 | 必需 | 默认值 |
|-----------|------|------|--------|
| `portal.example.com/visible` | 是否在门户中显示 | 是 | true |
| `portal.example.com/group` | 分组名称 | 否 | - |
| `portal.example.com/description` | 服务描述 | 否 | - |
| `portal.example.com/team` | 所属团队 | 否 | - |
| `portal.example.com/priority` | 显示优先级 | 否 | 0 |

## API 接口

### 获取 Ingress 列表
```
GET /api/ingresses
```

响应示例：
```json
{
  "ingresses": [
    {
      "name": "my-app",
      "namespace": "default",
      "host": "my-app.example.com",
      "path": "/",
      "service": "my-app-svc:80",
      "visible": true,
      "group": "开发环境",
      "description": "我的应用服务",
      "team": "后端团队",
      "priority": 100,
      "backendStatus": "healthy",
      "discoveredAt": "2024-01-01T00:00:00Z"
    }
  ],
  "count": 1
}
```

### 刷新 Ingress 列表
```
GET /api/ingresses/refresh
```

### Super Mode 认证
```
POST /api/auth/super-mode
```

## 项目结构

```
ingress-portal/
├── cmd/
│   └── server/          # 主程序入口
│       └── main.go
├── internal/
│   ├── api/             # API 路由和处理
│   ├── auth/            # 认证逻辑
│   └── k8s/             # Kubernetes 客户端
├── pkg/
│   └── models/          # 数据模型
├── frontend/            # 前端代码
│   ├── src/
│   │   ├── components/  # Vue 组件
│   │   ├── views/       # 页面视图
│   │   ├── stores/      # Pinia 状态管理
│   │   ├── router/      # 路由配置
│   │   └── types/       # TypeScript 类型
│   └── package.json
├── web/
│   └── dist/            # 前端构建产物
├── Makefile             # 构建脚本
├── Dockerfile           # Docker 构建文件
└── go.mod               # Go 依赖管理
```

## 开发指南

### 构建命令

```bash
# 构建后端
make build

# 构建前端
make build-frontend

# 构建所有
make build-all

# 运行开发服务器
make run

# 清理构建产物
make clean
```

### 环境变量

- `KUBECONFIG`: 指定 kubeconfig 文件路径

## 后端健康状态

后端健康状态通过 Kubernetes Endpoints 检查：

- `healthy`: 所有后端 Pod 就绪
- `degraded`: 部分后端 Pod 就绪
- `unhealthy`: 没有后端 Pod 就绪
- `unknown`: 无法获取后端状态

## 许可证

[MIT License](LICENSE)

## 贡献指南

欢迎提交 Issue 和 Pull Request！

1. Fork 本仓库
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 联系方式

项目维护者: [Your Name](mailto:your.email@example.com)
