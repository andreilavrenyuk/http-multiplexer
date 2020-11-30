.PHONY: build run lint help

.DEFAULT_GOAL := help

SERVICE := http_multiplexer

build: ## Build service
	go build -o build/$(SERVICE) main.go 

run: ## Run server task
	go run main.go

lint: ## Run golangci linter 
	golangci-lint run --timeout 10m

help: ## Display callable targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
