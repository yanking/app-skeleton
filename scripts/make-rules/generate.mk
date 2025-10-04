## ================================================================================
## Protobuf generation
## ================================================================================
.PHONY: gen.protoc
gen.protoc: buf.lint buf.format ## Generate protobuf files
	@echo "Generating protobuf..."
	buf generate $(APIROOT)
	@echo "Protobuf generated successfully!"

## ================================================================================
## Buf commands
## ================================================================================
.PHONY: buf.lint
buf.lint: ## Lint protobuf files
	@echo "Linting protobuf files..."
	buf lint $(APIROOT)
	@echo "Linting completed!"

.PHONY: buf.format
buf.format: ## Format protobuf files
	@echo "Formatting protobuf files..."
	buf format -w $(APIROOT)
	@echo "Formatting completed!"

.PHONY: buf.dep.update
buf.dep.update: ## Update buf dependencies
	@echo "Updating buf dependencies..."
	buf dep update $(APIROOT)
	@echo "Buf dependencies updated!"

.PHONY: buf.check
buf.check: ## Check breaking changes
	@echo "Checking breaking changes..."
	buf breaking $(APIROOT) --against '.git#branch=main'
	@echo "Breaking changes check completed!"