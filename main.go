package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "modernc.org/sqlite"

	"github.com/aadityya4real/task-manager/internal/handler"
	"github.com/aadityya4real/task-manager/internal/middleware"
	"github.com/aadityya4real/task-manager/internal/storage"

	"github.com/redis/go-redis/v9"
)

func main() {
	log.Println("🚀 Starting Task Manager Server...")

	// Database connection
	db, err := sql.Open("sqlite", "tasks.db")
	if err != nil {
		log.Fatalf("❌ Failed to open database: %v", err)
	}
	defer db.Close()

	// Set connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("❌ Failed to connect to database: %v", err)
	}
	log.Println("✅ Database connected")

	// Create tables with indexes
	if err := initDatabase(db); err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}

	// Redis setup
	var rdb *redis.Client
	redisURL := os.Getenv("REDIS_URL")

	if redisURL != "" {
		opt, err := redis.ParseURL(redisURL)
		if err != nil {
			log.Fatalf("❌ Failed to parse Redis URL: %v", err)
		}
		rdb = redis.NewClient(opt)
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr:         "localhost:6379",
			Password:     "",
			DB:           0,
			DialTimeout:  5 * time.Second,
			ReadTimeout:  3 * time.Second,
			WriteTimeout: 3 * time.Second,
		})
	}

	// Test Redis connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Printf("⚠️ Redis connection failed: %v (caching will be disabled)", err)
	} else {
		log.Println("✅ Redis connected")
	}

	// Initialize store
	store := storage.New(db)

	// Initialize rate limiter (10 requests per minute for auth endpoints)
	authRateLimiter := middleware.NewRateLimiter(10, time.Minute)

	// Setup router
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("/health", handler.HealthHandler(db, rdb))

	// Serve frontend
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "frontend/index.html")
			return
		}
		http.NotFound(w, r)
	})

	// Static files
	fs := http.FileServer(http.Dir("frontend"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// API routes with rate limiting and validation
	mux.HandleFunc("/signup",
		authRateLimiter.Limit(
			middleware.ValidateContentType(
				handler.SignupHandler(store),
			),
		),
	)

	mux.HandleFunc("/login",
		authRateLimiter.Limit(
			middleware.ValidateContentType(
				handler.LoginHandler(store),
			),
		),
	)

	mux.HandleFunc("/tasks",
		middleware.AuthMiddleware(
			handler.TaskHandler(store, rdb),
		),
	)

	// Wrap with CORS and logging middleware
	finalHandler := middleware.Logging(middleware.CORS(mux))

	// Get port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Create server with timeouts
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      finalHandler,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		log.Printf("🌍 Server running on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("❌ Server failed: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("⚠️ Server forced to shutdown: %v", err)
	}

	// Close Redis connection
	if err := rdb.Close(); err != nil {
		log.Printf("⚠️ Error closing Redis: %v", err)
	}

	log.Println("✅ Server stopped gracefully")
}

func initDatabase(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return err
	}

	// Create index on username
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_username ON users(username)`)
	if err != nil {
		return err
	}

	// Create tasks table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS tasks (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			title TEXT NOT NULL,
			done BOOLEAN DEFAULT 0,
			user_id INTEGER NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
		)
	`)
	if err != nil {
		return err
	}

	// Create index on user_id for faster queries
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tasks_user_id ON tasks(user_id)`)
	if err != nil {
		return err
	}

	// Create composite index for user_id and done status
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_tasks_user_done ON tasks(user_id, done)`)
	if err != nil {
		return err
	}

	log.Println("✅ Database tables and indexes created")
	return nil
}
