# Quick Start Guide

Get your Task Manager API up and running in 5 minutes!

## Prerequisites

- Go 1.25.2+ installed
- Redis installed (optional, but recommended)

## Step 1: Clone & Setup

```bash
# Navigate to project directory
cd task-manager

# Install dependencies
go mod download
```

## Step 2: Configure Environment

```bash
# Copy the example environment file
cp .env.example .env
```

Edit `.env` and set your configuration:

```env
# REQUIRED for production
JWT_SECRET=your-super-secret-key-here

# Optional - defaults shown
PORT=8080
REDIS_URL=redis://localhost:6379
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:8080
```

**🔒 Security Tip:** Generate a strong JWT secret:
```bash
# On Linux/macOS
openssl rand -base64 32

# On Windows (PowerShell)
[Convert]::ToBase64String((1..32 | ForEach-Object { Get-Random -Maximum 256 }))
```

## Step 3: Start Redis (Optional)

### Windows
```bash
# Download from https://redis.io/download
# Or use WSL: wsl redis-server
```

### Linux
```bash
sudo systemctl start redis
# or
redis-server
```

### macOS
```bash
brew services start redis
# or
redis-server
```

**Note:** The app works without Redis, but caching will be disabled.

## Step 4: Run the Server

```bash
go run main.go
```

You should see:
```
🚀 Starting Task Manager Server...
✅ Database connected
✅ Database tables and indexes created
✅ Redis connected
🌍 Server running on port 8080
```

## Step 5: Test the API

### Health Check
```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "database": "healthy",
  "redis": "healthy",
  "uptime": "5.2s"
}
```

### Create a User
```bash
curl -X POST http://localhost:8080/signup \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

### Login
```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","password":"password123"}'
```

Save the token from the response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Create a Task
```bash
curl -X POST http://localhost:8080/tasks \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"title":"My first task"}'
```

### Get Tasks
```bash
curl http://localhost:8080/tasks \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

### Update a Task
```bash
curl -X PUT "http://localhost:8080/tasks?id=1" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE" \
  -d '{"title":"Updated task","done":true}'
```

### Delete a Task
```bash
curl -X DELETE "http://localhost:8080/tasks?id=1" \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

## Step 6: Access the Frontend

Open your browser and navigate to:
```
http://localhost:8080
```

The frontend interface will be served automatically.

## Common Issues

### Port Already in Use
```bash
# Change the port in .env
PORT=3000
```

### Redis Connection Failed
```
⚠️ Redis connection failed: dial tcp [::1]:6379: connect: connection refused
```

**Solution:** 
- Start Redis server, or
- The app will continue without caching

### Database Locked
```
database is locked
```

**Solution:** 
- Close any other connections to `tasks.db`
- Delete `tasks.db-shm` and `tasks.db-wal` files

### JWT Token Invalid
```
Invalid token
```

**Solution:**
- Check if token is expired (24 hour lifetime)
- Ensure `JWT_SECRET` hasn't changed
- Login again to get a new token

## Development Tips

### Hot Reload
Install `air` for automatic reloading:
```bash
go install github.com/cosmtrek/air@latest
air
```

### View Logs
All requests are logged with:
- Method
- URI
- Status code
- Duration
- Client IP

Example:
```
2026/06/02 10:30:45 POST /login 200 45.2ms 127.0.0.1:54321
```

### Database Inspection
Use SQLite browser to inspect the database:
```bash
# Install sqlite3
# Windows: Download from https://sqlite.org/download.html
# Linux: sudo apt-get install sqlite3
# macOS: brew install sqlite3

# Open database
sqlite3 tasks.db

# View tables
.tables

# Query users
SELECT * FROM users;

# Query tasks
SELECT * FROM tasks;

# Exit
.quit
```

### Clear Cache
```bash
# Connect to Redis
redis-cli

# Clear all cache
FLUSHALL

# Or clear specific user's cache
KEYS tasks:1:*
DEL tasks:1:10:0 tasks:1:10:10
```

## Building for Production

```bash
# Build binary
go build -o task-manager main.go

# Run binary
./task-manager
```

### Linux/macOS
```bash
# Build
GOOS=linux GOARCH=amd64 go build -o task-manager-linux main.go

# Make executable
chmod +x task-manager-linux

# Run
./task-manager-linux
```

### Windows
```bash
# Build
GOOS=windows GOARCH=amd64 go build -o task-manager.exe main.go

# Run
task-manager.exe
```

## Next Steps

- Read [README.md](README.md) for full documentation
- Check [SECURITY.md](SECURITY.md) for security best practices
- Review [CONTRIBUTING.md](CONTRIBUTING.md) if you want to contribute
- See [CHANGELOG.md](CHANGELOG.md) for version history

## Need Help?

- Check the [README.md](README.md) for detailed documentation
- Open an issue on GitHub
- Review existing issues for solutions

---

**Happy coding! 🚀**
