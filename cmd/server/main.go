package main

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "strings"
    "time"

    "github.com/gin-gonic/gin"

    "topstrem/internal/api"
    "topstrem/internal/auth"
    "topstrem/internal/cache"
    "topstrem/internal/handlers"
    "topstrem/internal/middleware"
    "topstrem/internal/storage"
)

func main() {

    // Inicializar Redis  
    redisAddr := os.Getenv("REDIS_ADDR")  
    if redisAddr == "" {  
        redisAddr = "localhost:6379"  
    }  
    redisCache := cache.NewRedisCache(redisAddr) 

    // Inicializar Bolt DB para login/favoritos
    boltPath := os.Getenv("BOLT_DB_PATH")
    if boltPath == "" {
        homeDir, err := os.UserHomeDir()
        if err == nil {
            boltDir := filepath.Join(homeDir, ".topstrem")
            if mkdirErr := os.MkdirAll(boltDir, 0o700); mkdirErr == nil {
                boltPath = filepath.Join(boltDir, "topstrem.db")
            }
        }
        if boltPath == "" {
            boltPath = filepath.Join(".", "topstrem.db")
        }
    }
    boltStore, err := storage.Open(boltPath)
    if err != nil {
        panic(fmt.Sprintf("Failed to open Bolt DB: %v", err))
    }
    defer boltStore.Close()
    auth.SetStorage(boltStore)

    // ========== 1. INICIALIZAÇÃO DOS CLIENTES ==========
    apiClient := api.NewClient()
    watchClient := api.NewWatchHubClient()
    tmdbClient, err := api.NewTMDBClient()
    if err != nil {
        panic(fmt.Sprintf("Failed to initialize TMDB client: %v", err))
    }
    
    // Clientes com cache  
    cachedApiClient := api.NewCachedCinemetaClient(apiClient, redisCache)  
    cachedWatchClient := api.NewCachedWatchHubClient(watchClient, redisCache)  
    cachedTmdbClient := api.NewCachedTMDBClient(tmdbClient, redisCache)

    // ========== 2. DETECÇÃO DO DIRETÓRIO DE ASSETS ==========
    cwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    assetsDir := ""
    if _, err := os.Stat(filepath.Join(cwd, "cmd/app/assets")); err == nil {
        assetsDir = filepath.Join(cwd, "cmd/app/assets")
    } else if _, err := os.Stat(filepath.Join(cwd, "assets")); err == nil {
        assetsDir = filepath.Join(cwd, "assets")
    } else {
        panic("Assets directory not found")
    }

    // ========== 3. RATE LIMITERS ==========
    rateLimiter := middleware.NewRateLimiter(100, time.Minute)   // rotas gerais
    apiRateLimiter := middleware.NewRateLimiter(60, time.Minute) // rotas de API

    // ========== 4. INICIALIZAÇÃO DO GIN ==========
    gin.SetMode(gin.ReleaseMode)
    r := gin.New()
    if err := r.SetTrustedProxies([]string{"127.0.0.1"}); err != nil {
        panic(fmt.Sprintf("Failed to set trusted proxies: %v", err))
    }

    // Middleware global
    r.Use(gin.Logger(), gin.Recovery())

    // ========== 5. ARQUIVOS ESTÁTICOS (COM GZIP E CACHE) ==========
    staticDir := filepath.Join(assetsDir, "static")

    // Handler para arquivos estáticos com cache
    r.Static("/static", staticDir)
    r.Use(func(c *gin.Context) {
        if strings.HasPrefix(c.Request.URL.Path, "/static/") {
            c.Header("Cache-Control", "public, max-age=31536000") // 1 ano
        }
        c.Next()
    })

    r.GET("/robots.txt", func(c *gin.Context) {
        c.File(filepath.Join(assetsDir, "robots.txt"))
    })
    r.StaticFS("/privacy", http.Dir(filepath.Join(assetsDir, "privacy")))
    r.GET("/callback.html", func(c *gin.Context) {
        c.File(filepath.Join(assetsDir, "callback.html"))
    })

    // Sitemap
    r.GET("/sitemap.xml", middleware.CORS(), rateLimiter.Middleware(), handlers.SitemapHandler(cachedApiClient))

    // ========== 6. ROTAS DINÂMICAS COM MIDDLEWARES ==========
    // A ordem de aplicação: primeiro rate limit, depois CSRF (se necessário), depois CORS

    // Rota inicial
    r.GET("/", middleware.CORS(), rateLimiter.Middleware(), handlers.HomeHandler)

    // Catálogo e detalhes
    r.GET("/catalog/*path", middleware.CORS(), rateLimiter.Middleware(), handlers.CatalogHandler(cachedApiClient))
    r.GET("/detail/*path", middleware.CORS(), rateLimiter.Middleware(), handlers.DetailHandler(cachedApiClient, cachedTmdbClient))

    // Favoritos (precisa de CSRF para POST/DELETE)
    r.GET("/favorites", middleware.CORS(), rateLimiter.Middleware(), handlers.FavoritesHandler(cachedApiClient))

    // Endpoints de API com rate limit mais restrito
    r.GET("/api/episodes/*path", middleware.CORS(), apiRateLimiter.Middleware(), handlers.EpisodesHandler(cachedApiClient, cachedTmdbClient))
    r.GET("/api/watch/*path", middleware.CORS(), apiRateLimiter.Middleware(), handlers.WatchHandler(cachedWatchClient))
    r.GET("/api/favorites", middleware.CORS(), apiRateLimiter.Middleware(), handlers.FavoritesAPIHandler(boltStore))
    r.POST("/api/favorites", middleware.CORS(), apiRateLimiter.Middleware(), handlers.UpdateFavoritesAPIHandler(boltStore))
    r.GET("/api/watched-episodes", middleware.CORS(), apiRateLimiter.Middleware(), handlers.WatchedEpisodesAPIHandler(boltStore))
    r.POST("/api/watched-episodes", middleware.CORS(), apiRateLimiter.Middleware(), handlers.UpdateWatchedEpisodesAPIHandler(boltStore))

    // Autenticação
    r.GET("/auth/login", middleware.CORS(), rateLimiter.Middleware(), auth.LoginHandler)
    r.GET("/auth/callback", middleware.CORS(), rateLimiter.Middleware(), auth.CallbackHandler)
    
    r.GET("/api/me", middleware.CORS(), auth.MeHandler)
    r.POST("/auth/logout", middleware.CORS(), auth.LogoutHandler)

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{
            "status":   "ok",
            "timestamp": time.Now().Format(time.RFC3339),
        })
    })

    // ========== 7. CONFIGURAÇÃO DO AUTH (GOOGLE OAUTH) ==========
    baseURL := os.Getenv("BASE_URL")
    if baseURL == "" {
        baseURL = "http://localhost:8080"
    }
    redirectURL := baseURL + "/auth/callback"

    auth.InitAuth(
        os.Getenv("GOOGLE_CLIENT_ID"),
        os.Getenv("GOOGLE_CLIENT_SECRET"),
        redirectURL,
    )

    // ========== 8. INICIALIZAÇÃO DO SERVIDOR ==========
    fmt.Println("Servidor rodando na porta 8080...")
    r.Run(":8080")
}