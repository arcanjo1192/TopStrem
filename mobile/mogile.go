// mobile/mobile.go
package mobile

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"

    "github.com/gin-gonic/gin"

    "topstrem/internal/api"
    "topstrem/internal/auth"
    "topstrem/internal/crypto"
    "topstrem/internal/handlers"
    "topstrem/internal/middleware"
    "topstrem/internal/storage"
)

// Server é a struct exportada que será gerada como classe Java/Kotlin pelo gomobile.
type Server struct{}

// Start inicia o servidor HTTP local. Ele recebe:
//   - dataDir: diretório privado do app (ex.: /data/data/br.com.topstrem/files)
//   - assetsDir: diretório onde os assets do frontend foram extraídos
//   - secretsJSON: string JSON com os segredos (chave‑valor)
func (s *Server) Start(dataDir, assetsDir, secretsJSON string) {
    startServer(dataDir, assetsDir, secretsJSON)
}

// startServer contém toda a lógica de inicialização.
func startServer(dataDir, assetsDir, secretsJSON string) {
    // ---- Configuração de variáveis de ambiente a partir de segredos fornecidos pelo Android ----
    var secrets map[string]string
    if err := json.Unmarshal([]byte(secretsJSON), &secrets); err != nil {
        fmt.Println("Aviso: erro ao interpretar segredos:", err)
    } else {
        for key, value := range secrets {
            os.Setenv(key, value)
        }
    }

    // Valores fixos para o mobile
    os.Setenv("BASE_URL", "http://localhost:8080")
    os.Setenv("PORT", "8080")

    // ---- Banco de dados BoltDB ----
    boltPath := dataDir + "/topstrem.db"
    boltStore, err := storage.Open(boltPath)
    if err != nil {
        panic(fmt.Sprintf("Failed to open Bolt DB: %v", err))
    }
    defer boltStore.Close()
    auth.SetStorage(boltStore)

    // ---- Inicialização dos clientes (sem cache) ----
    apiClient := api.NewClient()
    watchClient := api.NewWatchHubClient()
    tmdbClient, err := api.NewTMDBClient()
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize TMDB client: %v", err))
    }

    // ---- Rate limiters ----
    rateLimiter := middleware.NewRateLimiter(30, time.Minute)
    apiRateLimiter := middleware.NewRateLimiter(20, time.Minute)

    // ---- Configuração do Gin ----
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    r.Use(gin.Recovery())

    // ---- Servir arquivos estáticos ----
    staticDir := assetsDir + "/static"
    r.Static("/static", staticDir)

    r.Use(func(c *gin.Context) {
        if strings.HasPrefix(c.Request.URL.Path, "/static/") {
            c.Header("Cache-Control", "public, max-age=31536000")
        }
        c.Next()
    })

    r.GET("/robots.txt", func(c *gin.Context) {
        c.File(assetsDir + "/robots.txt")
    })
    r.StaticFS("/privacy", http.Dir(assetsDir+"/privacy"))
    r.GET("/callback.html", func(c *gin.Context) {
        c.File(assetsDir + "/callback.html")
    })

    // ---- Rotas dinâmicas ----
    r.GET("/", middleware.CORS(), rateLimiter.Middleware(), handlers.HomeHandler)

    r.GET("/catalog/*path", middleware.CORS(), rateLimiter.Middleware(), handlers.CatalogHandler(apiClient))
    r.GET("/detail/*path", middleware.CORS(), rateLimiter.Middleware(), handlers.DetailHandler(apiClient, tmdbClient))
    r.GET("/favorites", middleware.CORS(), rateLimiter.Middleware(), handlers.FavoritesHandler(apiClient))
    r.GET("/lists", middleware.CORS(), rateLimiter.Middleware(), handlers.ListsHandler(apiClient, boltStore))

    r.GET("/api/episodes/*path", middleware.CORS(), apiRateLimiter.Middleware(), handlers.EpisodesHandler(apiClient, tmdbClient))
    r.GET("/api/watch/*path", middleware.CORS(), apiRateLimiter.Middleware(), handlers.WatchHandler(watchClient, tmdbClient))
    r.GET("/api/favorites", middleware.CORS(), apiRateLimiter.Middleware(), handlers.FavoritesAPIHandler(boltStore))
    r.POST("/api/favorites", middleware.CORS(), apiRateLimiter.Middleware(), handlers.UpdateFavoritesAPIHandler(boltStore))
    r.GET("/api/lists", middleware.CORS(), apiRateLimiter.Middleware(), handlers.ListsAPIHandler(boltStore))
    r.POST("/api/lists", middleware.CORS(), apiRateLimiter.Middleware(), handlers.UpdateListsAPIHandler(boltStore))
    r.GET("/api/search", middleware.CORS(), apiRateLimiter.Middleware(), handlers.SearchHandler(apiClient))
    r.GET("/api/watched-episodes", middleware.CORS(), apiRateLimiter.Middleware(), handlers.WatchedEpisodesAPIHandler(boltStore))
    r.POST("/api/watched-episodes", middleware.CORS(), apiRateLimiter.Middleware(), handlers.UpdateWatchedEpisodesAPIHandler(boltStore))
    r.GET("/api/share-token", middleware.CORS(), handlers.ShareTokenHandler())

    // Autenticação
    r.GET("/auth/login", middleware.CORS(), rateLimiter.Middleware(), auth.LoginHandler)
    r.GET("/auth/callback", middleware.CORS(), rateLimiter.Middleware(), auth.CallbackHandler)
    r.GET("/api/me", middleware.CORS(), auth.MeHandler)
    r.POST("/auth/logout", middleware.CORS(), auth.LogoutHandler)

    // Health check
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"status": "ok"})
    })

    // ---- Configuração de autenticação Google (somente se as chaves foram fornecidas) ----
    googleID := os.Getenv("GOOGLE_CLIENT_ID")
    googleSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    if googleID != "" && googleSecret != "" {
        auth.InitAuth(
            googleID,
            googleSecret,
            "http://localhost:8080/auth/callback",
        )
    } else {
        fmt.Println("Aviso: autenticação Google não configurada")
    }

    // ---- Criptografia para listas compartilháveis ----
    if key := os.Getenv("LIST_SHARE_SECRET"); key != "" {
        if err := crypto.Init(key); err != nil {
            fmt.Println("Aviso: criptografia de listas não inicializada:", err)
        }
    }

    // ---- Inicia o servidor ----
    fmt.Println("Servidor mobile rodando em http://localhost:8080")
    r.Run("0.0.0.0:8080")
}