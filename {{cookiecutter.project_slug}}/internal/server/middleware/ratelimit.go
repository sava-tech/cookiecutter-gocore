package middleware

import (
	"math"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/{{ cookiecutter.github_username }}/{{ cookiecutter.project_name }}/pkg/response"
	"golang.org/x/time/rate"
)

// idleTTL is how long a client's limiter is kept after its last request.
// Once idle for this long its bucket would be full again anyway, so
// dropping it loses nothing and keeps memory bounded.
const idleTTL = 10 * time.Minute

type limiterEntry struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter hands out one token-bucket limiter per key (e.g. client IP)
// and periodically evicts limiters that have gone idle.
//
// State is in-memory, so limits apply per process: with N replicas the
// effective limit is N times higher. Move this to Redis if that matters.
type RateLimiter struct {
	mu      sync.Mutex
	entries map[string]*limiterEntry
	limit   rate.Limit
	burst   int
}

// NewRateLimiter creates a limiter allowing `limit` events per second with
// the given burst, and starts a background goroutine that evicts idle keys.
func NewRateLimiter(limit rate.Limit, burst int) *RateLimiter {
	rl := &RateLimiter{
		entries: make(map[string]*limiterEntry),
		limit:   limit,
		burst:   burst,
	}
	go rl.cleanupLoop()
	return rl
}

func (rl *RateLimiter) get(key string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry, exists := rl.entries[key]
	if !exists {
		entry = &limiterEntry{limiter: rate.NewLimiter(rl.limit, rl.burst)}
		rl.entries[key] = entry
	}
	entry.lastSeen = time.Now()
	return entry.limiter
}

func (rl *RateLimiter) cleanupLoop() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		cutoff := time.Now().Add(-idleTTL)
		rl.mu.Lock()
		for key, entry := range rl.entries {
			if entry.lastSeen.Before(cutoff) {
				delete(rl.entries, key)
			}
		}
		rl.mu.Unlock()
	}
}

// retryAfterSeconds is roughly how long until one token is available again.
func (rl *RateLimiter) retryAfterSeconds() int {
	if rl.limit <= 0 {
		return 60
	}
	return int(math.Ceil(1 / float64(rl.limit)))
}

// Middleware limits requests per client IP across every route it wraps.
//
// c.ClientIP() only honours X-Forwarded-For from proxies configured via
// router.SetTrustedProxies, so clients can't spoof their way around this.
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return rl.handler(func(c *gin.Context) string {
		return c.ClientIP()
	})
}

// PerRouteMiddleware limits requests per client IP *and* route, so each
// endpoint it wraps gets its own budget (hitting /login doesn't use up
// /forgot-password's allowance).
func (rl *RateLimiter) PerRouteMiddleware() gin.HandlerFunc {
	return rl.handler(func(c *gin.Context) string {
		return c.ClientIP() + "|" + c.FullPath()
	})
}

func (rl *RateLimiter) handler(keyFn func(*gin.Context) string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if !rl.get(keyFn(c)).Allow() {
			c.Header("Retry-After", strconv.Itoa(rl.retryAfterSeconds()))
			response.TooManyRequests(c, "Too many requests, slow down.")
			c.Abort()
			return
		}
		c.Next()
	}
}
