help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

mockgen: ## Generate mock
	go install go.uber.org/mock/mockgen@latest
	mockgen -destination=internal/mock/worker_pool.go -source=internal/worker_pool.go
	mockgen -destination=internal/pkg/worker/mock/worker.go -source=internal/pkg/worker/interface.go

httpPort := 8080
run: ## Start the server with default config
	go run ./cmd/...

run-with-config: ## Start the server with a custom config
	export $(shell cat custom.env | xargs) && go run ./cmd/...

tests: ## Run the tests
	go test -v -cover -race -count=1 ./...

submit-greeting-jobs: ## Submit greeting jobs with optional names
	@read -p "Enter names (comma separated, default: Joe,Lily,Bob,Kevin,Henry,Denies,Tom): " names; \
	names=$${names:-Joe,Lily,Bob,Kevin,Henry,Denies,Tom}; \
	url="http://localhost:$(httpPort)/greeting?$$(echo $$names | sed 's/[^,]*/name=&/g' | tr ',' '&')"; \
	curl --location --request POST "$$url"
