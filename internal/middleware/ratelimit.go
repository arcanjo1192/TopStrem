package middleware  
  
import (  
    "net/http"  
    "sync"  
    "time"  
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
  
func (rl *RateLimiter) Middleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        clientIP := r.RemoteAddr

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
            http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
            return
        }

        validRequests = append(validRequests, now)

        // Armazenar nova lista ou deletar se vazia (evita memory leak)
        if len(validRequests) > 0 {
            rl.clients[clientIP] = validRequests
        } else {
            delete(rl.clients, clientIP)
        }

        next(w, r)
    }
}