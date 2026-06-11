package middleware

import (
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	mu       sync.Mutex
	requests map[string]int
	limit    int
	window   time.Duration
}

func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	return &RateLimiter{requests: make(map[string]int), limit: limit, window: window}
}

func (rl *RateLimiter) Allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	rl.requests[key]++
	if rl.requests[key] > rl.limit {
		return false
	}
	go func() {
		time.Sleep(rl.window)
		rl.mu.Lock()
		rl.requests[key]--
		if rl.requests[key] <= 0 {
			delete(rl.requests, key)
		}
		rl.mu.Unlock()
	}()
	return true
}

var coachLimiter = NewRateLimiter(10, time.Minute)

func Chain(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("Referrer-Policy", "no-referrer")
		w.Header().Set("Permissions-Policy", "geolocation=(), microphone=(), camera=()")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if r.URL.Path == "/api/coach" && r.Method == http.MethodPost {
			ip := r.RemoteAddr
			if !coachLimiter.Allow(ip) {
				http.Error(w, `{"success":false,"message":"Too many requests. Please wait before asking Wizz again.","error":"rate_limited"}`, http.StatusTooManyRequests)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
