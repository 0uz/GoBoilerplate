# Go Boilerplate Makefile

# Variables
APP_NAME=goauthboiler
DOCKER_COMPOSE=docker-compose
MONITORING_DIR=monitoring

# Colors for output
GREEN=\033[0;32m
YELLOW=\033[1;33m
RED=\033[0;31m
NC=\033[0m # No Color

.PHONY: help build run stop clean test monitoring-start monitoring-stop logs dev

# Default target
help: ## Show this help message
	@echo "$(GREEN)Go Boilerplate - Available Commands$(NC)"
	@echo ""
	@echo "  $(YELLOW)Application:$(NC)"
	@echo "    run                Start the application"
	@echo "    stop               Stop the application"
	@echo "    restart            Restart the application"
	@echo "    rebuild            Rebuild and restart"
	@echo ""
	@echo "  $(YELLOW)Development:$(NC)"
	@echo "    dev                Start dependencies only"
	@echo "    test               Run tests"
	@echo "    test-coverage      Run tests with coverage"
	@echo "    logs               Show application logs"
	@echo ""
	@echo "  $(YELLOW)Monitoring:$(NC)"
	@echo "    monitoring-start   Start monitoring stack"
	@echo "    monitoring-stop    Stop monitoring stack"
	@echo "    start-all          Start app + monitoring"
	@echo "    stop-all           Stop everything"
	@echo ""
	@echo "  $(YELLOW)Utilities:$(NC)"
	@echo "    health             Check service health"
	@echo "    status             Show service status"
	@echo "    clean              Clean up containers"
	@echo "    clean-all          Clean up everything"

# Application commands
build: ## Build the application
	@echo "$(GREEN)Building application...$(NC)"
	go build -o bin/app cmd/main.go

run: ## Start the application with Docker
	@echo "$(GREEN)Starting application...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Application started!$(NC)"
	@echo "  📱 App: http://localhost:8080/api/v1/ready"
	@echo "  📊 Metrics: http://localhost:8080/api/v1/metrics"

stop: ## Stop the application
	@echo "$(YELLOW)Stopping application...$(NC)"
	$(DOCKER_COMPOSE) down

restart: ## Restart the application
	@echo "$(YELLOW)Restarting application...$(NC)"
	$(DOCKER_COMPOSE) restart

rebuild: ## Rebuild and restart the application
	@echo "$(YELLOW)Rebuilding application...$(NC)"
	$(DOCKER_COMPOSE) up -d --build

# Development commands
dev: ## Start application for development (dependencies only)
	@echo "$(GREEN)Starting development environment...$(NC)"
	$(DOCKER_COMPOSE) up -d postgres valkey
	@echo "$(GREEN)Dependencies started. Run 'go run cmd/main.go' to start the app locally$(NC)"

test: ## Run tests
	@echo "$(GREEN)Running tests...$(NC)"
	go test ./...

test-coverage: ## Run tests with coverage
	@echo "$(GREEN)Running tests with coverage...$(NC)"
	go test -cover ./...

# Monitoring commands
monitoring-start: ## Start monitoring stack
	@echo "$(GREEN)Starting monitoring stack...$(NC)"
	@if ! docker network ls --filter name=goboilerplate_app-network --quiet | grep -q .; then \
		echo "$(RED)❌ Main application network not found.$(NC)"; \
		echo "$(YELLOW)   Please start the main application first with 'make run'$(NC)"; \
		exit 1; \
	fi
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)✨ Monitoring stack started!$(NC)"
	@echo ""
	@echo "$(GREEN)📊 Monitoring URLs:$(NC)"
	@echo "   🎨 Grafana: http://localhost:3000 (admin/admin123)"
	@echo "   📈 Prometheus: http://localhost:9090"
	@echo "   📋 Logs: http://localhost:3100"
	@echo ""
	@echo "$(GREEN)📱 Application URLs:$(NC)"
	@echo "   📊 App Metrics: http://localhost:8080/api/v1/metrics"
	@echo "   💚 Health Check: http://localhost:8080/api/v1/ready"

monitoring-stop: ## Stop monitoring stack
	@echo "$(YELLOW)Stopping monitoring stack...$(NC)"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) down

monitoring-restart: ## Restart monitoring stack
	@echo "$(YELLOW)Restarting monitoring stack...$(NC)"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) restart

# Logs
logs: ## Show application logs
	@echo "$(GREEN)Application logs:$(NC)"
	$(DOCKER_COMPOSE) logs -f $(APP_NAME)

logs-monitoring: ## Show monitoring logs
	@echo "$(GREEN)Monitoring logs:$(NC)"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) logs -f

# Cleanup commands
clean: ## Clean up containers and volumes
	@echo "$(YELLOW)Cleaning up...$(NC)"
	$(DOCKER_COMPOSE) down -v
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) down -v
	docker system prune -f

clean-all: ## Clean up everything including images
	@echo "$(RED)Cleaning up everything...$(NC)"
	$(DOCKER_COMPOSE) down -v --rmi all
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) down -v --rmi all
	docker system prune -af

# Combined commands
start-all: run monitoring-start ## Start application and monitoring
	@echo "$(GREEN)🚀 Everything is up and running!$(NC)"

stop-all: stop monitoring-stop ## Stop application and monitoring
	@echo "$(YELLOW)🛑 Everything stopped$(NC)"

status: ## Show status of all services
	@echo "$(GREEN)Application status:$(NC)"
	@$(DOCKER_COMPOSE) ps
	@echo ""
	@echo "$(GREEN)Monitoring status:$(NC)"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) ps

# Health checks
health: ## Check application health
	@echo "$(GREEN)Checking application health...$(NC)"
	@curl -s http://localhost:8080/api/v1/ready > /dev/null && echo "$(GREEN)✅ Application is healthy$(NC)" || echo "$(RED)❌ Application is not responding$(NC)"
	@curl -s http://localhost:3000/api/health > /dev/null && echo "$(GREEN)✅ Grafana is healthy$(NC)" || echo "$(YELLOW)⚠️  Grafana is not responding$(NC)"
	@curl -s http://localhost:9090/-/healthy > /dev/null && echo "$(GREEN)✅ Prometheus is healthy$(NC)" || echo "$(YELLOW)⚠️  Prometheus is not responding$(NC)"
