# NetEngine Cheatsheet

## Quick Start

```bash
# Navigate to project
cd ~/golang/skynet-net-engine-api

# Run server
./skynet-net-engine-api

# Or build and run
go build -o skynet-net-engine-api cmd/server/main.go
./skynet-net-engine-api
```

## Database Setup

```bash
# Create database
mysql -u fairusinampratama -p
CREATE DATABASE netengine;

# Run migrations
mysql -u fairusinampratama netengine < schema.sql
mysql -u fairusinampratama netengine < schema_users.sql

# Seed routers
go run cmd/seeder/main.go

# Sync users from MikroTik (imports all PPPoE accounts)
go run cmd/sync-users/main.go
```

## Development

```bash
# Frontend development (hot reload)
cd web
npm install
npm run dev  # Runs on http://localhost:5173

# Build frontend for production
npm run build

# Deploy frontend to Go binary
cp -r web/dist/* internal/assets/dist/
go build -o skynet-net-engine-api cmd/server/main.go
```

## API Endpoints

### Dashboard
- `GET /` - Web dashboard (React app)

### Monitoring
- `GET /api/v1/routers` - List all routers
- `GET /api/v1/router/:id/health` - CPU, Memory, Uptime
- `GET /api/v1/router/:id/users` - **All users with status** (connected/offline)
- `GET /api/v1/router/:id/traffic?user=USERNAME` - Live traffic (bits/sec)
- `GET /api/v1/monitoring/targets` - Active sessions only

### Management
- `POST /api/v1/secret` - Create PPPoE account
- `POST /api/v1/isolate` - Isolate/unisolate customer
- `POST /api/v1/router/:id/backup` - Trigger config backup

**Auth**: All `/api/v1/*` routes require header: `X-App-Key: netengine_secret_key_123`

## Testing

```bash
# Health check
curl http://localhost:8080/api/v1/health

# Get all users (with auth)
curl "http://localhost:8080/api/v1/router/1/users" \
  -H "X-App-Key: netengine_secret_key_123"

# Get live traffic
curl "http://localhost:8080/api/v1/router/1/traffic?user=USERNAME" \
  -H "X-App-Key: netengine_secret_key_123"
```

## Dashboard Features

### Active Sessions Table
- **Column Sorting**: Click headers to sort by User, IP, or Status
- **Advanced Pagination**: 
  - Entry counter (e.g., "Showing 1-10 of 538")
  - Page size selector (10, 25, 50, 100 per page)
  - Quick navigation (First, Previous, Next, Last)
- **Status Badges**: 
  - 🟢 Connected (green)
  - ⚫ Offline (gray)
- **Click-to-Monitor**: Click any user to view live traffic graph

### Live Traffic Widget
- Real-time bandwidth graph (updates every 1s)
- Shows Download (blue) and Upload (green)
- Skeleton loader during data fetch
- Supports both connected and offline users (shows 0 for offline)

### System Health Widget
- CPU usage percentage
- Memory usage (used/total)
- System uptime
- Auto-refresh every 30s

## Architecture

### Backend Flow
```
HTTP Request → Gin Router → API Handler → Worker Pool → MikroTik
                                              ↓
                                         Command Queue
                                              ↓
                                         Worker (Serialized)
```

### Key Features
- **Startup Warmup**: Server blocks until routers connect and cache data
- **Thread Safety**: All router operations serialized through command queue
- **PPPoE Detection**: Smart fallback for `<pppoe-USERNAME>` queue naming
- **Error Resilience**: Frontend handles transient errors (500-504) gracefully

## Troubleshooting

```bash
# Check if server is running
ps aux | grep skynet-net-engine-api

# Kill existing process
pkill -f skynet-net-engine-api

# View logs (if running in background)
tail -f /var/log/skynet-net-engine.log

# Check database connection
mysql -u fairusinampratama netengine -e "SELECT COUNT(*) FROM pppoe_users;"
```

## Git Workflow

```bash
# Check status
git status

# Stage and commit
git add .
git commit -m "feat: description"

# Push to GitHub
git push
```
