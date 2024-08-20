package middlewares

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

var mu sync.Mutex
var requestCounts = make(map[string]int)

const maxRequests = 10

func resetRequestCounts(ip string) {
	for {
		loops := 0
		time.Sleep(time.Minute * 10)
		mu.Lock()
		defer mu.Unlock()
		if requestCounts[ip] == 0 {
			loops++
			if loops > 2 {
				delete(requestCounts, ip)
				return
			}
		} else {
			loops = 0
		}
		requestCounts[ip] = 0
	}
}

func RateLimitMiddleware(c *gin.Context) {
	ip := c.ClientIP()
	mu.Lock()
	if _, exists := requestCounts[ip]; !exists {
		requestCounts[ip] = 0
		go resetRequestCounts(ip)
	}
	if requestCounts[ip] >= maxRequests {
		c.JSON(http.StatusTooManyRequests, gin.H{"message": "Rate limit exceeded"})
		c.Abort()
		return
	}

	requestCounts[ip]++
	mu.Unlock()
	c.Next()
}
