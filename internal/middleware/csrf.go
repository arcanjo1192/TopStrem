package middleware

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
)

func CSRF() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Para GET, HEAD, OPTIONS - gerar novo token
        if c.Request.Method == "GET" || c.Request.Method == "HEAD" || c.Request.Method == "OPTIONS" {
            token := generateCSRFToken()
            c.SetCookie("csrf_token", token, 3600, "/", "", true, true)
            c.Next()
            return
        }

        // Para POST, PUT, DELETE - validar token
        cookieToken, _ := c.Cookie("csrf_token")
        headerToken := c.GetHeader("X-CSRF-Token")
        formToken := c.PostForm("csrf_token")

        // Aceitar token do header ou form
        token := headerToken
        if token == "" {
            token = formToken
        }

        // Validar que temos ambos os tokens e que são iguais
        if cookieToken == "" || token == "" || cookieToken != token {
            c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": fmt.Sprintf("CSRF token invalid or missing")})
            return
        }

        c.Next()
    }
}

func generateCSRFToken() string {
    b := make([]byte, 32)
    if _, err := rand.Read(b); err != nil {
        // Fallback se rand.Read falhar
        return base64.StdEncoding.EncodeToString([]byte(time.Now().String()))
    }
    return base64.StdEncoding.EncodeToString(b)
}

func getCSRFTokenFromCookie(r *http.Request) string {
    cookie, err := r.Cookie("csrf_token")
    if err != nil {
        return ""
    }
    return cookie.Value
}