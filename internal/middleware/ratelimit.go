package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimiter struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     int
	burst    int
	window   time.Duration
}

type visitor struct {
	tokens   int
	lastSeen time.Time
}

func newRateLimiter(rate, burst int, window time.Duration) *rateLimiter {
	rl := &rateLimiter{
		visitors: make(map[string]*visitor),
		rate:     rate,
		burst:    burst,
		window:   window,
	}

	go rl.cleanupVisitors()
	return rl
}

func (rl *rateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > rl.window {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *rateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{tokens: rl.burst - 1, lastSeen: time.Now()}
		return true
	}

	elapsed := time.Since(v.lastSeen)
	v.lastSeen = time.Now()

	tokensToAdd := int(elapsed.Seconds()) * rl.rate
	v.tokens = min(v.tokens+tokensToAdd, rl.burst)

	if v.tokens > 0 {
		v.tokens--
		return true
	}

	return false
}

func RateLimitMiddleware(requestsPerSecond, burst int) gin.HandlerFunc {
	limiter := newRateLimiter(requestsPerSecond, burst, time.Minute*5)

	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
			c.Abort()
			return
		}
		c.Next()
	}
}
