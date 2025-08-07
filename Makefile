# --- Load .env file if it exists ---
	include deploy/.env
	export

# --- Phony Targets ---
# Always declare targets that don't produce a file of the same name as .PHONY
# This ensures make runs the recipe even if a file with that name exists.
.PHONY: help run setup build run_build test lint nilaway migrate_new migrate_up packages_install packages_update dev swagger_docgen docgen cloc compose_up compose_down compose_restart compose_build compose_logs compose_exec compose_down_volume local_migrate_up

# --- Help Target ---
help: ## Show this help message
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	# Use \ to continue a line in make, and ensure each @echo has proper quotes.
	# The \t ensures a tab for alignment, but you still need actual tabs for recipe indentation.
	@echo "	run                	Run server"
	@echo "	setup              	Install pre-requisites"
	@echo "	build              	Build Binaries for linux, windows and mac"
	@echo "	run_build          	Run server from build"
	@echo "	compose_up         	Start Docker Compose services"
	@echo "	compose_down       	Stop Docker Compose services"
	@echo "	compose_down_volume    	Remove Docker Compose services and volumes"
	@echo "	compose_restart    	Restart Docker Compose services"
	@echo "	compose_build      	Rebuild Docker Compose services"
	@echo "	compose_logs       	Show Docker Compose logs"
	@echo "	compose_exec       	Execute command in Docker Compose service"
	@echo "	test               	Run tests"
	@echo "	lint               	Run linter"
	@echo "	nilaway            	Run nilaway"
	@echo "	migrate_new        	Create new migration"
	@echo "	migrate_up         	Apply pending migrations"
	@echo "	local_migrate_up   	Apply pending migrations (local)"
	@echo "	packages_install   	Install packages"
	@echo "	packages_update    	Update packages"
	@echo "	dev                	Run dev server"
	@echo "	swagger_docgen     	Generate Swagger Docs"
	@echo "	docgen             	Generate OpenAPIv3 Docs and Swagger Docs"
	@echo "	cloc               	Count lines of code"

# --- Other Targets (ensure these are also indented with TABS, not spaces) ---

run: ## Run server
	@echo "Running server..."
	@go run main.go

setup: ## Install pre-requisites
	@go install go.uber.org/nilaway/cmd/nilaway@latest
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/golang/mock/mockgen@latest
	@go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@go install github.com/oxisto/air@latest

build: ## Build Binaries for linux, windows and mac
	@echo "Building server..."
	@GOOS=linux GOARCH=amd64 go build -o bin/server main.go
	@GOOS=windows GOARCH=amd64 go build -o bin/server.exe main.go
	@GOOS=darwin GOARCH=amd64 go build -o bin/server.darwin main.go
	@echo "Done building for linux, windows and mac (x64 only)"
	@echo "Binaries are in bin/ directory"

run_build: ## Run server from build
	@echo "Running server..."
	@./bin/server

compose_up: ## Start Docker Compose services
	@echo "Starting Docker Compose services..."
	@docker compose -f deploy/docker-compose.yml up -d

compose_down: ## Stop Docker Compose services
	@echo "Stopping Docker Compose services..."
	@docker compose -f deploy/docker-compose.yml down

compose_down_volume: ## Remove Docker Compose services and volumes
	@echo "Removing Docker Compose services and volumes..."
	@docker compose -f deploy/docker-compose.yml down -v

compose_restart: ## Restart Docker Compose services
	@echo "Restarting Docker Compose services..."
	@docker compose -f deploy/docker-compose.yml restart

compose_build: ## Rebuild Docker Compose services
	@echo "Rebuilding Docker Compose services..."
	@docker compose -f deploy/docker-compose.yml up --build

compose_logs: ## Show Docker Compose logs
	@echo "Showing Docker Compose logs..."
	@docker compose -f deploy/docker-compose.yml logs -f

compose_exec: ## Execute command in Docker Compose service
	@echo "Executing command in Docker Compose service..."
	@docker compose -f deploy/docker-compose.yml exec app sh
	@echo "Use 'exit' to leave the container shell."

test: ## Run tests
	@echo "Running tests..."
	@go test -v ./...

lint: ## Run linter
	@echo "Running linter..."
	@golangci-lint run

nilaway: ## Run nilaway
	@echo "Running nilaway..."
	@nilaway ../server/ # Adjust path if nilaway is run from root

migrate_new: ## Create new migration
	@echo "Creating new migration..."
	@read -p "Enter migration name: " name; \
	migrate create -ext sql -dir internal/db/migrations -seq $$name

migrate_up: ## Apply pending migrations
	@echo "Applying pending migrations..."
	migrate -path internal/db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable" up

local_migrate_up: ## Apply pending migrations
	@echo "Applying pending migrations..."
	migrate -path internal/db/migrations -database "postgres://${DB_USER}:${DB_PASSWORD}@localhost:${DB_PORT}/${DB_NAME}?sslmode=disable" up


packages_install: ## Install packages
	@echo "Installing packages..."
	@go mod tidy

packages_update: ## Update packages
	@echo "Updating packages..."
	@go get -u

dev: ## Run dev server
	@air

swagger_docgen: ## Generate Swagger Docs
	@echo "Generating Swagger Docs..."
	@echo ""
	@rm -rf docs/
	@swag init -dir ./cmd
	@echo ""
	@echo "Done generating docs."

docgen: swagger_docgen ## Generate OpenAPIv3 Docs and Swagger Docs
	@echo ""
	@echo "Generating OpenAPIv3 Docs..."
	@echo ""
	@rm -rf docs/openapi.yaml
	@npx -p swagger2openapi swagger2openapi --yaml --outfile docs/openapi.yaml "http://localhost:${PORT}/swagger/doc.json"
	@echo "Done generating OpenAPIv3 Docs."

cloc: ## Count lines of code
	@echo "Counting lines of code..."
	@cloc .