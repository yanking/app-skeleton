GO := go

GO_BUILD_FLAGS += -ldflags "$(GO_LDFLAGS)"

BINARY ?= $(shell basename $(PWD))

# 如果未通过命令行指定 MAIN_PATH，则使用默认值
ifndef MAIN_PATH
MAIN_PATH := cmd/demo_server/main.go
endif

# 如果未通过命令行指定 BINARY_NAME，则从入口文件目录名获取
ifndef BINARY_NAME
BINARY_NAME := $(notdir $(patsubst %/,%,$(dir $(MAIN_PATH))))
endif

.PHONY: go.build
go.build:
	@echo "===========> Building $(BINARY_NAME)"
	@mkdir -p $(OUTPUT_DIR)
	CGO_ENABLED=0 $(GO) build $(GO_BUILD_FLAGS) -o $(OUTPUT_DIR)/$(BINARY_NAME) $(MAIN_PATH)

.PHONY: go.format
go.format: tools.verify.goimports ## 格式化 Go 源码.
	@echo "===========> Running formaters to format codes"
	@$(FIND) -type f -name '*.go' | $(XARGS) gofmt -s -w
	@$(FIND) -type f -name '*.go' | $(XARGS) goimports -w -local $(ROOT_PACKAGE)
	@$(GO) mod edit -fmt

.PHONY: go.tidy
go.tidy: ## 自动添加/移除依赖包.
	@echo "===========> Running 'go mod tidy'..."
	@$(GO) mod tidy