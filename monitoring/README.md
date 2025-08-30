# 🚀 GoAuth Boilerplate - Modern Monitoring Stack

Bu proje için kurulmuş olan comprehensive monitoring ve observability altyapısı.

## 📊 Stack Bileşenleri

### Core Monitoring
- **Prometheus** - Metrics toplama ve alert sistemi
- **Grafana** - Görselleştirme ve dashboard
- **Loki** - Log aggregation
- **Promtail** - Log collector

### System Monitoring  
- **Node Exporter** - Sistem metrics'leri
- **cAdvisor** - Container metrics'leri
- **PostgreSQL Exporter** - Database metrics'leri
- **Redis Exporter** - Cache metrics'leri

## 🎯 İzlenen Metrikler

### Application Metrics
- HTTP request rate, duration, error rate
- Authentication attempts ve success rate
- Active sessions
- Cache hit/miss oranları
- Database connection pool

### System Metrics
- CPU, Memory, Disk kullanımı
- Network I/O
- Container resource kullanımı
- Database performansı

## 🚀 Kurulum ve Başlatma

### Ön Gereksinimler
1. Docker ve Docker Compose kurulu olmalı
2. Ana uygulama çalışıyor olmalı (`docker-compose up` ile)

### Başlatma
```powershell
# Windows PowerShell
cd monitoring
.\start-monitoring.ps1

# veya manuel olarak
docker-compose up -d
```

### Erişim URL'leri
- **Grafana**: http://localhost:3000 (admin/admin123)
- **Prometheus**: http://localhost:9090  
- **Loki**: http://localhost:3100
- **Node Exporter**: http://localhost:9100
- **cAdvisor**: http://localhost:8081

## 📈 Dashboardlar

### Ana Dashboard: "GoAuth Boilerplate - Modern Dashboard"
- 🚀 **Application Overview**: Uptime, request rate, response time
- 🔐 **Authentication Metrics**: Login attempts, success rates
- 💾 **Cache Performance**: Hit/miss ratio, performance
- 🐛 **Error Monitoring**: Error rates by endpoint
- 🔧 **System Resources**: CPU, memory, database connections

## 🚨 Alert Kuralları

### Critical Alerts
- Application down (1 dakika)
- High error rate (>10% for 2 min)

### Warning Alerts  
- High response time (>1s for 3 min)
- High auth failure rate (>5/sec for 1 min)
- Low cache hit ratio (<70% for 5 min)
- High database connections (>20 for 2 min)
- High memory usage (>80% for 5 min)
- High CPU usage (>80% for 5 min)

## 🔧 Konfigürasyon Dosyaları

```
monitoring/
├── docker-compose.yml          # Ana compose dosyası
├── prometheus/
│   ├── prometheus.yml          # Prometheus konfigürasyonu
│   └── alert_rules.yml         # Alert kuralları
├── grafana/
│   └── provisioning/
│       ├── datasources/        # Veri kaynakları
│       └── dashboards/         # Dashboard tanımları
├── loki/
│   └── loki-config.yml         # Loki konfigürasyonu
└── promtail/
    └── promtail-config.yml     # Log collector config
```

## 📊 Custom Metrics

Go uygulamasında aşağıdaki custom metrics'ler tanımlı:

```go
// HTTP Metrics
http_requests_total
http_request_duration_seconds
http_requests_in_flight

// Auth Metrics  
auth_attempts_total
active_sessions

// Cache Metrics
cache_hits_total
cache_misses_total

// Database Metrics
database_connections_active
```

## 🛠️ Yönetim Komutları

```bash
# Servisleri başlat
docker-compose up -d

# Servisleri durdur
docker-compose down

# Logları izle
docker-compose logs -f

# Servis durumunu kontrol et
docker-compose ps

# Metrics'leri kontrol et
curl http://localhost:8080/metrics

# Health check
curl http://localhost:8080/ready
```

## 🎨 Dashboard Özellikleri

- **Modern Dark Theme** - Göz yormayan tasarım
- **Real-time Updates** - 5 saniye refresh
- **Interactive Panels** - Drill-down capability
- **Responsive Design** - Mobile friendly
- **Emoji Icons** - Kolay tanımlama

## 🔍 Troubleshooting

### Grafana'ya erişemiyorum
```bash
docker-compose logs grafana
# Default: admin/admin123
```

### Metrics görünmüyor
```bash
# App metrics endpoint'ini kontrol et
curl http://localhost:8080/metrics

# Prometheus targets'ı kontrol et  
# http://localhost:9090/targets
```

### Loglar gelmiyor
```bash
docker-compose logs promtail
# Log path'leri kontrol et: ../logs/*.log
```

## 📝 Notlar

- **PostgreSQL/Redis Exporters**: Otomatik bağlantı kurulur
- **Log Rotation**: Otomatik log temizleme aktif
- **Data Retention**: Prometheus 30 gün data tutar
- **Resource Limits**: Production'da resource limit'leri ekleyin

## 🔧 Geliştirme

Yeni metrics eklemek için:

1. `internal/adapters/api/middleware/metrics.go` dosyasını düzenleyin
2. Dashboard'a yeni panel ekleyin
3. Gerekirse alert kuralı ekleyin

Örnek metric ekleme:
```go
myCustomMetric = promauto.NewCounterVec(
    prometheus.CounterOpts{
        Name: "my_custom_metric_total",
        Help: "Description of my metric",
    },
    []string{"label1", "label2"},
)
```
