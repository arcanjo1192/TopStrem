package middleware

import (
    "net/http"
    "os"
    "strings"
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

func CORS(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        
        // Verificar se a origem é permitida
        if origin != "" && allowedOrigins[origin] {
            w.Header().Set("Access-Control-Allow-Origin", origin)
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-CSRF-Token")
            w.Header().Set("Access-Control-Max-Age", "86400") // 24 horas
        }
        
        // Responder a preflight requests
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }
        
        next(w, r)
    }
}