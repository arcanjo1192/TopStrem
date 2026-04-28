package auth

import (
    "context"
    "fmt"
    "net/http"
    "os"

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
    authConfig = AuthConfig{
        GoogleClientID:     clientID,
        GoogleClientSecret: clientSecret,
        RedirectURL:        redirectURL,
        CookieSecure:       false,
        JWTSecret:          []byte(os.Getenv("JWT_SECRET")),
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
        http.Error(w, "Falha ao trocar o token: "+err.Error(), http.StatusInternalServerError)
        return
    }
    idToken, ok := token.Extra("id_token").(string)
    if !ok {
        http.Error(w, "ID Token não encontrado", http.StatusInternalServerError)
        return
    }
    userInfo, err := validateGoogleIDToken(idToken)
    if err != nil {
        http.Error(w, "ID Token inválido: "+err.Error(), http.StatusUnauthorized)
        return
    }
    jwtToken, err := generateJWT(userInfo.Email, userInfo.Name)
    if err != nil {
        http.Error(w, "Erro ao gerar token de sessão", http.StatusInternalServerError)
        return
    }
    callbackURL := fmt.Sprintf("/callback.html?token=%s", jwtToken)
    http.Redirect(w, r, callbackURL, http.StatusFound)
}