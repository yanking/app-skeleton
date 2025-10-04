## ==============================================================================
## 工具相关的 Makefile
## ==============================================================================

TOOLS ?= golangci-lint goimports swagger buf redocly

.PHONY: tools.verify
tools.verify: $(addprefix tools.verify., $(TOOLS))

.PHONY: tools.install
tools.install: $(addprefix tools.install., $(TOOLS))

.PHONY: tools.install.%
tools.install.%:
	@echo "===========> Installing $*"
	@$(MAKE) install.$*
	@echo "===========> $* installed successfully!"

.PHONY: tools.verify.%
tools.verify.%:
	@if ! which $* &>/dev/null; then $(MAKE) tools.install.$*; fi

.PHONY: install.golangci-lint
install.golangci-lint: ## Install golangci-lint
	@$(GO) install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@golangci-lint completion bash > $(HOME)/.golangci-lint.bash
	@if ! grep -q .golangci-lint.bash $(HOME)/.bashrc; then echo "source \$$HOME/.golangci-lint.bash" >> $(HOME)/.bashrc; fi

.PHONY: install.goimports
install.goimports: ## Install goimports
	@$(GO) install golang.org/x/tools/cmd/goimports@latest

.PHONY: install.protoc-plugins
install.protoc-plugins: ## Install protoc-plugins
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	@$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2@latest
	@$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: install.swagger
install.swagger: ## Install swagger
	@$(GO) install github.com/go-swagger/go-swagger/cmd/swagger@latest

.PHONY: install.buf
install.buf: ## Install buf
	@$(GO) install github.com/bufbuild/buf/cmd/buf@latest

.PHONY: install.grpcurl
install.grpcurl: ## Install grpcurl
	@$(GO) install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest

.PHONY: install.redocly
install.redocly: ## Install Redocly CLI
	@$(GO) install github.com/redocly/redoc/cmd/redoc@latest
