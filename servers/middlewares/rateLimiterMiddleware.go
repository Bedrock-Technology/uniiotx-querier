package middlewares

import (
	"golang.org/x/time/rate"
	"net"
	"net/http"
	"sync"
)

// limiterInfo holds the rate limiter and the last time it was used
type limiterInfo struct {
	Limiter *rate.Limiter
}

// ipRateLimiter stores a map of IP addresses to their respective limiter and a mutex
type ipRateLimiter struct {
	mu       sync.RWMutex
	limiters map[string]*limiterInfo
	rate     rate.Limit
	burst    int
}

func newIPRateLimiter(rate rate.Limit, burst int) *ipRateLimiter {
	return &ipRateLimiter{
		limiters: make(map[string]*limiterInfo),
		rate:     rate,
		burst:    burst,
	}
}

func (i *ipRateLimiter) getLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.limiters[ip]
	if !exists {
		limiter = &limiterInfo{
			Limiter: rate.NewLimiter(i.rate, i.burst),
		}
		i.limiters[ip] = limiter
	}

	return limiter.Limiter
}

// middleware function for IP based rate limiting
func (i *ipRateLimiter) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := readUserIP(r) // Custom function to accurately read the user's IP address

		limiter := i.getLimiter(ip)
		if !limiter.Allow() {
			http.Error(w, http.StatusText(http.StatusTooManyRequests), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// readUserIP extracts the IP address from the request, taking into account proxies.
func readUserIP(r *http.Request) string {
	// Check X-Real-IP and X-Forwarded-For headers first
	if ip := r.Header.Get("X-Real-IP"); ip != "" {
		return ip
	}
	if ip := r.Header.Get("X-Forwarded-For"); ip != "" {
		return ip
	}

	// Fallback to using the remote address
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}

// RateLimiterMiddleware creates a new rate limiter that allows requests up to 'r' tokens per second with a maximum burst size of 'b'.
func RateLimiterMiddleware(r rate.Limit, b int) func(http.Handler) http.Handler {
	ipRateLimiter := newIPRateLimiter(r, b)

	return func(next http.Handler) http.Handler {
		return ipRateLimiter.middleware(next)
	}
}
