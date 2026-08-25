package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type rateLimitEntry struct {
	count     int
	resetTime time.Time
}

var (
	rateLimitStore = sync.Map{}
	rateLimitMax   = 10
	rateLimitWindow = 15 * time.Minute
)

// RateLimiter limits requests per IP within a time window.
// Default: 10 requests per 15 minutes per IP.
func RateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()

		now := time.Now()
		val, _ := rateLimitStore.LoadOrStore(ip, &rateLimitEntry{
			count:     0,
			resetTime: now.Add(rateLimitWindow),
		})

		entry := val.(*rateLimitEntry)

		// Reset window if expired
		if now.After(entry.resetTime) {
			entry.count = 0
			entry.resetTime = now.Add(rateLimitWindow)
		}

		entry.count++

		if entry.count > rateLimitMax {
			remaining := entry.resetTime.Sub(now).Seconds()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error":       "too many requests, please try again later",
				"retry_after": int(remaining),
			})
			c.Abort()
			return
		}

		rateLimitStore.Store(ip, entry)
		c.Next()
	}
}
