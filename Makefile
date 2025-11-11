.PHONY: build clean install run test help

# 项目名称
BINARY_NAME=todo
BINARY_PATH=./bin/$(BINARY_NAME)

# Go 参数
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# 构建参数
BUILD_FLAGS=-ldflags="-s -w"
BUILD_DIR=./cmd/todo

all: clean build

## build: 编译项目
build:
	@echo "正在编译..."
	@mkdir -p bin
	CGO_ENABLED=1 $(GOBUILD) $(BUILD_FLAGS) -o $(BINARY_PATH) $(BUILD_DIR)
	@echo "编译完成: $(BINARY_PATH)"

## clean: 清理构建文件
clean:
	@echo "正在清理..."
	@$(GOCLEAN)
	@rm -rf bin/
	@echo "清理完成"

## install: 安装到 $GOPATH/bin
install: build
	@echo "正在安装..."
	@cp $(BINARY_PATH) $(GOPATH)/bin/$(BINARY_NAME)
	@echo "安装完成: $(GOPATH)/bin/$(BINARY_NAME)"

## run: 运行程序
run: build
	@$(BINARY_PATH)

## test: 运行测试
test:
	@echo "正在运行测试..."
	@$(GOTEST) -v ./...

## deps: 下载依赖
deps:
	@echo "正在下载依赖..."
	@$(GOMOD) download
	@$(GOMOD) tidy
	@echo "依赖下载完成"

## help: 显示帮助信息
help:
	@echo "TodoList Makefile 命令:"
	@echo ""
	@grep -E '^##' Makefile | sed 's/## /  /'
	@echo ""

.DEFAULT_GOAL := help
