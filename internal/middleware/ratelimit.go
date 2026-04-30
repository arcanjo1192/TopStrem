package middleware  
  
import (  
    "net/http"  
    "sync"  
    "time"  

    "github.com/gin-gonic/gin"
)  
  
type RateLimiter struct {  
    clients map[string][]time.Time  
    mutex   sync.Mutex  
    limit   int  
    window  time.Duration  
}  
  
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {  
    return &RateLimiter{  
        clients: make(map[string][]time.Time),  
        limit:   limit,  
        window:  window,  
    }  
}  
  
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        clientIP := c.ClientIP()

        rl.mutex.Lock()
        defer rl.mutex.Unlock()

        now := time.Now()
        requests := rl.clients[clientIP]

        // Remove requests antigos
        validRequests := make([]time.Time, 0)
        for _, reqTime := range requests {
            if now.Sub(reqTime) < rl.window {
                validRequests = append(validRequests, reqTime)
            }
        }

        if len(validRequests) >= rl.limit {
            c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{"error": "Rate limit exceeded"})
            return
        }

        validRequests = append(validRequests, now)

        // Armazenar nova lista ou deletar se vazia (evita memory leak)
        if len(validRequests) > 0 {
            rl.clients[clientIP] = validRequests
        } else {
            delete(rl.clients, clientIP)
        }

        c.Next()
    }
}