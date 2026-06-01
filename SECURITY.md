# Security Documentation

This document outlines the security measures implemented in the Task Manager API.

## Security Fixes Applied

### 1. JWT Secret Management ✅
**Issue:** Hardcoded JWT secret in multiple files  
**Fix:** 
- Centralized secret management in `utils/jwt.go`
- Environment variable support via `JWT_SECRET`
- Falls back to default only in development
- Single source of truth for token generation and validation

**Files Modified:**
- `internal/utils/jwt.go` - Added `GetSecretKey()` and `ValidateToken()`
- `internal/middleware/auth.go` - Uses centralized validation

### 2. Input Validation & Sanitization ✅
**Issue:** No input validation or sanitization  
**Fix:**
- Created `middleware/validation.go` with sanitization functions
- Removes null bytes and trims whitespace
- Length validation on all inputs
- Content-Type validation for POST/PUT requests

**Validation Rules:**
- Username: 3-50 characters, unique
- Password: minimum 8 characters
- Task title: maximum 500 characters
- All inputs sanitized before processing

### 3. Rate Limiting ✅
**Issue:** No protection against brute force attacks  
**Fix:**
- Implemented `middleware/ratelimit.go`
- 10 requests per minute on auth endpoints
- Per-IP tracking with automatic cleanup
- Returns 429 (Too Many Requests) when exceeded

### 4. CORS Configuration ✅
**Issue:** Wide-open CORS (`Access-Control-Allow-Origin: *`)  
**Fix:**
- Created `middleware/cors.go` with configurable origins
- Environment variable `ALLOWED_ORIGINS` for whitelist
- No wildcard in production
- Proper preflight handling

### 5. Error Handling & Logging ✅
**Issue:** Using `fmt.Println` and exposing internal errors  
**Fix:**
- Replaced with proper `log` package
- Generic error messages to prevent information leakage
- Structured logging with timestamps
- Request/response logging middleware

**Security Improvements:**
- Login errors don't reveal if username exists (prevents enumeration)
- Internal errors logged but not exposed to clients
- All database errors wrapped with context

### 6. Database Security ✅
**Issue:** Missing error handling, no indexes, no constraints  
**Fix:**
- Added proper error handling in all queries
- Created indexes for performance and security:
  - `idx_users_username` - Fast username lookups
  - `idx_tasks_user_id` - User-scoped queries
  - `idx_tasks_user_done` - Composite index
- Foreign key constraints with CASCADE delete
- User-scoped queries prevent unauthorized access
- Connection pooling with limits

### 7. Password Security ✅
**Already Implemented (Verified):**
- Bcrypt hashing with default cost (10)
- No plaintext password storage
- Password not returned in API responses

**Added:**
- Minimum 8 character requirement
- Password validation before hashing

### 8. Redis Connection Handling ✅
**Issue:** No connection verification, silent failures  
**Fix:**
- Connection test at startup with timeout
- Graceful degradation if Redis unavailable
- Proper error logging
- Connection timeouts configured

### 9. Graceful Shutdown ✅
**Issue:** No cleanup on shutdown  
**Fix:**
- Signal handling (SIGINT, SIGTERM)
- 10-second graceful shutdown timeout
- Proper resource cleanup (DB, Redis)
- In-flight requests completed before shutdown

### 10. Health Check Endpoint ✅
**Issue:** No monitoring capability  
**Fix:**
- Created `/health` endpoint
- Checks database connectivity
- Checks Redis connectivity
- Returns degraded status if issues detected
- Useful for load balancers and monitoring

## Security Best Practices

### For Production Deployment

1. **Environment Variables** (CRITICAL)
   ```bash
   # Generate a strong secret (example)
   JWT_SECRET=$(openssl rand -base64 32)
   
   # Set allowed origins
   ALLOWED_ORIGINS=https://yourdomain.com,https://app.yourdomain.com
   
   # Redis URL (if using external Redis)
   REDIS_URL=redis://user:password@redis-host:6379
   ```

2. **HTTPS Only**
   - Use a reverse proxy (nginx, Caddy) for TLS termination
   - Set `Secure` flag on cookies if using them
   - Enable HSTS headers

3. **Database**
   - Regular backups of `tasks.db`
   - File permissions: `chmod 600 tasks.db`
   - Consider PostgreSQL for production

4. **Redis**
   - Enable authentication (`requirepass`)
   - Use TLS for Redis connections
   - Limit network access

5. **Rate Limiting**
   - Consider using a reverse proxy for global rate limiting
   - Adjust limits based on your use case
   - Monitor for abuse patterns

6. **Monitoring**
   - Set up health check monitoring
   - Log aggregation and analysis
   - Alert on repeated auth failures

## Remaining Considerations

### Not Implemented (Future Enhancements)

1. **Email Verification** - Users can sign up without email verification
2. **Password Reset** - No password recovery mechanism
3. **Account Lockout** - No temporary lockout after failed attempts
4. **2FA/MFA** - No multi-factor authentication
5. **API Versioning** - No version prefix on endpoints
6. **Request ID Tracking** - No correlation IDs for debugging
7. **Audit Logging** - No audit trail for sensitive operations
8. **Session Management** - JWT tokens can't be revoked before expiry
9. **CSRF Protection** - Not needed for stateless JWT API, but consider if adding cookies
10. **Content Security Policy** - Frontend security headers

### Known Limitations

1. **JWT Expiration** - Tokens valid for 24 hours, can't be revoked
   - Consider shorter expiration + refresh tokens
   - Or implement token blacklist in Redis

2. **Rate Limiting** - In-memory, resets on restart
   - Consider Redis-based rate limiting for multi-instance deployments

3. **SQLite** - Single-writer limitation
   - Fine for small-medium apps
   - Consider PostgreSQL for high concurrency

## Vulnerability Reporting

If you discover a security vulnerability, please email security@yourdomain.com instead of opening a public issue.

## Security Checklist for Deployment

- [ ] Set strong `JWT_SECRET` environment variable
- [ ] Configure `ALLOWED_ORIGINS` with specific domains
- [ ] Enable HTTPS via reverse proxy
- [ ] Set up Redis authentication
- [ ] Configure firewall rules
- [ ] Set proper file permissions on database
- [ ] Enable logging and monitoring
- [ ] Set up automated backups
- [ ] Review and test rate limits
- [ ] Scan dependencies for vulnerabilities (`go list -m all | nancy sleuth`)

## Dependencies Security

Current dependencies are from trusted sources:
- `golang.org/x/crypto` - Official Go crypto library
- `github.com/golang-jwt/jwt/v5` - Well-maintained JWT library
- `github.com/redis/go-redis/v9` - Official Redis client
- `modernc.org/sqlite` - Pure Go SQLite implementation

**Recommendation:** Regularly update dependencies and scan for vulnerabilities.

```bash
# Check for known vulnerabilities
go list -json -m all | nancy sleuth

# Update dependencies
go get -u ./...
go mod tidy
```

## Compliance Notes

- **GDPR**: No personal data collected beyond username
- **Password Storage**: Compliant with OWASP guidelines (bcrypt)
- **Data Retention**: No automatic deletion implemented
- **Right to Deletion**: Implement user deletion endpoint if needed

---

**Last Updated:** June 2, 2026  
**Security Review:** Completed
