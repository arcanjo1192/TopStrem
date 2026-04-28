// mobile/mobile.go
package mobile

import (
    "net/http"

    "topstrem/internal/api"
    "topstrem/internal/handlers"
)

// StartServer inicia o servidor HTTP.
func StartServer() {
    apiClient := api.NewClient()
    watchClient := api.NewWatchHubClient()

    // Rotas públicas
    http.HandleFunc("/", handlers.HomeHandler)
    http.HandleFunc("/catalog/", handlers.CatalogHandler(apiClient))
    http.HandleFunc("/detail/", handlers.DetailHandler(apiClient))
    http.HandleFunc("/favorites", handlers.FavoritesHandler(apiClient))
    http.HandleFunc("/api/episodes/", handlers.EpisodesHandler(apiClient))
    http.HandleFunc("/api/watch/", handlers.WatchHandler(watchClient))

    // Inicia o servidor
    go http.ListenAndServe(":8080", nil)
}