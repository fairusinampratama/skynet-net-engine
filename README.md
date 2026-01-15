# Skynet NetEngine API 🚀

**NetEngine** is a high-performance middleware designed for ISP management. It acts as the "Muscle" between your "Brain" (e.g., Laravel/PHP App) and your MikroTik router fleet.

It maintains persistent, self-healing TCP connections to hundreds of routers, allowing for millisecond-latency commands and real-time monitoring without the overhead of establishing new connections for every request.

## ✨ Features

*   **Persistent Connections**: Dedicated Goroutine per router with automatic reconnection logic.
*   **High Concurrency**: Thread-safe worker pool handling concurrent API requests.
*   **REST API Bridge**: Simple HTTP endpoints for complex RouterOS commands (Sync, Secret, Queues).
*   **Real-time Monitoring**: Fetch CPU, Memory, and live Queue Traffic stats instantly.
*   **Customer Isolation**: One-click isolation via Firewall Address Lists.
*   **Enterprise Security**: API Key authentication (`X-App-Key`) and Webhook event dispatching.
*   **Swagger Documentation**: Built-in interactive API docs.

## 🛠️ Tech Stack

*   **Language**: Go (Golang) 1.22+
*   **Framework**: Gin Gonic
*   **Database**: MySQL (for router credentials)
*   **MikroTik Lib**: `go-routeros`
*   **Logging**: Uber Zap

## 🚀 Deployment

### Option 1: Docker (Recommended)

```bash
# 1. Start Services
docker-compose up -d

# 2. Check Logs
docker-compose logs -f net-engine
```

### Option 2: Manual (Linux/Systemd)

```bash
# 1. Build
go build -o skynet-net-engine-api cmd/server/main.go

# 2. Setup Service
sudo cp net-engine.service /etc/systemd/system/
sudo systemctl enable net-engine
sudo systemctl start net-engine
```

## ⚙️ Configuration

Create a `.env` file (or rely on defaults/docker-compose):

```ini
DB_DSN="username:password@tcp(127.0.0.1:3306)/netengine?parseTime=true"
API_PORT=":8080"
APP_KEY="your_secure_random_key"
```

## 📚 API Documentation

Once running, access the full Swagger UI at:
**[http://localhost:8080/api/v1/swagger/index.html](http://localhost:8080/api/v1/swagger/index.html)**

### Core Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/api/v1/health` | System health check |
| `POST` | `/api/v1/sync/:id` | Force router sync |
| `POST` | `/api/v1/secret` | Create PPPoE secret |
| `POST` | `/api/v1/isolate` | Isolate/Unisolate customer |
| `GET` | `/api/v1/router/:id/health` | Get Router CPU/Mem |
| `GET` | `/api/v1/router/:id/traffic` | Get Live User Traffic |

## 🧪 Development

### Running Tests
```bash
go test ./... -v
```

### Database Seeding
To populate the database with initial router data:
```bash
go run cmd/seeder/main.go
```

## 📝 License
Proprietary / Internal Use Only.
