# Development tasks
setup:
	@echo "Setting up project..."
	cd apps/api && go mod tidy
	cd apps/web && npm install

# Start development servers
dev:
	docker-compose -f docker-compose.dev.yml up --build

# Run tests
test:
	cd apps/api && go test ./...
	cd apps/web && npm run test

# Build for production
build:
	cd apps/api && go build -o bin/server cmd/server/main.go
	cd apps/web && npm run build
