# 引入 golang.mk 规则
include scripts/make-rules/all.mk

# 可以通过命令行参数覆盖这些默认值
.PHONY: build
build: go.tidy ## 编译源码，依赖 tidy 目标自动添加/移除依赖包 例如: make build MAIN_PATH=cmd/other_server/main.go BINARY_NAME=other_server
	$(MAKE) go.build

.PHONY: protoc
protoc: ## 生成 protobuf 文件
	$(MAKE) gen.protoc


## --------------------------------------
## Cleanup
## --------------------------------------

.PHONY: clean
clean: ## 清理构建产物、临时文件等. 例如 _output 目录.
	@echo "===========> Cleaning all build output"
	@-rm -vrf $(OUTPUT_DIR)

## --------------------------------------
## Lint / Verification
## --------------------------------------

.PHONY: lint tidy format
lint: ## 执行静态代码检查.
	@$(MAKE) go.lint

tidy: ## 自动添加/移除依赖包.
	@$(MAKE) go.tidy

## --------------------------------------
## Hack / Tools
## --------------------------------------

.PHONY: swagger
swagger: ## 聚合 swagger 文档到一个 openapi.yaml 文件中.
	@$(MAKE) swagger.run

.PHONY: serve-swagger
serve-swagger: ## 运行 Swagger 文档服务器.
	@$(MAKE) swagger.serve

#================================================================================
# Help
.PHONY: help
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "General targets:"
	@grep -h '##' $(MAKEFILE_LIST) | grep -v 'grep' | sed 's/:.*## /\t/' | sed 's/## /\t/' | awk -F'\t' '{printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
