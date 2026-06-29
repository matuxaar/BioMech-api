package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	rateLimiters   []func()
	rateLimitersMu sync.Mutex
)

func StopRateLimiters() {
	rateLimitersMu.Lock()
	stops := make([]func(), len(rateLimiters))
	copy(stops, rateLimiters)
	rateLimiters = nil
	rateLimitersMu.Unlock()
	for _, stop := range stops {
		stop()
	}
}

type rateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*clientBucket
	rate     int
	burst    int
	cleanup  time.Duration
	stopCh   chan struct{}
}

type clientBucket struct {
	tokens   int
	last     time.Time
}

func RateLimit(rate, burst int) gin.HandlerFunc {
	rl := &rateLimiter{
		clients: make(map[string]*clientBucket),
		rate:    rate,
		burst:   burst,
		cleanup: time.Minute,
		stopCh:  make(chan struct{}),
	}

	go func() {
		for {
			select {
			case <-rl.stopCh:
				return
			case <-time.After(rl.cleanup):
			}
			rl.mu.Lock()
			for ip, b := range rl.clients {
				if time.Since(b.last) > rl.cleanup*2 {
					delete(rl.clients, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()

	rateLimitersMu.Lock()
	rateLimiters = append(rateLimiters, func() { close(rl.stopCh) })
	rateLimitersMu.Unlock()

	return func(c *gin.Context) {
		ip := c.ClientIP()
		rl.mu.Lock()
		b, ok := rl.clients[ip]
		now := time.Now()
		if !ok {
			b = &clientBucket{tokens: rl.burst, last: now}
			rl.clients[ip] = b
		}
		elapsed := now.Sub(b.last)
		b.last = now
		b.tokens += int(elapsed.Seconds() * float64(rl.rate))
		if b.tokens > rl.burst {
			b.tokens = rl.burst
		}
		if b.tokens > 0 {
			b.tokens--
			rl.mu.Unlock()
			c.Next()
		} else {
			rl.mu.Unlock()
			rateLimitedRequests.Inc()
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
		}
	}
}
