# Skynet NetEngine API 🚀

**NetEngine** is a high-performance middleware designed for ISP management. It acts as the "Muscle" between your "Brain" (e.g., Laravel/PHP App) and your MikroTik router fleet.

It maintains persistent, self-healing TCP connections to hundreds of routers, allowing for millisecond-latency commands and real-time monitoring without the overhead of establishing new connections for every request.

## ✨ Features

### Backend
*   **Persistent Connections**: Dedicated Goroutine per router with automatic reconnection logic
*   **High Concurrency**: Thread-safe worker pool with serialized router operations
*   **Startup Warmup**: Ensures data is cached before serving traffic (eliminates race conditions)
*   **Smart Queue Detection**: Automatic fallback for PPPoE queue naming patterns
*   **REST API Bridge**: Simple HTTP endpoints for complex RouterOS commands
*   **Real-time Monitoring**: CPU, Memory, and live Queue Traffic stats with 10s timeout
*   **Enterprise Security**: API Key authentication (`X-App-Key`) and Webhook event dispatching
*   **Swagger Documentation**: Built-in interactive API docs at `/api/v1/swagger/index.html`

### Dashboard (Frontend)
> **Note**: The React frontend has been moved to the `feature/full-stack` branch. The `main` branch is now a dedicated Go backend API.

To run the full stack (Backend + Frontend), checkout the branch:
```bash
git checkout feature/full-stack
```

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

## 🖥️ Dashboard

Access the web dashboard at: **[http://localhost:8080](http://localhost:8080)**

Features:
- Real-time traffic graphs with click-to-monitor
- System resource monitoring (CPU/Memory/Uptime)
- Active user sessions with pagination
- Multi-router selector dropdown
- Automatic skeleton loaders during data fetch

## 📚 API Documentation

Once running, access the full Swagger UI at:
**[http://localhost:8080/api/v1/swagger/index.html](http://localhost:8080/api/v1/swagger/index.html)**

### Core Endpoints

| Method | Endpoint | Description |
| :--- | :--- | :--- |
| `GET` | `/` | Web Dashboard (React App) |
| `GET` | `/api/v1/health` | System health check |
| `GET` | `/api/v1/routers` | List all configured routers |
| `GET` | `/api/v1/monitoring/targets` | Get all active users |
| `GET` | `/api/v1/router/:id/health` | Router CPU/Memory stats |
| `GET` | `/api/v1/router/:id/traffic?user=USERNAME` | Live user traffic (bits/sec) |
| `POST` | `/api/v1/sync/:id` | Force router sync |
| `POST` | `/api/v1/secret` | Create PPPoE secret |
| `POST` | `/api/v1/isolate` | Isolate/Unisolate customer |

**Authentication**: All `/api/v1/*` routes require header: `X-App-Key: netengine_secret_key_123`

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

## 🏗️ Architecture

### Backend Flow
```
HTTP Request → Gin Router → API Handler → Worker Pool → MikroTik Router
                                              ↓
                                         Command Queue (Buffered Channel)
                                              ↓
                                         Worker Goroutine (Serialized)
                                              ↓
                                         RouterOS API Response
```

### Key Improvements
- **Warmup Phase**: Server blocks until routers connect and cache initial data
- **Thread Safety**: All router operations serialized through command queue
- **PPPoE Detection**: Smart fallback for `<pppoe-USERNAME>` queue naming
- **Error Resilience**: Frontend treats 500-504 errors as transient loading states

## 📝 License
Proprietary / Internal Use Only.
