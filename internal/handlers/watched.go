package handlers

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "topstrem/internal/auth"
    "topstrem/internal/storage"
)

type WatchedUpdateRequest struct {
    Action    string `json:"action"`
    EpisodeID string `json:"episodeId"`
}

func WatchedEpisodesAPIHandler(store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := auth.GetEmailFromRequest(c.Request)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
            return
        }

        watched, err := store.GetWatchedEpisodes(email)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar episódios assistidos"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"watchedEpisodes": watched})
    }
}

func UpdateWatchedEpisodesAPIHandler(store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := auth.GetEmailFromRequest(c.Request)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
            return
        }

        var payload WatchedUpdateRequest
        if err := c.BindJSON(&payload); err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
            return
        }

        if payload.EpisodeID == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "episodeId é obrigatório"})
            return
        }

        switch strings.ToLower(payload.Action) {
        case "add":
            if err := store.AddWatchedEpisode(email, payload.EpisodeID); err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao marcar episódio como assistido"})
                return
            }
        case "remove":
            if err := store.RemoveWatchedEpisode(email, payload.EpisodeID); err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover episódio assistido"})
                return
            }
        default:
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ação inválida"})
            return
        }

        watched, err := store.GetWatchedEpisodes(email)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar episódios assistidos"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"watchedEpisodes": watched})
    }
}
