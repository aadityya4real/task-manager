"# Task Manager

A secure, production-ready task management API built with Go, featuring JWT authentication, Redis caching, and SQLite database.

**Started on:** 04/13/2026  
**Status:** Production Ready ✅

## Features

- ✅ **User Authentication** - Secure signup/login with bcrypt password hashing
- ✅ **JWT Authorization** - Token-based authentication with configurable secrets
- ✅ **Task Management** - Full CRUD operations for tasks
- ✅ **Redis Caching** - Intelligent caching with pagination support
- ✅ **Rate Limiting** - Protection against brute force attacks
- ✅ **Input Validation** - Comprehensive input sanitization and validation
- ✅ **Security Hardening** - CORS configuration, SQL injection protection
- ✅ **Health Checks** - Monitor database and Redis connectivity
- ✅ **Graceful Shutdown** - Proper cleanup of resources
- ✅ **Structured Logging** - Request/response logging with timestamps
- ✅ **Database Indexes** - Optimized queries for performance

## Quick Start

### Prerequisites

- Go 1.25.2 or higher
- Redis (optional, for caching)

### Installation

1. Clone the repository
```bash
git clone https://github.com/aadityya4real/task-manager.git
cd task-manager
```

2. Install dependencies
```bash
go mod download
```

3. Configure environment variables
```bash
cp .env.example .env
# Edit .env with your configuration
```

4. Run the server
```bash
go run main.go
```

The server will start on `http://localhost:8080`

## Configuration

### Environment Variables

| Variable | Description | Default | Required |
|----------|-------------|---------|----------|
| `PORT` | Server port | `8080` | No |
| `JWT_SECRET` | Secret key for JWT tokens | `mysecretkey` | **Yes (production)** |
| `REDIS_URL` | Redis connection URL | `redis://localhost:6379` | No |
| `ALLOWED_ORIGINS` | CORS allowed origins (comma-separated) | `http://localhost:3000,http://localhost:8080` | No |

**⚠️ Security Warning:** Always set a strong `JWT_SECRET` in production!

## API Endpoints

### Health Check
```
GET /health
```
Returns server health status including database and Redis connectivity.

### Authentication

#### Signup
```
POST /signup
Content-Type: application/json

{
  "username": "user123",
  "password": "securepassword"
}
```

**Validation:**
- Username: 3-50 characters
- Password: minimum 8 characters
- Rate limit: 10 requests/minute

#### Login
```
POST /login
Content-Type: application/json

{
  "username": "user123",
  "password": "securepassword"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Rate limit:** 10 requests/minute

### Tasks

All task endpoints require authentication via `Authorization: Bearer <token>` header.

#### Create Task
```
POST /tasks
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Complete project documentation"
}
```

#### Get Tasks (with pagination)
```
GET /tasks?limit=10&offset=0
Authorization: Bearer <token>
```

**Query Parameters:**
- `limit`: Number of tasks per page (default: 10)
- `offset`: Number of tasks to skip (default: 0)

#### Update Task
```
PUT /tasks?id=1
Authorization: Bearer <token>
Content-Type: application/json

{
  "title": "Updated task title",
  "done": true
}
```

#### Delete Task
```
DELETE /tasks?id=1
Authorization: Bearer <token>
```

## Security Features

### Implemented Protections

1. **Password Security**
   - Bcrypt hashing with default cost
   - Minimum 8 character requirement
   - No plaintext password storage

2. **JWT Security**
   - Configurable secret key via environment
   - Token expiration (24 hours)
   - Secure token validation

3. **Rate Limiting**
   - 10 requests/minute on auth endpoints
   - Per-IP tracking
   - Automatic cleanup of old visitors

4. **Input Validation**
   - SQL injection protection via parameterized queries
   - Input sanitization (null byte removal, trimming)
   - Length validation on all inputs
   - Content-Type validation

5. **CORS Configuration**
   - Configurable allowed origins
   - No wildcard (*) in production
   - Preflight request handling

6. **Database Security**
   - Foreign key constraints
   - User-scoped queries (prevents unauthorized access)
   - Connection pooling with limits
   - Proper error handling

## Performance Optimizations

- **Redis Caching**: Pagination-aware caching with 5-minute TTL
- **Database Indexes**: 
  - `idx_users_username` on users.username
  - `idx_tasks_user_id` on tasks.user_id
  - `idx_tasks_user_done` composite index
- **Connection Pooling**: Max 25 open connections, 5 idle
- **Graceful Shutdown**: 10-second timeout for cleanup

## Development

### Project Structure
```
task-manager/
├── frontend/           # Frontend files
├── internal/
│   ├── handler/       # HTTP handlers
│   ├── middleware/    # Middleware (auth, logging, rate limiting)
│   ├── storage/       # Database operations
│   ├── types/         # Data models
│   └── utils/         # Utilities (JWT)
├── main.go            # Application entry point
├── go.mod             # Go dependencies
└── tasks.db           # SQLite database
```

### Running Tests
```bash
go test ./...
```

### Building for Production
```bash
go build -o task-manager main.go
./task-manager
```

## Deployment

### Environment Setup

1. Set strong `JWT_SECRET`
2. Configure `ALLOWED_ORIGINS` with your frontend domain
3. Set up Redis for production caching
4. Use a reverse proxy (nginx) for HTTPS

### Docker (Optional)
```dockerfile
FROM golang:1.25-alpine
WORKDIR /app
COPY . .
RUN go build -o task-manager
CMD ["./task-manager"]
```

## License

MIT License - feel free to use this project for learning or production.

## Contributing

Contributions are welcome! Please open an issue or submit a pull request.

## Author

[@aadityya4real](https://github.com/aadityya4real)

