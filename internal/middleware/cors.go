package middleware

import (
    "net/http"
    "os"
    "strings"

    "github.com/gin-gonic/gin"
)

// allowedOrigins contém as origens CORS permitidas
var allowedOrigins map[string]bool

func init() {
    allowedOrigins = make(map[string]bool)
    
    // Origem de desenvolvimento
    if devOrigin := os.Getenv("DEV_ORIGIN"); devOrigin != "" {
        allowedOrigins[devOrigin] = true
    }
    
    // Origens de produção (separadas por vírgula)
    if prodOrigins := os.Getenv("ALLOWED_ORIGINS"); prodOrigins != "" {
        for _, origin := range strings.Split(prodOrigins, ",") {
            origin = strings.TrimSpace(origin)
            if origin != "" {
                allowedOrigins[origin] = true
            }
        }
    }
    
    // Se não houver nenhuma origem configurada, aceitar localhost por padrão
    if len(allowedOrigins) == 0 {
        allowedOrigins["http://localhost:8080"] = true
        allowedOrigins["http://localhost:3000"] = true
    }
}

func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        origin := c.GetHeader("Origin")
        
        // Verificar se a origem é permitida
        if origin != "" && allowedOrigins[origin] {
            c.Header("Access-Control-Allow-Origin", origin)
            c.Header("Access-Control-Allow-Credentials", "true")
            c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
            c.Header("Access-Control-Max-Age", "86400") // 24 horas
        }
        
        // Responder a preflight requests
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(http.StatusOK)
            return
        }
        
        c.Next()
    }
}