package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "topstrem/internal/api"
)

func WatchHandler(watchClient api.WatchHubClientInterface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 5 {
            http.Error(w, "Requisição inválida", http.StatusBadRequest)
            return
        }
        mediaType := pathParts[3]
        id := pathParts[4]
		
		if mediaType != "movie" && mediaType != "series" {  
			http.Error(w, "Tipo de mídia inválido", http.StatusBadRequest)  
			return  
		}  
		if id == "" {  
			http.Error(w, "ID não fornecido", http.StatusBadRequest)  
			return  
		}

        streams, err := watchClient.GetStreams(mediaType, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(streams)
    }
}