# Contributing to Task Manager

Thank you for considering contributing to Task Manager! This document provides guidelines and instructions for contributing.

## Code of Conduct

- Be respectful and inclusive
- Provide constructive feedback
- Focus on what is best for the community
- Show empathy towards other community members

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check existing issues. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce**
- **Expected behavior**
- **Actual behavior**
- **Environment details** (OS, Go version, etc.)
- **Logs or error messages**

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion, include:

- **Clear title and description**
- **Use case** - Why is this enhancement needed?
- **Proposed solution**
- **Alternative solutions** you've considered

### Pull Requests

1. **Fork the repository**
2. **Create a feature branch** (`git checkout -b feature/amazing-feature`)
3. **Make your changes**
4. **Test your changes**
5. **Commit with clear messages** (`git commit -m 'Add amazing feature'`)
6. **Push to your branch** (`git push origin feature/amazing-feature`)
7. **Open a Pull Request**

#### Pull Request Guidelines

- Follow the existing code style
- Add tests for new features
- Update documentation as needed
- Keep PRs focused on a single feature/fix
- Write clear commit messages
- Reference related issues

## Development Setup

### Prerequisites

```bash
# Install Go 1.25.2 or higher
go version

# Install Redis (optional, for caching)
# Windows: Download from https://redis.io/download
# Linux: sudo apt-get install redis-server
# macOS: brew install redis
```

### Local Development

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/task-manager.git
cd task-manager

# Install dependencies
go mod download

# Copy environment template
cp .env.example .env

# Edit .env with your settings
# Set JWT_SECRET, ALLOWED_ORIGINS, etc.

# Run the server
go run main.go

# In another terminal, start Redis (if using caching)
redis-server
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with verbose output
go test -v ./...

# Run specific package tests
go test ./internal/handler/...
```

### Code Style

We follow standard Go conventions:

```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Install and run golangci-lint (recommended)
golangci-lint run
```

### Project Structure

```
task-manager/
├── frontend/              # Frontend files (HTML, CSS, JS)
├── internal/
│   ├── handler/          # HTTP request handlers
│   │   ├── auth.go       # Authentication endpoints
│   │   ├── task.go       # Task CRUD endpoints
│   │   └── health.go     # Health check endpoint
│   ├── middleware/       # HTTP middleware
│   │   ├── auth.go       # JWT authentication
│   │   ├── cors.go       # CORS configuration
│   │   ├── logging.go    # Request logging
│   │   ├── ratelimit.go  # Rate limiting
│   │   └── validation.go # Input validation
│   ├── storage/          # Database layer
│   │   ├── sqlite.go     # Task operations
│   │   └── user.go       # User operations
│   ├── types/            # Data models
│   │   ├── task.go       # Task struct
│   │   └── user.go       # User struct
│   └── utils/            # Utilities
│       └── jwt.go        # JWT token management
├── main.go               # Application entry point
├── go.mod                # Go dependencies
├── go.sum                # Dependency checksums
├── .env.example          # Environment template
├── .gitignore            # Git ignore rules
├── README.md             # Project documentation
├── SECURITY.md           # Security documentation
├── CHANGELOG.md          # Version history
└── CONTRIBUTING.md       # This file
```

## Coding Standards

### Go Best Practices

1. **Error Handling**
   ```go
   // Good
   if err != nil {
       log.Printf("Error doing something: %v", err)
       return fmt.Errorf("failed to do something: %w", err)
   }
   
   // Bad
   if err != nil {
       panic(err)
   }
   ```

2. **Input Validation**
   ```go
   // Always sanitize and validate user input
   username = middleware.SanitizeString(username)
   if len(username) < 3 || len(username) > 50 {
       return errors.New("invalid username length")
   }
   ```

3. **Context Usage**
   ```go
   // Use request context, not background context
   ctx := r.Context()
   result, err := db.QueryContext(ctx, query, args...)
   ```

4. **Logging**
   ```go
   // Use structured logging
   log.Printf("Action completed | user: %d | resource: %s", userID, resourceID)
   
   // Don't expose sensitive data
   log.Printf("User logged in: %s", username) // OK
   log.Printf("Password: %s", password)       // NEVER
   ```

### Security Guidelines

1. **Never hardcode secrets**
   ```go
   // Good
   secret := os.Getenv("JWT_SECRET")
   
   // Bad
   secret := "hardcoded-secret"
   ```

2. **Always use parameterized queries**
   ```go
   // Good
   db.Query("SELECT * FROM users WHERE id = ?", userID)
   
   // Bad
   db.Query(fmt.Sprintf("SELECT * FROM users WHERE id = %d", userID))
   ```

3. **Validate and sanitize all inputs**
   ```go
   // Always sanitize before using
   input = middleware.SanitizeString(input)
   ```

4. **Use proper error messages**
   ```go
   // Good - Generic message
   http.Error(w, "Invalid credentials", http.StatusUnauthorized)
   
   // Bad - Reveals information
   http.Error(w, "User not found", http.StatusUnauthorized)
   ```

## Testing Guidelines

### Writing Tests

```go
func TestCreateUser(t *testing.T) {
    // Setup
    db := setupTestDB(t)
    defer db.Close()
    store := storage.New(db)
    
    // Test
    user := types.User{
        Username: "testuser",
        Password: "hashedpassword",
    }
    
    id, err := store.CreateUser(user)
    
    // Assert
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if id == 0 {
        t.Error("Expected valid ID, got 0")
    }
}
```

### Test Coverage

- Aim for >80% code coverage
- Test happy paths and error cases
- Test edge cases and boundary conditions
- Mock external dependencies (Redis, etc.)

## Documentation

### Code Comments

```go
// CreateUser inserts a new user into the database with validation.
// Returns the user ID on success or an error if validation fails or
// the username already exists.
func (s *Store) CreateUser(u types.User) (int64, error) {
    // Implementation
}
```

### API Documentation

When adding new endpoints, update README.md with:
- Endpoint path and method
- Request body format
- Response format
- Status codes
- Authentication requirements
- Example usage

## Commit Message Guidelines

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Code style changes (formatting, etc.)
- **refactor**: Code refactoring
- **test**: Adding or updating tests
- **chore**: Maintenance tasks

### Examples

```
feat(auth): add password reset functionality

Implement password reset via email with secure tokens.
Tokens expire after 1 hour.

Closes #123
```

```
fix(tasks): prevent unauthorized task deletion

Add user_id check in DeleteTask to ensure users can only
delete their own tasks.

Fixes #456
```

## Release Process

1. Update CHANGELOG.md
2. Update version in documentation
3. Create git tag (`git tag -a v2.1.0 -m "Release v2.1.0"`)
4. Push tag (`git push origin v2.1.0`)
5. Create GitHub release with changelog

## Questions?

- Open an issue for questions
- Check existing issues and documentation
- Reach out to maintainers

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

Thank you for contributing to Task Manager! 🚀
