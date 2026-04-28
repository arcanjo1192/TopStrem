package auth

import (
    "context"
    "fmt"
    "net/http"
    "strings"

    "github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, "Token não fornecido", http.StatusUnauthorized)
            return
        }
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        claims := &jwt.MapClaims{}
        token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
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