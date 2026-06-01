package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
)

type HealthStatus struct {
	Status   string            `json:"status"`
	Database string            `json:"database"`
	Redis    string            `json:"redis"`
	Uptime   string            `json:"uptime"`
	Details  map[string]string `json:"details,omitempty"`
}

var startTime = time.Now()

// HealthHandler provides a health check endpoint
func HealthHandler(db *sql.DB, rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
		defer cancel()

		health := HealthStatus{
			Status:  "healthy",
			Uptime:  time.Since(startTime).String(),
			Details: make(map[string]string),
		}

		// Check database
		if err := db.PingContext(ctx); err != nil {
			health.Database = "unhealthy"
			health.Status = "degraded"
			health.Details["database_error"] = err.Error()
		} else {
			health.Database = "healthy"
		}

		// Check Redis
		if err := rdb.Ping(ctx).Err(); err != nil {
			health.Redis = "unhealthy"
			health.Status = "degraded"
			health.Details["redis_error"] = err.Error()
		} else {
			health.Redis = "healthy"
		}

		statusCode := http.StatusOK
		if health.Status == "degraded" {
			statusCode = http.StatusServiceUnavailable
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(health)
	}
}
