# Changelog

All notable changes to this project will be documented in this file.

## [2.0.0] - 2026-06-02

### 🔒 Security Enhancements

#### Critical Fixes
- **JWT Secret Management**: Moved hardcoded secrets to environment variables with centralized management
- **Rate Limiting**: Added rate limiting (10 req/min) on authentication endpoints to prevent brute force attacks
- **CORS Configuration**: Replaced wildcard CORS with configurable origin whitelist
- **Input Validation**: Implemented comprehensive input sanitization and validation across all endpoints
- **Error Handling**: Prevented information leakage through generic error messages

#### Authentication & Authorization
- Added username enumeration protection in login endpoint
- Implemented minimum password length requirement (8 characters)
- Added username uniqueness validation
- Centralized JWT token validation logic
- Enhanced password validation before hashing

#### Database Security
- Added foreign key constraints with CASCADE delete
- Created performance indexes:
  - `idx_users_username` for fast username lookups
  - `idx_tasks_user_id` for user-scoped queries
  - `idx_tasks_user_done` composite index
- Improved error handling in all database operations
- Added connection pooling with limits (25 max open, 5 idle)
- Enhanced user-scoped queries to prevent unauthorized access

### 🚀 New Features

#### Middleware
- **Logging Middleware**: Request/response logging with timestamps and status codes
- **Rate Limiting Middleware**: Configurable per-IP rate limiting with automatic cleanup
- **CORS Middleware**: Environment-based origin whitelisting
- **Validation Middleware**: Content-Type validation and input sanitization

#### Endpoints
- **Health Check** (`GET /health`): Monitor database and Redis connectivity
  - Returns service status
  - Checks database ping
  - Checks Redis connectivity
  - Reports uptime

#### Infrastructure
- **Graceful Shutdown**: Proper cleanup of resources on SIGINT/SIGTERM
- **Connection Timeouts**: Configured timeouts for database and Redis
- **Server Timeouts**: Read (15s), Write (15s), Idle (60s) timeouts

### 📝 Code Quality Improvements

#### Error Handling
- Replaced `fmt.Println` with proper `log` package
- Added context to all error messages
- Implemented proper error wrapping with `%w`
- Added error checking for `rows.Scan()` operations

#### Input Validation
- Username: 3-50 characters, trimmed, sanitized
- Password: minimum 8 characters, sanitized
- Task title: maximum 500 characters, sanitized
- Null byte removal from all inputs

#### Logging
- Structured logging with consistent format
- Request method, URI, status code, duration, and IP
- Error logging with context
- Cache hit/miss logging for debugging

### 🗑️ Removed

- **Unused Code**: Removed `internal/storage/redis.go` (Redis initialized in main.go)
- **Duplicate Secret**: Removed duplicate JWT secret definition in middleware

### 📚 Documentation

#### New Files
- **SECURITY.md**: Comprehensive security documentation
- **CHANGELOG.md**: This file
- **.env.example**: Environment variable template
- **README.md**: Complete rewrite with:
  - Feature list
  - API documentation
  - Security features
  - Configuration guide
  - Deployment instructions

#### Updated Files
- **.gitignore**: Added .env, database files, logs, and IDE files

### 🔧 Configuration

#### Environment Variables
- `JWT_SECRET`: JWT signing key (required in production)
- `ALLOWED_ORIGINS`: Comma-separated list of allowed CORS origins
- `REDIS_URL`: Redis connection URL (optional)
- `PORT`: Server port (default: 8080)

### 🐛 Bug Fixes

- Fixed missing error handling in `GetTasks()` row scanning
- Fixed Redis connection not being tested at startup
- Fixed database connection not being closed on shutdown
- Fixed potential race conditions in rate limiter
- Fixed cache invalidation for paginated results

### ⚡ Performance Improvements

- Added database indexes for faster queries
- Implemented connection pooling with optimal settings
- Added pagination-aware Redis caching
- Optimized cache invalidation strategy

### 📦 Dependencies

No new dependencies added. All fixes use existing libraries:
- `golang.org/x/crypto` - Password hashing
- `github.com/golang-jwt/jwt/v5` - JWT tokens
- `github.com/redis/go-redis/v9` - Redis client
- `modernc.org/sqlite` - SQLite driver

### 🔄 Migration Notes

#### Breaking Changes
None - All changes are backward compatible with existing API contracts.

#### Recommended Actions
1. Create `.env` file from `.env.example`
2. Set strong `JWT_SECRET` in production
3. Configure `ALLOWED_ORIGINS` for your frontend
4. Review rate limiting settings for your use case
5. Set up health check monitoring

#### Database Migration
No manual migration needed. New indexes and constraints are created automatically on startup.

### 📊 Statistics

- **Files Modified**: 12
- **Files Added**: 8
- **Files Removed**: 1
- **Security Issues Fixed**: 10
- **New Features**: 5
- **Performance Improvements**: 4

---

## [1.0.0] - 2026-04-13

### Initial Release
- Basic task management API
- User authentication with JWT
- SQLite database
- Redis caching
- CRUD operations for tasks
- Frontend interface

---

**Note**: Version 2.0.0 represents a major security and quality overhaul while maintaining API compatibility.
