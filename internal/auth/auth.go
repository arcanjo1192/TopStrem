package auth

import (
    "context"
    "errors"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/oauth2"
    "golang.org/x/oauth2/google"
    "topstrem/internal/storage"
)

// ClientType define o tipo de cliente que se autentica
type ClientType string

const (
    ClientTypeWebBrowser ClientType = "web_browser"      // Desktop ou Mobile Browser
    ClientTypeNativeApp  ClientType = "native_app"       // App nativo iOS/Android
    ClientTypeCustomTab  ClientType = "custom_tab"       // Custom Tab Android (WebView)
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
var authStore *storage.Storage

func SetStorage(store *storage.Storage) {
    authStore = store
}

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

// detectClientType identifica o tipo de cliente baseado no User-Agent e headers
func detectClientType(r *http.Request) ClientType {
    userAgent := r.Header.Get("User-Agent")
    
    // Header custom enviado por app nativo
    if clientType := r.Header.Get("X-Client-Type"); clientType != "" {
        if clientType == "native" {
            return ClientTypeNativeApp
        }
    }
    
    // Detectar WebView Android (Custom Tab)
    if strings.Contains(userAgent, "wv") || 
       strings.Contains(userAgent, "WebView") {
        return ClientTypeCustomTab
    }
    
    // Padrão: Web Browser (desktop ou mobile)
    return ClientTypeWebBrowser
}

func LoginHandler(c *gin.Context) {
    url := authConfig.OAuth2Config.AuthCodeURL("random-state", oauth2.AccessTypeOffline)
    clientType := detectClientType(c.Request)

    switch clientType {
    case ClientTypeNativeApp:
        // App nativo: retorna JSON com URL de autenticação
        c.JSON(http.StatusOK, map[string]interface{}{
            "status":  "ok",
            "authUrl": url,
            "message": "Abra esta URL no seu navegador padrão para autenticar",
        })
        
    case ClientTypeCustomTab:
        // Custom Tab Android: retorna JSON com URL
        c.JSON(http.StatusOK, map[string]interface{}{
            "status":  "ok",
            "authUrl": url,
        })
        
    default: // ClientTypeWebBrowser
        // Web Browser: redireciona direto para Google
        c.Redirect(http.StatusTemporaryRedirect, url)
    }
}

func CallbackHandler(c *gin.Context) {
    code := c.Query("code")
    if code == "" {
        c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Código de autorização não encontrado"})
        return
    }
    
    token, err := authConfig.OAuth2Config.Exchange(context.Background(), code)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Falha ao trocar o token"})
        return
    }
    
    idToken, ok := token.Extra("id_token").(string)
    if !ok {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "ID Token não encontrado"})
        return
    }
    
    userInfo, err := validateGoogleIDToken(idToken)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "ID Token inválido"})
        return
    }
    
    jwtToken, err := generateJWT(userInfo.Email, userInfo.Name)
    if err != nil {
        c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao gerar token de sessão"})
        return
    }

    if authStore != nil {
        authStore.SaveUser(storage.UserProfile{
            Email:     userInfo.Email,
            Name:      userInfo.Name,
            LastLogin: time.Now().UTC(),
        })
    }
    
    // Armazenar token em HttpOnly secure cookie (funciona para web e custom tab)
    c.SetCookie("auth_token", jwtToken, 72*3600, "/", "", authConfig.CookieSecure, true)
    
    // Detectar se veio de app nativo (via parâmetro na URL)
    clientType := c.Query("client_type")
    
    if clientType == "native" {
        // App nativo: retornar HTML que passa token para app nativo via scheme customizado
        html := `<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>Autenticação</title>
    <style>
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif; display: flex; align-items: center; justify-content: center; height: 100vh; margin: 0; background: #f5f5f5; }
        .container { text-align: center; }
        h1 { color: #333; margin-bottom: 10px; }
        p { color: #666; }
    </style>
</head>
<body>
    <div class="container">
        <h1>✅ Autenticação bem-sucedida!</h1>
        <p>Você será redirecionado para o aplicativo...</p>
    </div>
    <script>
        const token = '` + jwtToken + `';
        
        // Tentar abrir scheme customizado do app nativo com o token
        const deepLink = 'topstrem://auth?token=' + encodeURIComponent(token);
        
        // Tentar abrir
        setTimeout(() => {
            window.location.href = deepLink;
        }, 500);
        
        // Se não abrir, exibir token para copiar manualmente
        setTimeout(() => {
            document.querySelector('.container').innerHTML += '<div style="margin-top: 20px; padding: 20px; background: white; border-radius: 8px; max-width: 500px;"><p>Se o aplicativo não abrir, copie este token:</p><input type="text" value="' + token + '" readonly style="width: 100%; padding: 10px; margin-top: 10px; border: 1px solid #ddd; border-radius: 4px;"></div>';
        }, 2000);
    </script>
</body>
</html>`
        c.Header("Content-Type", "text/html; charset=utf-8")
        c.String(http.StatusOK, html)
        return
    }
    
    // Web Browser: redirecionar para home (cookie já foi setado)
    c.Redirect(http.StatusFound, "/")
}

func MeHandler(c *gin.Context) {
    tokenString := getTokenFromRequest(c.Request)
    if tokenString == "" {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
        return
    }

    // Validar e parsear JWT
    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return authConfig.JWTSecret, nil
    })

    if err != nil || !token.Valid {
        c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token inválido"})
        return
    }

    // Retornar dados do usuário
    c.JSON(http.StatusOK, map[string]interface{}{
        "email": claims["email"],
        "name":  claims["name"],
    })
}

// LogoutHandler limpa o cookie de autenticação (para web) ou retorna sucesso (para app nativo)
func LogoutHandler(c *gin.Context) {
    if c.Request.Method != "POST" {
        c.AbortWithStatusJSON(http.StatusMethodNotAllowed, gin.H{"error": "Método não permitido"})
        return
    }

    // Limpar cookie definindo MaxAge negativo (funciona para web)
    c.SetCookie("auth_token", "", -1, "/", "", authConfig.CookieSecure, true)

    c.JSON(http.StatusOK, gin.H{"status": "logged out"})
}

// getTokenFromRequest extrai o JWT do cookie ou do header Authorization
func getTokenFromRequest(r *http.Request) string {
    // Prioridade 1: Cookie HttpOnly (web browser)
    if cookie, err := r.Cookie("auth_token"); err == nil && cookie.Value != "" {
        return cookie.Value
    }
    
    // Prioridade 2: Header Authorization (app nativo)
    // Formato: "Bearer <token>"
    authHeader := r.Header.Get("Authorization")
    if authHeader != "" {
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) == 2 && parts[0] == "Bearer" {
            return parts[1]
        }
    }
    
    return ""
}

func GetEmailFromRequest(r *http.Request) (string, error) {
    tokenString := getTokenFromRequest(r)
    if tokenString == "" {
        return "", errors.New("token não encontrado")
    }

    claims := jwt.MapClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return authConfig.JWTSecret, nil
    })
    if err != nil || !token.Valid {
        return "", errors.New("token inválido")
    }

    email, ok := claims["email"].(string)
    if !ok || email == "" {
        return "", errors.New("email não encontrado no token")
    }
    return email, nil
}