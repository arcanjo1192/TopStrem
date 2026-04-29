package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

    "topstrem/internal/api"
    "topstrem/internal/models"
    "topstrem/internal/templates"
)

// WatchDataResponse contém os dados brutos de streams para JSON
type WatchDataResponse struct {
    MediaType string            `json:"mediaType"`
    ID        string            `json:"id"`
    Streams   []models.Stream   `json:"streams"`
}

// getWatchData extrai os dados brutos de streams
// Isso reutiliza a mesma lógica para HTML e JSON
func getWatchData(watchClient api.WatchHubClientInterface, mediaType, id string) ([]models.Stream, error) {
    response, err := watchClient.GetStreams(mediaType, id)
    if err != nil {
        return nil, err
    }
    if response == nil {
        return []models.Stream{}, nil
    }
    return response.Streams, nil
}

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

        // Obter dados (reutilizado para HTML e JSON)
        streams, err := getWatchData(watchClient, mediaType, id)
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }

        // Negociar formato de resposta
        if IsJSONRequest(r) {
            // Retornar JSON para aplicativos mobile/frontend
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(WatchDataResponse{
                MediaType: mediaType,
                ID:        id,
                Streams:   streams,
            })
            return
        }

        // Retornar HTML (template) para web
        templates.WatchPage(streams, mediaType, id).Render(r.Context(), w)
    }
}