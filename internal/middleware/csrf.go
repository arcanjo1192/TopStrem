package middleware  
  
import (  
    "crypto/rand"  
    "encoding/base64"  
    "net/http"  
)  
  
func CSRF(next http.HandlerFunc) http.HandlerFunc {  
    return func(w http.ResponseWriter, r *http.Request) {  
        if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {  
            next(w, r)  
            return  
        }  
          
        token := r.Header.Get("X-CSRF-Token")  
        sessionToken := getSessionCSRFToken(r)  
          
        if token == "" || sessionToken == "" || token != sessionToken {  
            http.Error(w, "CSRF token invalid", http.StatusForbidden)  
            return  
        }  
          
        next(w, r)  
    }  
}  
  
func generateCSRFToken() string {  
    b := make([]byte, 32)  
    rand.Read(b)  
    return base64.StdEncoding.EncodeToString(b)  
}  
  
func getSessionCSRFToken(r *http.Request) string {  
    // Implementar lógica para obter token da sessão/cookie  
    return ""  
}