package main

import (
    "fmt"
    "net/http"
    "os"
    "path/filepath"

    "topstrem/internal/api"
    "topstrem/internal/auth"
    "topstrem/internal/handlers"
)

func main() {
    apiClient := api.NewClient()
    watchClient := api.NewWatchHubClient()
    tmdbClient := api.NewTMDBClient()

    cwd, err := os.Getwd()
    if err != nil {
        panic(err)
    }

    // Detecta assetsDir
    assetsDir := ""
    if _, err := os.Stat(filepath.Join(cwd, "cmd/app/assets")); err == nil {
        assetsDir = filepath.Join(cwd, "cmd/app/assets")
    } else if _, err := os.Stat(filepath.Join(cwd, "assets")); err == nil {
        assetsDir = filepath.Join(cwd, "assets")
    } else {
        panic("Assets directory not found")
    }

    // Servir arquivos estáticos
    staticDir := filepath.Join(assetsDir, "static")
    http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(staticDir))))

    // Rotas existentes
    http.HandleFunc("/", handlers.HomeHandler)
    http.HandleFunc("/catalog/", handlers.CatalogHandler(apiClient))
    // CORREÇÃO: use apenas o handler com dois argumentos
    http.HandleFunc("/detail/", handlers.DetailHandler(apiClient, tmdbClient))
    http.HandleFunc("/favorites", handlers.FavoritesHandler(apiClient))
    http.HandleFunc("/api/episodes/", handlers.EpisodesHandler(apiClient, tmdbClient))
    http.HandleFunc("/api/watch/", handlers.WatchHandler(watchClient))
    http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filepath.Join(assetsDir, "robots.txt"))
    })
    http.Handle("/privacy/", http.StripPrefix("/privacy/", http.FileServer(http.Dir(filepath.Join(assetsDir, "privacy")))))

    // Configura autenticação
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

    // Rotas de autenticação
    http.HandleFunc("/auth/login", auth.LoginHandler)
    http.HandleFunc("/auth/callback", auth.CallbackHandler)

    // Página de callback
    http.HandleFunc("/callback.html", func(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, filepath.Join(assetsDir, "callback.html"))
    })

    fmt.Println("Servidor rodando na porta 8080...")
    http.ListenAndServe(":8080", nil)
}