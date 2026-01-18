# Go Boilerplate Makefile

# Variables
APP_NAME=goboilerplate
DOCKER_COMPOSE=docker compose
MONITORING_DIR=monitoring

# Colors for output (ANSI escape codes - works on macOS/Linux)
GREEN=\033[0;32m
YELLOW=\033[0;33m
RED=\033[0;31m
NC=\033[0m

.PHONY: help build run stop clean test monitoring-start monitoring-stop logs dev

# Default target
help: ## Show this help message
	@printf "$(GREEN)Go Boilerplate - Available Commands$(NC)\n"
	@printf "\n"
	@printf "  $(YELLOW)Application:$(NC)\n"
	@printf "    run                Start the application\n"
	@printf "    stop               Stop the application\n"
	@printf "    restart            Restart the application\n"
	@printf "    rebuild            Rebuild and restart\n"
	@printf "\n"
	@printf "  $(YELLOW)Development:$(NC)\n"
	@printf "    dev                Start dependencies only\n"
	@printf "    test               Run tests\n"
	@printf "    test-coverage      Run tests with coverage\n"
	@printf "    logs               Show application logs\n"
	@printf "\n"
	@printf "  $(YELLOW)Monitoring:$(NC)\n"
	@printf "    monitoring-start   Start monitoring stack\n"
	@printf "    monitoring-stop    Stop monitoring stack\n"
	@printf "    start-all          Start app + monitoring\n"
	@printf "    stop-all           Stop everything\n"
	@printf "\n"
	@printf "  $(YELLOW)Utilities:$(NC)\n"
	@printf "    health             Check service health\n"
	@printf "    status             Show service status\n"
	@printf "    clean              Clean up containers\n"
	@printf "    clean-all          Clean up everything\n"

# Application commands
build: ## Build the application
	@printf "$(GREEN)Building application...$(NC)\n"
	go build -o bin/app cmd/main.go

run: ## Start the application with Docker
	@printf "$(GREEN)Starting application...$(NC)\n"
	$(DOCKER_COMPOSE) up -d
	@printf "$(GREEN)Application started!$(NC)\n"
	@printf "  >> Health Check: http://localhost:8080/ready\n"
	@printf "  >> Metrics: http://localhost:8080/metrics\n"

stop: ## Stop the application
	@printf "$(YELLOW)Stopping application...$(NC)\n"
	$(DOCKER_COMPOSE) down

restart: ## Restart the application
	@printf "$(YELLOW)Restarting application...$(NC)\n"
	$(DOCKER_COMPOSE) restart

rebuild: ## Rebuild and restart the application
	@printf "$(YELLOW)Rebuilding application...$(NC)\n"
	$(DOCKER_COMPOSE) up -d --build

# Development commands
dev: ## Start application for development (dependencies only)
	@printf "$(GREEN)Starting development environment...$(NC)\n"
	$(DOCKER_COMPOSE) up -d postgres valkey
	@printf "$(GREEN)Dependencies started. Run 'go run cmd/main.go' to start the app locally$(NC)\n"

test: ## Run tests
	@printf "$(GREEN)Running tests...$(NC)\n"
	go test ./...

test-coverage: ## Run tests with coverage
	@printf "$(GREEN)Running tests with coverage...$(NC)\n"
	go test -cover ./...

# Monitoring commands
monitoring-start: ## Start monitoring stack
	@printf "$(GREEN)Starting monitoring stack...$(NC)\n"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) up -d
	@printf "$(GREEN)Monitoring stack started!$(NC)\n"
	@printf "\n"
	@printf "$(GREEN)Monitoring URLs:$(NC)\n"
	@printf "   >> Grafana: http://localhost:3000 (admin/admin123)\n"
	@printf "   >> Prometheus: http://localhost:9090\n"
	@printf "   >> Logs: http://localhost:3100\n"
	@printf "\n"
	@printf "$(GREEN)Application URLs:$(NC)\n"
	@printf "   >> App Metrics: http://localhost:8080/metrics\n"
	@printf "   >> Health Check: http://localhost:8080/ready\n"

monitoring-stop: ## Stop monitoring stack
	@printf "$(YELLOW)Stopping monitoring stack...$(NC)\n"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) down

monitoring-restart: ## Restart monitoring stack
	@printf "$(YELLOW)Restarting monitoring stack...$(NC)\n"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) restart

# Logs
logs: ## Show application logs
	@printf "$(GREEN)Application logs:$(NC)\n"
	$(DOCKER_COMPOSE) logs -f $(APP_NAME)

logs-monitoring: ## Show monitoring logs
	@printf "$(GREEN)Monitoring logs:$(NC)\n"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) logs -f

# Cleanup commands
clean: ## Clean up containers and volumes
	@printf "$(YELLOW)Cleaning up...$(NC)\n"
	$(DOCKER_COMPOSE) down -v
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) down -v
	docker system prune -f

clean-all: ## Clean up everything including images
	@printf "$(RED)Cleaning up everything...$(NC)\n"
	$(DOCKER_COMPOSE) down -v --rmi all
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) down -v --rmi all
	docker system prune -af

# Combined commands
start-all: run monitoring-start ## Start application and monitoring
	@printf "$(GREEN)Everything is up and running!$(NC)\n"

stop-all: stop monitoring-stop ## Stop application and monitoring
	@printf "$(YELLOW)Everything stopped$(NC)\n"

status: ## Show status of all services
	@printf "$(GREEN)Application status:$(NC)\n"
	@$(DOCKER_COMPOSE) ps
	@printf "\n"
	@printf "$(GREEN)Monitoring status:$(NC)\n"
	@cd $(MONITORING_DIR) && $(DOCKER_COMPOSE) ps

# Health checks
health: ## Check application health
	@printf "$(GREEN)Checking application health...$(NC)\n"
	@curl -sf --max-time 5 http://localhost:8080/ready > /dev/null 2>&1 && printf "$(GREEN)[OK] Application is healthy$(NC)\n" || printf "$(RED)[ERROR] Application is not responding$(NC)\n"
	@curl -sf --max-time 5 http://localhost:3000/api/health > /dev/null 2>&1 && printf "$(GREEN)[OK] Grafana is healthy$(NC)\n" || printf "$(YELLOW)[WARN] Grafana is not responding$(NC)\n"
	@curl -sf --max-time 5 http://localhost:9090/-/healthy > /dev/null 2>&1 && printf "$(GREEN)[OK] Prometheus is healthy$(NC)\n" || printf "$(YELLOW)[WARN] Prometheus is not responding$(NC)\n"
