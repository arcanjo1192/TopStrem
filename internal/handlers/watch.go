package handlers

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
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

func WatchHandler(watchClient api.WatchHubClientInterface) gin.HandlerFunc {
    return func(c *gin.Context) {
        pathParts := strings.Split(c.Request.URL.Path, "/")
        if len(pathParts) < 5 {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Requisição inválida"})
            return
        }
        mediaType := pathParts[3]
        id := pathParts[4]
		
		if mediaType != "movie" && mediaType != "series" {  
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Tipo de mídia inválido"})
			return  
		}  
		if id == "" {  
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID não fornecido"})
			return  
		}

        // Obter dados (reutilizado para HTML e JSON)
        streams, err := getWatchData(watchClient, mediaType, id)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
            return
        }

        // Negociar formato de resposta
        if IsJSONRequest(c.Request) {
            // Retornar JSON para aplicativos mobile/frontend
            c.JSON(http.StatusOK, WatchDataResponse{
                MediaType: mediaType,
                ID:        id,
                Streams:   streams,
            })
            return
        }

        // Retornar HTML (template) para web
        templates.WatchPage(streams, mediaType, id).Render(c.Request.Context(), c.Writer)
    }
}