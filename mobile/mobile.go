package mobile

import (
    "net/http"

    "topstrem/internal/api"
    "topstrem/internal/handlers"
)

func StartServer() {
    // Cliente para catálogo
    cinemetaClient := api.NewClient()

    // Cliente TMDB
    tmdbClient, err := api.NewTMDBClient()
    if err != nil {
        panic("falha ao criar cliente TMDB: " + err.Error())
    }

    // Cliente WatchHub
    watchClient := api.NewWatchHubClient()

    // Rotas - agora com os argumentos corretos
    http.HandleFunc("/", handlers.HomeHandler)
    http.HandleFunc("/catalog/", handlers.CatalogHandler(cinemetaClient))
    http.HandleFunc("/detail/", handlers.DetailHandler(cinemetaClient, tmdbClient))
    http.HandleFunc("/favorites", handlers.FavoritesHandler(cinemetaClient))
    http.HandleFunc("/api/episodes/", handlers.EpisodesHandler(cinemetaClient, tmdbClient))
    http.HandleFunc("/api/watch/", handlers.WatchHandler(watchClient))

    http.ListenAndServe(":8080", nil)
}