package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "topstrem/internal/api"
)

func WatchHandler(watchClient *api.WatchHubClient) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 5 {
            http.Error(w, "Requisição inválida", http.StatusBadRequest)
            return
        }
        mediaType := pathParts[3]
        id := pathParts[4]

        streams, err := watchClient.GetStreams(mediaType, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(streams)
    }
}