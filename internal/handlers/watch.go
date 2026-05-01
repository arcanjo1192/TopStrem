package handlers

import (
    "log"
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "topstrem/internal/api"
    "topstrem/internal/models"
    "topstrem/internal/templates"
)

type WatchDataResponse struct {
    MediaType string          `json:"mediaType"`
    ID        string          `json:"id"`
    Streams   []models.Stream `json:"streams"`
}

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

func WatchHandler(watchClient api.WatchHubClientInterface, tmdbClient api.TMDBClientInterface) gin.HandlerFunc {
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

        streams, err := getWatchData(watchClient, mediaType, id)
        if err != nil || len(streams) == 0 {
            log.Printf("WatchHub vazio ou erro para %s/%s, tentando TMDB", mediaType, id)
            tmdbStreams, tmdbErr := tmdbClient.GetStreamsFromTMDB(id, mediaType)
            if tmdbErr == nil {
                streams = tmdbStreams
            } else {
                log.Printf("TMDB também falhou: %v", tmdbErr)
            }
        }
        if streams == nil {
            streams = []models.Stream{}
        }

        if IsJSONRequest(c.Request) {
            c.JSON(http.StatusOK, WatchDataResponse{
                MediaType: mediaType,
                ID:        id,
                Streams:   streams,
            })
            return
        }

        templates.WatchPage(streams, mediaType, id).Render(c.Request.Context(), c.Writer)
    }
}