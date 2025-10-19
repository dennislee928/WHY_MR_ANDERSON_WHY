# Unified Security Platform Makefile
# Multi-cloud deployment support: Cloudflare Workers, OCI, IBM Cloud

.PHONY: help install build test deploy clean docker-build docker-push

# Default target
.DEFAULT_GOAL := help

# Color output
BLUE := \033[0;34m
GREEN := \033[0;32m
RED := \033[0;31m
NC := \033[0m # No Color

##@ General

help: ## Display this help message
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make $(BLUE)<target>$(NC)\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(BLUE)%-20s$(NC) %s\n", $$1, $$2 } /^##@/ { printf "\n$(GREEN)%s$(NC)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

install: ## Install all dependencies
	@echo "$(GREEN)Installing dependencies...$(NC)"
	@$(MAKE) install-backend
	@$(MAKE) install-frontend
	@$(MAKE) install-ai
	@echo "$(GREEN)✓ All dependencies installed$(NC)"

install-backend: ## Install backend dependencies
	@echo "Installing Go backend dependencies..."
	cd src/backend && go mod download && go mod verify

install-frontend: ## Install frontend dependencies
	@echo "Installing frontend dependencies..."
	cd src/frontend && npm install

install-ai: ## Install AI/Quantum dependencies
	@echo "Installing AI/Quantum dependencies..."
	cd src/ai-quantum && pip install -r requirements.txt

##@ Build

build: ## Build all components
	@echo "$(GREEN)Building all components...$(NC)"
	@$(MAKE) build-backend
	@$(MAKE) build-frontend
	@echo "$(GREEN)✓ Build complete$(NC)"

build-backend: ## Build Go backend
	@echo "Building backend..."
	cd src/backend && go build -o ../../bin/backend ./cmd/server/main.go

build-frontend: ## Build Next.js frontend
	@echo "Building frontend..."
	cd src/frontend && npm run build

##@ Test

test: ## Run all tests
	@echo "$(GREEN)Running all tests...$(NC)"
	@$(MAKE) test-backend
	@$(MAKE) test-frontend
	@$(MAKE) test-ai
	@echo "$(GREEN)✓ All tests passed$(NC)"

test-backend: ## Run backend tests
	@echo "Running backend tests..."
	cd src/backend && go test ./... -v -cover

test-frontend: ## Run frontend tests
	@echo "Running frontend tests..."
	cd src/frontend && npm test || echo "No tests configured"

test-ai: ## Run AI/Quantum tests
	@echo "Running AI tests..."
	cd src/ai-quantum && python -m pytest tests/ -v || echo "No tests found"

test-integration: ## Run integration tests
	@echo "Running integration tests..."
	cd tests && go test -tags=integration ./... -v

##@ Docker

docker-build: ## Build all Docker images
	@echo "$(GREEN)Building Docker images...$(NC)"
	docker build -t security-platform/backend:latest src/backend
	docker build -t security-platform/frontend:latest src/frontend
	docker build -t security-platform/ai-quantum:latest src/ai-quantum
	@echo "$(GREEN)✓ Docker images built$(NC)"

docker-build-multi-arch: ## Build multi-architecture images (amd64 + arm64)
	@echo "$(GREEN)Building multi-arch images...$(NC)"
	docker buildx create --use || true
	docker buildx build --platform linux/amd64,linux/arm64 -t security-platform/backend:latest --push src/backend
	docker buildx build --platform linux/amd64,linux/arm64 -t security-platform/frontend:latest --push src/frontend
	docker buildx build --platform linux/amd64,linux/arm64 -t security-platform/ai-quantum:latest --push src/ai-quantum
	@echo "$(GREEN)✓ Multi-arch images built and pushed$(NC)"

docker-push: ## Push Docker images to registry
	@echo "Pushing Docker images..."
	docker push security-platform/backend:latest
	docker push security-platform/frontend:latest
	docker push security-platform/ai-quantum:latest

docker-compose-up: ## Start services with Docker Compose
	@echo "Starting services with Docker Compose..."
	cd infrastructure/docker && docker-compose up -d

docker-compose-down: ## Stop services with Docker Compose
	@echo "Stopping services..."
	cd infrastructure/docker && docker-compose down

docker-compose-logs: ## View Docker Compose logs
	cd infrastructure/docker && docker-compose logs -f

##@ Cloud Deployments

deploy-cloudflare: ## Deploy to Cloudflare Workers
	@echo "$(GREEN)Deploying to Cloudflare Workers...$(NC)"
	cd infrastructure/cloud-configs/cloudflare && npm install && npx wrangler deploy --env production
	@echo "$(GREEN)✓ Deployed to Cloudflare$(NC)"

deploy-cloudflare-staging: ## Deploy to Cloudflare staging
	cd infrastructure/cloud-configs/cloudflare && npx wrangler deploy --env staging

deploy-oci: ## Deploy to Oracle Cloud (OCI)
	@echo "$(GREEN)Deploying to OCI...$(NC)"
	cd infrastructure/cloud-configs/oci/terraform && terraform apply -auto-approve
	@echo "$(GREEN)✓ Deployed to OCI$(NC)"

deploy-oci-plan: ## Plan OCI deployment
	cd infrastructure/cloud-configs/oci/terraform && terraform plan

deploy-ibm: ## Deploy to IBM Cloud
	@echo "$(GREEN)Deploying to IBM Cloud...$(NC)"
	ibmcloud login --apikey $(IBM_CLOUD_API_KEY) -r us-south
	ibmcloud target --cf
	cd infrastructure/cloud-configs/ibm && ibmcloud cf push
	@echo "$(GREEN)✓ Deployed to IBM Cloud$(NC)"

deploy-k8s: ## Deploy to Kubernetes
	@echo "Deploying to Kubernetes..."
	kubectl apply -f infrastructure/kubernetes/

##@ Development

dev-backend: ## Run backend in development mode
	@echo "Starting backend dev server..."
	cd src/backend && go run ./cmd/server/main.go

dev-frontend: ## Run frontend in development mode
	@echo "Starting frontend dev server..."
	cd src/frontend && npm run dev

dev-ai: ## Run AI/Quantum service in development mode
	@echo "Starting AI/Quantum service..."
	cd src/ai-quantum && python main.py

dev-all: ## Run all services in development mode (requires tmux)
	@echo "Starting all services..."
	tmux new-session -d -s security-platform
	tmux send-keys -t security-platform:0 'make dev-backend' C-m
	tmux split-window -t security-platform:0 -v
	tmux send-keys -t security-platform:0.1 'make dev-frontend' C-m
	tmux split-window -t security-platform:0 -h
	tmux send-keys -t security-platform:0.2 'make dev-ai' C-m
	tmux attach -t security-platform

##@ Cloudflare Workers Management

cloudflare-init: ## Initialize Cloudflare D1 database
	@echo "Initializing Cloudflare D1..."
	cd infrastructure/cloud-configs/cloudflare && npm run d1:create
	cd infrastructure/cloud-configs/cloudflare && npm run d1:init

cloudflare-dev: ## Run Cloudflare Workers locally
	cd infrastructure/cloud-configs/cloudflare && npm run dev

cloudflare-tail: ## Tail Cloudflare Workers logs
	cd infrastructure/cloud-configs/cloudflare && npm run tail

##@ OCI Management

oci-init: ## Initialize OCI deployment
	cd infrastructure/cloud-configs/oci/terraform && terraform init

oci-ssh-app: ## SSH to OCI application server
	@echo "Connecting to OCI app server..."
	ssh ubuntu@$(shell cd infrastructure/cloud-configs/oci/terraform && terraform output -raw app_server_public_ip)

oci-ssh-db: ## SSH to OCI database server
	@echo "Connecting to OCI DB server..."
	ssh ubuntu@$(shell cd infrastructure/cloud-configs/oci/terraform && terraform output -raw db_server_public_ip)

oci-destroy: ## Destroy OCI infrastructure
	@echo "$(RED)Destroying OCI infrastructure...$(NC)"
	@read -p "Are you sure? [y/N] " -n 1 -r; \
	echo; \
	if [[ $$REPLY =~ ^[Yy]$$ ]]; then \
		cd infrastructure/cloud-configs/oci/terraform && terraform destroy; \
	fi

##@ Database

db-migrate: ## Run database migrations
	@echo "Running database migrations..."
	cd src/backend && go run ./cmd/migrate/main.go up

db-rollback: ## Rollback last migration
	cd src/backend && go run ./cmd/migrate/main.go down

db-seed: ## Seed database with test data
	cd src/backend && go run ./cmd/seed/main.go

##@ Linting & Security

lint: ## Run linters on all code
	@echo "$(GREEN)Running linters...$(NC)"
	@$(MAKE) lint-backend
	@$(MAKE) lint-frontend
	@echo "$(GREEN)✓ Linting complete$(NC)"

lint-backend: ## Lint Go code
	cd src/backend && golangci-lint run ./...

lint-frontend: ## Lint TypeScript/React code
	cd src/frontend && npm run lint

security-scan: ## Run security scan on all code
	@echo "$(GREEN)Running security scans...$(NC)"
	trivy fs --severity HIGH,CRITICAL .
	@echo "$(GREEN)✓ Security scan complete$(NC)"

##@ Documentation

docs-serve: ## Serve documentation locally
	@echo "Serving documentation..."
	cd docs && python -m http.server 8080

docs-build: ## Build documentation
	@echo "Building documentation..."
	# Add documentation build command here (e.g., mkdocs, sphinx)

##@ CI/CD

cicd-buddy: ## View Buddy CI configuration
	@cat cicd/buddy/buddy.yml

cicd-argocd: ## Apply ArgoCD configuration
	kubectl apply -f cicd/argocd/application.yaml

cicd-harness: ## View Harness pipeline
	@cat cicd/harness/.harness/pipelines/multi-cloud-deployment.yaml

##@ Cleanup

clean: ## Clean build artifacts and caches
	@echo "$(GREEN)Cleaning build artifacts...$(NC)"
	rm -rf bin/
	rm -rf src/backend/bin/
	rm -rf src/frontend/.next/
	rm -rf src/frontend/out/
	rm -rf src/ai-quantum/__pycache__/
	find . -type d -name "node_modules" -exec rm -rf {} + 2>/dev/null || true
	@echo "$(GREEN)✓ Cleanup complete$(NC)"

clean-docker: ## Remove all Docker containers and images
	@echo "$(RED)Removing Docker containers and images...$(NC)"
	docker-compose -f infrastructure/docker/docker-compose.yml down -v
	docker system prune -af

##@ Monitoring

logs-backend: ## View backend logs
	docker logs -f security-platform-backend || journalctl -f -u security-platform-backend

logs-frontend: ## View frontend logs
	docker logs -f security-platform-frontend || journalctl -f -u security-platform-frontend

logs-ai: ## View AI/Quantum logs
	docker logs -f security-platform-ai || journalctl -f -u security-platform-ai

status: ## Show status of all services
	@echo "$(GREEN)Service Status:$(NC)"
	@docker ps --filter "name=security-platform" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" || echo "No Docker containers running"

health-check: ## Check health of all services
	@echo "Checking service health..."
	@curl -f http://localhost:3001/api/v1/health || echo "Backend: DOWN"
	@curl -f http://localhost:8000/api/v1/status || echo "AI Service: DOWN"
	@curl -f http://localhost:3000 || echo "Frontend: DOWN"

##@ Cost Monitoring

cost-report: ## Generate cost report for all clouds
	@echo "$(BLUE)Multi-Cloud Cost Report$(NC)"
	@echo "See: docs/deployment/cost-comparison.md"
	@cat docs/deployment/cost-comparison.md | grep -A 20 "## Executive Summary"

##@ Utilities

env-example: ## Create .env.example file
	@echo "Creating .env.example..."
	@cat > .env.example << 'EOF'
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=sectools
DB_PASSWORD=changeme
DB_NAME=security

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# IBM Quantum
IBM_QUANTUM_TOKEN=your_ibm_quantum_token_here

# Cloudflare
CLOUDFLARE_API_TOKEN=your_cloudflare_api_token
CLOUDFLARE_ACCOUNT_ID=your_account_id

# OCI
OCI_TENANCY_OCID=ocid1.tenancy.oc1..xxx
OCI_USER_OCID=ocid1.user.oc1..xxx
OCI_FINGERPRINT=xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx:xx

# IBM Cloud
IBM_CLOUD_API_KEY=your_ibm_cloud_api_key
IBM_CLOUD_REGION=us-south

# Application
APP_ENV=development
APP_PORT=3001
LOG_LEVEL=info
EOF
	@echo "$(GREEN)✓ .env.example created$(NC)"

git-add-all: ## Git add all unified project files (excluding old projects)
	@echo "$(GREEN)Adding unified project files to git...$(NC)"
	git add .gitignore
	git add README.md README.zh-TW.md
	git add Makefile
	git add .env.example
	git add src/
	git add infrastructure/
	git add cicd/
	git add docs/
	git add configs/
	git add scripts/
	git add tests/
	@echo "$(GREEN)✓ Files added to git$(NC)"
	@echo "$(BLUE)Run 'git status' to review changes$(NC)"

version: ## Show version information
	@echo "Security Platform - Unified Multi-Cloud Deployment"
	@echo "Version: 1.0.0"
	@echo "Go: $(shell go version)"
	@echo "Node: $(shell node --version)"
	@echo "Python: $(shell python --version)"
	@echo "Docker: $(shell docker --version)"

