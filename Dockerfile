# 阶段 1: 构建前端
FROM node:18-alpine AS frontend-builder

WORKDIR /app/frontend

# 复制前端依赖文件
COPY frontend/package*.json ./

# 安装依赖
RUN npm ci

# 复制前端源代码
COPY frontend/ ./

# 构建前端
RUN npm run build

# 阶段 2: 构建后端
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# 安装构建依赖
RUN apk add --no-cache git

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . ./

# 复制前端构建产物
COPY --from=frontend-builder /app/frontend/dist ./web/dist

# 构建后端
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/ingress-portal ./cmd/server

# 阶段 3: 运行时镜像
FROM alpine:latest

WORKDIR /app

# 安装 ca-certificates (用于 HTTPS 连接)
RUN apk --no-cache add ca-certificates tzdata

# 从构建阶段复制二进制文件
COPY --from=backend-builder /app/bin/ingress-portal .

# 从前端构建阶段复制静态文件
COPY --from=frontend-builder /app/frontend/dist ./web/dist

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/ingresses || exit 1

# 运行服务
ENTRYPOINT ["./ingress-portal"]
CMD ["--port", "8080"]
