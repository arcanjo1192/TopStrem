package middleware

import (
    "crypto/rand"
    "encoding/base64"
    "fmt"
    "net/http"
    "time"
)

func CSRF(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Para GET, HEAD, OPTIONS - gerar novo token
        if r.Method == "GET" || r.Method == "HEAD" || r.Method == "OPTIONS" {
            token := generateCSRFToken()
            cookie := &http.Cookie{
                Name:     "csrf_token",
                Value:    token,
                Path:     "/",
                HttpOnly: true,
                Secure:   true, // HTTPS em produção
                SameSite: http.SameSiteStrictMode,
                MaxAge:   3600, // 1 hora
            }
            http.SetCookie(w, cookie)
            next(w, r)
            return
        }

        // Para POST, PUT, DELETE - validar token
        cookieToken := getCSRFTokenFromCookie(r)
        headerToken := r.Header.Get("X-CSRF-Token")
        formToken := r.FormValue("_csrf_token")

        // Aceitar token do header ou form
        token := headerToken
        if token == "" {
            token = formToken
        }

        // Validar que temos ambos os tokens e que são iguais
        if cookieToken == "" || token == "" || cookieToken != token {
            http.Error(w, fmt.Sprintf("CSRF token invalid or missing"), http.StatusForbidden)
            return
        }

        next(w, r)
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