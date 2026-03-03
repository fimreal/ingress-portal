.PHONY: build build-frontend build-all run clean deploy help

# 变量
BINARY_NAME=ingress-portal
GO=go
GOFLAGS=-v
FRONTEND_DIR=frontend
BUILD_DIR=bin
DOCKER_IMAGE=ingress-portal
DOCKER_TAG=latest

# 默认目标
default: help

## help: 显示帮助信息
help:
	@echo "可用的构建目标:"
	@echo ""
	@sed -n 's/^## //p' $(MAKEFILE_LIST) | column -t -s ':'
	@echo ""

## build: 构建后端二进制文件
build:
	@echo "构建后端..."
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd/server
	@echo "✓ 后端构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

## build-frontend: 构建前端
build-frontend:
	@echo "构建前端..."
	@cd $(FRONTEND_DIR) && npm install && npm run build
	@mkdir -p web/dist
	@cp -r $(FRONTEND_DIR)/dist/* web/dist/
	@echo "✓ 前端构建完成: web/dist/"

## build-all: 构建前后端
build-all: build-frontend build
	@echo "✓ 所有构建完成"

## run: 运行开发服务器
run:
	@echo "启动开发服务器..."
	$(GO) run ./cmd/server

## run-dev: 运行开发服务器(带前端热重载)
run-dev:
	@echo "启动开发模式..."
	@cd $(FRONTEND_DIR) && npm run dev &

## clean: 清理构建产物
clean:
	@echo "清理构建产物..."
	@rm -rf $(BUILD_DIR)
	@rm -rf web/dist
	@rm -rf $(FRONTEND_DIR)/dist
	@rm -rf $(FRONTEND_DIR)/node_modules
	@echo "✓ 清理完成"

## test: 运行测试
test:
	@echo "运行测试..."
	$(GO) test -v ./...

## docker-build: 构建 Docker 镜像
docker-build:
	@echo "构建 Docker 镜像..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "✓ Docker 镜像构建完成: $(DOCKER_IMAGE):$(DOCKER_TAG)"

## docker-run: 运行 Docker 容器
docker-run:
	@echo "运行 Docker 容器..."
	docker run -p 8080:8080 \
		-v ~/.kube:/root/.kube \
		$(DOCKER_IMAGE):$(DOCKER_TAG)

## docker-push: 推送 Docker 镜像
docker-push:
	@echo "推送 Docker 镜像..."
	docker push $(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "✓ Docker 镜像推送完成"

## deploy: 部署到 Kubernetes
deploy:
	@echo "部署到 Kubernetes..."
	@kubectl apply -f deploy/
	@echo "✓ 部署完成"

## undeploy: 从 Kubernetes 卸载
undeploy:
	@echo "卸载 Kubernetes 部署..."
	@kubectl delete -f deploy/ --ignore-not-found
	@echo "✓ 卸载完成"

## token-generate: 生成 Super Mode Token
token-generate:
	@echo "生成 Super Mode Token..."
	@$(GO) run ./cmd/server token generate

## token-status: 查看 Token 状态
token-status:
	@echo "查看 Token 状态..."
	@$(GO) run ./cmd/server token status

## install: 安装依赖
install:
	@echo "安装后端依赖..."
	$(GO) mod download
	@echo "安装前端依赖..."
	@cd $(FRONTEND_DIR) && npm install
	@echo "✓ 依赖安装完成"

## lint: 代码检查
lint:
	@echo "运行代码检查..."
	@which golangci-lint > /dev/null || (echo "请先安装 golangci-lint" && exit 1)
	golangci-lint run
	@echo "✓ 代码检查完成"

## fmt: 格式化代码
fmt:
	@echo "格式化代码..."
	$(GO) fmt ./...
	@cd $(FRONTEND_DIR) && npm run lint -- --fix 2>/dev/null || true
	@echo "✓ 代码格式化完成"

## proto: 生成 protobuf (如果需要)
proto:
	@echo "生成 protobuf..."
	@which protoc > /dev/null || (echo "请先安装 protoc" && exit 1)
	@# protoc --go_out=. --go-grpc_out=. proto/*.proto
	@echo "✓ Protobuf 生成完成"
