# GoAuth Boilerplate Monitoring Stack Startup Script
# This script starts the complete monitoring infrastructure

Write-Host "🚀 Starting GoAuth Boilerplate Monitoring Stack..." -ForegroundColor Green

# Check if Docker is running
try {
    docker info | Out-Null
    Write-Host "✅ Docker is running" -ForegroundColor Green
} catch {
    Write-Host "❌ Docker is not running. Please start Docker first." -ForegroundColor Red
    exit 1
}

# Check if main application network exists
$networkExists = docker network ls --filter name=goboilerplate_app-network --quiet
if (-not $networkExists) {
    Write-Host "❌ Main application network 'goboilerplate_app-network' not found." -ForegroundColor Red
    Write-Host "   Please start the main application first with 'docker-compose up'" -ForegroundColor Yellow
    exit 1
}

Write-Host "✅ Application network found" -ForegroundColor Green

# Navigate to monitoring directory
$monitoringDir = Split-Path -Parent $MyInvocation.MyCommand.Path
Set-Location $monitoringDir

Write-Host "📂 Working directory: $monitoringDir" -ForegroundColor Blue

# Pull latest images
Write-Host "⬇️  Pulling latest monitoring images..." -ForegroundColor Blue
docker-compose pull

# Start monitoring stack
Write-Host "🎯 Starting monitoring services..." -ForegroundColor Blue
docker-compose up -d

# Wait for services to be healthy
Write-Host "⏳ Waiting for services to become healthy..." -ForegroundColor Blue
Start-Sleep -Seconds 10

# Check service status
$services = @("prometheus", "grafana", "loki", "promtail", "node-exporter", "cadvisor")

Write-Host "`n📊 Service Status:" -ForegroundColor Yellow
foreach ($service in $services) {
    $status = docker-compose ps --services --filter "status=running" | Where-Object { $_ -eq $service }
    if ($status) {
        Write-Host "✅ $service" -ForegroundColor Green
    } else {
        Write-Host "❌ $service" -ForegroundColor Red
    }
}

Write-Host "`n🌐 Access URLs:" -ForegroundColor Cyan
Write-Host "   📊 Grafana:    http://localhost:3000 (admin/admin123)" -ForegroundColor White
Write-Host "   📈 Prometheus: http://localhost:9090" -ForegroundColor White
Write-Host "   📋 Loki:      http://localhost:3100" -ForegroundColor White
Write-Host "   🖥️  Node Exporter: http://localhost:9100" -ForegroundColor White
Write-Host "   🐳 cAdvisor:   http://localhost:8081" -ForegroundColor White

Write-Host "`n🎯 Application Metrics:" -ForegroundColor Cyan
Write-Host "   📊 App Metrics: http://localhost:8080/api/v1/metrics" -ForegroundColor White
Write-Host "   💚 Health Check: http://localhost:8080/api/v1/ready" -ForegroundColor White

Write-Host "`n✨ Monitoring stack is ready!" -ForegroundColor Green
Write-Host "   Check the logs with: docker-compose logs -f" -ForegroundColor Yellow
Write-Host "   Stop the stack with: docker-compose down" -ForegroundColor Yellow
