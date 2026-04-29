package auth

import (
    "context"
    "fmt"
    "net/http"

    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        tokenString := getTokenFromRequest(r)
        if tokenString == "" {
            http.Error(w, "Token não fornecido", http.StatusUnauthorized)
            return
        }
        
        claims := &jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
            // Rejeitar algoritmo "none"
            if token.Method.Alg() == "none" {
                return nil, fmt.Errorf("algoritmo 'none' não permitido")
            }
            
            // Validar que é HMAC
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("método inesperado: %v", token.Header["alg"])
            }
            
            return authConfig.JWTSecret, nil
        })
        
        if err != nil || !token.Valid {
            http.Error(w, "Token inválido", http.StatusUnauthorized)
            return
        }
        
        ctx := context.WithValue(r.Context(), "user_email", (*claims)["email"])
        ctx = context.WithValue(ctx, "user_name", (*claims)["name"])
        next(w, r.WithContext(ctx))
    }
}