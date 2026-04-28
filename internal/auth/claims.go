package auth

import (
    "context"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
    "google.golang.org/api/idtoken"
)

type User struct {
    Email string
    Name  string
}

func validateGoogleIDToken(idToken string) (*User, error) {
    if authConfig.GoogleClientID == "" {
        return nil, fmt.Errorf("GoogleClientID não configurado")
    }
    payload, err := idtoken.Validate(context.Background(), idToken, authConfig.GoogleClientID)
    if err != nil {
        return nil, fmt.Errorf("token inválido: %v", err)
    }
    email, _ := payload.Claims["email"].(string)
    name, _ := payload.Claims["name"].(string)
    if email == "" {
        return nil, fmt.Errorf("email não encontrado no token")
    }
    if name == "" {
        name = email
    }
    return &User{Email: email, Name: name}, nil
}

func generateJWT(email, name string) (string, error) {
    claims := jwt.MapClaims{
        "email": email,
        "name":  name,
        "exp":   time.Now().Add(72 * time.Hour).Unix(),
        "iat":   time.Now().Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(authConfig.JWTSecret)
}