package middleware

import (
	"net/http"
	"sync"
	"time"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	window   time.Duration
}

type visitor struct {
	count      int
	lastReset  time.Time
	mu         sync.Mutex
}

func NewRateLimiter(requestsPerWindow int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     requestsPerWindow,
		window:   window,
	}

	// Cleanup old visitors every 5 minutes
	go func() {
		for {
			time.Sleep(5 * time.Minute)
			rl.cleanup()
		}
	}()

	return rl
}

func (rl *rateLimiter) cleanup() {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	for ip, v := range rl.visitors {
		v.mu.Lock()
		if time.Since(v.lastReset) > rl.window*2 {
			delete(rl.visitors, ip)
		}
		v.mu.Unlock()
	}
}

func (rl *rateLimiter) getVisitor(ip string) *visitor {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		v = &visitor{
			count:     0,
			lastReset: time.Now(),
		}
		rl.visitors[ip] = v
	}

	return v
}

func (rl *rateLimiter) Limit(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := r.RemoteAddr

		v := rl.getVisitor(ip)
		v.mu.Lock()

		// Reset counter if window has passed
		if time.Since(v.lastReset) > rl.window {
			v.count = 0
			v.lastReset = time.Now()
		}

		// Check if rate limit exceeded
		if v.count >= rl.rate {
			v.mu.Unlock()
			http.Error(w, "Rate limit exceeded. Try again later.", http.StatusTooManyRequests)
			return
		}

		v.count++
		v.mu.Unlock()

		next.ServeHTTP(w, r)
	}
}
