package auth

import (
    "context"
    "encoding/json"
    "net/http"
    "os"
    "strings"

    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
)

type AuthConfig struct {
    GoogleClientID     string
    GoogleClientSecret string
    CookieSecure       bool
    JWTSecret          []byte
    RedirectURL        string
    OAuth2Config       *oauth2.Config
}

var authConfig AuthConfig

func InitAuth(clientID, clientSecret, redirectURL string) {
    // Validar JWT_SECRET
    jwtSecret := os.Getenv("JWT_SECRET")
    if jwtSecret == "" || len(jwtSecret) < 32 {
        panic("ERRO: JWT_SECRET deve estar configurado e ter no mínimo 32 caracteres")
    }
    
    // Validar credenciais Google
    if clientID == "" || clientSecret == "" {
        panic("ERRO: GOOGLE_CLIENT_ID e GOOGLE_CLIENT_SECRET devem estar configurados")
    }
    
    authConfig = AuthConfig{
        GoogleClientID:     clientID,
        GoogleClientSecret: clientSecret,
        RedirectURL:        redirectURL,
        CookieSecure:       os.Getenv("ENVIRONMENT") == "production",
        JWTSecret:          []byte(jwtSecret),
        OAuth2Config: &oauth2.Config{
            ClientID:     clientID,
            ClientSecret: clientSecret,
            RedirectURL:  redirectURL,
            Scopes:       []string{"openid", "profile", "email"},
            Endpoint:     google.Endpoint,
        },
    }
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
    url := authConfig.OAuth2Config.AuthCodeURL("random-state", oauth2.AccessTypeOffline)

    userAgent := r.Header.Get("User-Agent")
    // Detectar webview Android: "wv", "WebView", "Android"
    isWebView := strings.Contains(userAgent, "wv") || 
                 strings.Contains(userAgent, "WebView") || 
                 strings.Contains(userAgent, "Android")

    // Se for webview, retornar JSON com a URL de login
    if isWebView {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]string{
            "url": url,
            "success": "true",
        })
        return
    }

    // Caso contrário, redirecionar diretamente
    http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {
    code := r.URL.Query().Get("code")
    if code == "" {
        http.Error(w, "Código de autorização não encontrado", http.StatusBadRequest)
        return
    }
    
    token, err := authConfig.OAuth2Config.Exchange(context.Background(), code)
    if err != nil {
        http.Error(w, "Falha ao trocar o token", http.StatusInternalServerError)
        return
    }
    
    idToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "ID Token não encontrado", http.StatusInternalServerError)
        return
    }
    
    userInfo, err := validateGoogleIDToken(idToken)
    if err != nil {
        http.Error(w, "ID Token inválido", http.StatusUnauthorized)
        return
    }
    
    jwtToken, err := generateJWT(userInfo.Email, userInfo.Name)
    if err != nil {
        http.Error(w, "Erro ao gerar token de sessão", http.StatusInternalServerError)
        return
    }
    
    // Armazenar token em HttpOnly secure cookie em vez de URL
    cookie := &http.Cookie{
        Name:     "auth_token",
        Value:    jwtToken,
        Path:     "/",
        HttpOnly: true,
        Secure:   authConfig.CookieSecure, // true em produção
        SameSite: http.SameSiteStrictMode,
        MaxAge:   72 * 3600, // 72 horas
    }
    http.SetCookie(w, cookie)
    
    // Redirecionar para home sem expor o token
    http.Redirect(w, r, "/", http.StatusFound)
}

// MeHandler retorna dados do usuário autenticado a partir do JWT no cookie
func MeHandler(w http.ResponseWriter, r *http.Request) {
    // Obter token do cookie
    cookie, err := r.Cookie("auth_token")
    if err != nil {
        http.Error(w, "Não autenticado", http.StatusUnauthorized)
        return
    }

    tokenString := cookie.Value
    if tokenString == "" {
        http.Error(w, "Não autenticado", http.StatusUnauthorized)
        return
    }

    // Validar e parsear JWT
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return authConfig.JWTSecret, nil
    })

    if err != nil || !token.Valid {
        http.Error(w, "Token inválido", http.StatusUnauthorized)
        return
    }

    // Retornar dados do usuário
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "email": claims["email"],
        "name":  claims["name"],
    })
}

// LogoutHandler limpa o cookie de autenticação
func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != "POST" {
        http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
        return
    }

    // Limpar cookie definindo MaxAge negativo
    cookie := &http.Cookie{
        Name:     "auth_token",
        Value:    "",
        Path:     "/",
        HttpOnly: true,
        Secure:   authConfig.CookieSecure,
        SameSite: http.SameSiteStrictMode,
        MaxAge:   -1,
    }
    http.SetCookie(w, cookie)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{"status": "logged out"})
}