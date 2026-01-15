# Skynet NetEngine Cheat Sheet 🏃

## 1. Production Mode (The "Single Binary")
This is how you run it on your server. It serves both the API and the React Dashboard.

```bash
# Build (if not already built)
go build -o skynet-net-engine-api cmd/server/main.go

# Run
./skynet-net-engine-api
```
👉 **Access Dashboard**: http://localhost:8080
👉 **Access API**: http://localhost:8080/api/v1/health

---

## 2. Development Mode (The "Hacker Way")
Use this if you want to edit the React code and see changes instantly.

**Terminal 1 (Backend)**:
```bash
# Run the API with CORS enabled for dev
go run cmd/server/main.go
```

**Terminal 2 (Frontend)**:
```bash
cd web
npm run dev
```
👉 **Access Dashboard (Hot Reload)**: http://localhost:5173

---

## 3. Troubleshooting

**"Port 8080 already in use"** / **"Failed to start API server"**
This means an old version of the app is still running in the background.

**Kill all running instances:**
```bash
pkill -f skynet-net-engine-api
```
*(This stops both the binary and the `go run` process)*
