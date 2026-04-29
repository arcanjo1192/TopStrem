package main

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"
    "time"

    "github.com/NYTimes/gziphandler" // ← import adicionado

    "topstrem/internal/cache"
    "topstrem/internal/api"
    "topstrem/internal/auth"
    "topstrem/internal/handlers"
    "topstrem/internal/middleware"
)

func main() {

    // Inicializar Redis  
    redisAddr := os.Getenv("REDIS_ADDR")  
    if redisAddr == "" {  
        redisAddr = "localhost:6379"  
    }  
    redisCache := cache.NewRedisCache(redisAddr) 

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
    apiRateLimiter := middleware.NewRateLimiter(30, time.Minute) // rotas de API

    // ========== 4. ARQUIVOS ESTÁTICOS (COM GZIP E CACHE) ==========
    staticDir := filepath.Join(assetsDir, "static")

    // Handler base para arquivos estáticos (sem middlewares extras)
    staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))

    // Aplica Gzip + Cache-Control via gziphandler e HandlerFunc personalizado
    http.Handle("/static/", gziphandler.GzipHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Cache-Control", "public, max-age=31536000") // 1 ano
        staticHandler.ServeHTTP(w, r)
    })))

    http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filepath.Join(assetsDir, "robots.txt"))
    })
    http.Handle("/privacy/", http.StripPrefix("/privacy/", http.FileServer(http.Dir(filepath.Join(assetsDir, "privacy")))))
    http.HandleFunc("/callback.html", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filepath.Join(assetsDir, "callback.html"))
    })

    // ========== 5. ROTAS DINÂMICAS COM MIDDLEWARES ==========
    // A ordem de aplicação: primeiro rate limit, depois CSRF (se necessário), depois CORS

    // Rota inicial
    http.HandleFunc("/", middleware.CORS(rateLimiter.Middleware(handlers.HomeHandler)))

    // Catálogo e detalhes
    http.HandleFunc("/catalog/", middleware.CORS(rateLimiter.Middleware(handlers.CatalogHandler(cachedApiClient))))
    http.HandleFunc("/detail/", middleware.CORS(rateLimiter.Middleware(handlers.DetailHandler(cachedApiClient, cachedTmdbClient))))

    // Favoritos (precisa de CSRF para POST/DELETE)
    http.HandleFunc("/favorites", middleware.CORS(middleware.CSRF(rateLimiter.Middleware(handlers.FavoritesHandler(cachedApiClient)))))

    // Endpoints de API com rate limit mais restrito
    http.HandleFunc("/api/episodes/", middleware.CORS(apiRateLimiter.Middleware(handlers.EpisodesHandler(cachedApiClient, cachedTmdbClient))))
    http.HandleFunc("/api/watch/", middleware.CORS(apiRateLimiter.Middleware(handlers.WatchHandler(cachedWatchClient))))

    // Autenticação
    http.HandleFunc("/auth/login", middleware.CORS(rateLimiter.Middleware(auth.LoginHandler)))
    http.HandleFunc("/auth/callback", middleware.CORS(rateLimiter.Middleware(auth.CallbackHandler)))
    
    http.HandleFunc("/api/me", middleware.CORS(auth.MeHandler))
    http.HandleFunc("/auth/logout", middleware.CORS(auth.LogoutHandler))

    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
    })

    // ========== 6. CONFIGURAÇÃO DO AUTH (GOOGLE OAUTH) ==========
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

    // ========== 7. INICIALIZAÇÃO DO SERVIDOR ==========
    fmt.Println("Servidor rodando na porta 8080...")
    http.ListenAndServe(":8080", nil)
}