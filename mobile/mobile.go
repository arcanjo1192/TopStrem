package mobile

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "path/filepath"
    "syscall"
    "time"

    "github.com/NYTimes/gziphandler"
    "github.com/joho/godotenv"

    "topstrem/internal/api"
    "topstrem/internal/auth"
    "topstrem/internal/cache"
    "topstrem/internal/handlers"
    "topstrem/internal/middleware"
)

func StartServer() {
    // ========== 0. CARREGAR VARIÁVEIS DE AMBIENTE ==========
    // Tentar carregar .env da raiz do projeto
    if err := godotenv.Load(); err != nil {
        fmt.Println("⚠️  .env não encontrado, usando variáveis de ambiente do sistema")
    } else {
        fmt.Println("✅ Variáveis de ambiente carregadas do .env")
    }

    // ========== 1. OBTER PORTA ==========
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    // ========== 2. INICIALIZAÇÃO DO CACHE REDIS ==========
    redisAddr := os.Getenv("REDIS_ADDR")
    if redisAddr == "" {
        redisAddr = "localhost:6379"
    }
    redisCache := cache.NewRedisCache(redisAddr)

    // ========== 3. INICIALIZAÇÃO DE AUTENTICAÇÃO OAUTH ==========
    googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
    googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
    redirectURL := os.Getenv("REDIRECT_URL")

    if googleClientID == "" || googleClientSecret == "" {
        fmt.Println("⚠️  AVISO: Variáveis Google OAuth não configuradas. Autenticação desabilitada.")
        fmt.Println("Configure GOOGLE_CLIENT_ID, GOOGLE_CLIENT_SECRET e REDIRECT_URL no .env")
    } else {
        auth.InitAuth(googleClientID, googleClientSecret, redirectURL)
    }

    // ========== 4. INICIALIZAÇÃO DOS CLIENTES DE API ==========
    apiClient := api.NewClient()
    watchClient := api.NewWatchHubClient()
    tmdbClient, err := api.NewTMDBClient()
    if err != nil {
        log.Printf("ERRO ao criar cliente TMDB: %v", err)
        log.Println("Continuando sem TMDB...")
    }

    // Clientes com cache
    cachedApiClient := api.NewCachedCinemetaClient(apiClient, redisCache)
    cachedWatchClient := api.NewCachedWatchHubClient(watchClient, redisCache)
    var cachedTmdbClient api.TMDBClient
    if tmdbClient != nil {
        cachedTmdbClient = api.NewCachedTMDBClient(tmdbClient, redisCache)
    }

    // ========== 5. DETECÇÃO DO DIRETÓRIO DE ASSETS ==========
    cwd, err := os.Getwd()
    if err != nil {
        log.Fatalf("ERRO ao obter diretório atual: %v", err)
    }

    assetsDir := ""
    if _, err := os.Stat(filepath.Join(cwd, "cmd/app/assets")); err == nil {
        assetsDir = filepath.Join(cwd, "cmd/app/assets")
    } else if _, err := os.Stat(filepath.Join(cwd, "assets")); err == nil {
        assetsDir = filepath.Join(cwd, "assets")
    } else {
        log.Fatalf("ERRO: Diretório de assets não encontrado")
    }

    // ========== 6. RATE LIMITERS ==========
    rateLimiter := middleware.NewRateLimiter(100, time.Minute)   // rotas gerais
    apiRateLimiter := middleware.NewRateLimiter(30, time.Minute) // rotas de API

    // ========== 7. ARQUIVOS ESTÁTICOS (COM GZIP E CACHE) ==========
    staticDir := filepath.Join(assetsDir, "static")
    staticHandler := http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir)))

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

    // Health check endpoint
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprintf(w, `{"status":"ok","timestamp":"%s"}`, time.Now().Format(time.RFC3339))
    })

    // ========== 8. ROTAS DINÂMICAS COM MIDDLEWARES ==========
    // Rota inicial
    http.HandleFunc("/", middleware.CORS(rateLimiter.Middleware(handlers.HomeHandler)))

    // Catálogo e detalhes
    http.HandleFunc("/catalog/", middleware.CORS(rateLimiter.Middleware(handlers.CatalogHandler(cachedApiClient))))
    if tmdbClient != nil {
        http.HandleFunc("/detail/", middleware.CORS(rateLimiter.Middleware(handlers.DetailHandler(cachedApiClient, cachedTmdbClient))))
    }

    // Favoritos
    http.HandleFunc("/favorites", middleware.CORS(middleware.CSRF(rateLimiter.Middleware(handlers.FavoritesHandler(cachedApiClient)))))

    // Endpoints de API
    if tmdbClient != nil {
        http.HandleFunc("/api/episodes/", middleware.CORS(apiRateLimiter.Middleware(handlers.EpisodesHandler(cachedApiClient, cachedTmdbClient))))
    }
    http.HandleFunc("/api/watch/", middleware.CORS(apiRateLimiter.Middleware(handlers.WatchHandler(cachedWatchClient))))

    // Autenticação
    http.HandleFunc("/auth/login", middleware.CORS(rateLimiter.Middleware(auth.LoginHandler)))
    http.HandleFunc("/auth/callback", middleware.CORS(rateLimiter.Middleware(auth.CallbackHandler)))

    // Endpoints adicionais de autenticação
    http.HandleFunc("/api/me", middleware.CORS(auth.MeHandler))
    http.HandleFunc("/auth/logout", middleware.CORS(auth.LogoutHandler))

    // ========== 9. CRIAR SERVIDOR COM GRACEFUL SHUTDOWN ==========
    server := &http.Server{
        Addr:         ":" + port,
        Handler:      http.DefaultServeMux,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    // Canal para sinais de encerramento
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Goroutine para shutdown graceful
    go func() {
        sig := <-sigChan
        log.Printf("📌 Recebido sinal: %v, iniciando shutdown gracioso...\n", sig)

        ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
        defer cancel()

        if err := server.Shutdown(ctx); err != nil {
            log.Printf("❌ Erro durante shutdown: %v\n", err)
        }
    }()

    // ========== 10. INICIAR SERVIDOR ==========
    fmt.Printf("🚀 Servidor TopStrem iniciado em http://localhost:%s\n", port)
    fmt.Println("📡 Health check: http://localhost:" + port + "/health")

    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("❌ Erro fatal do servidor: %v", err)
    }

    fmt.Println("✅ Servidor encerrado corretamente")
}