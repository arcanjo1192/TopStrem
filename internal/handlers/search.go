package handlers

import (
    "net/http"
    "net/url"
    "strings"

    "github.com/gin-gonic/gin"
    "topstrem/internal/api"
    "topstrem/internal/models"
)

// SearchResponse representa a resposta da API de pesquisa
type SearchResponse struct {
    Results []models.CatalogMeta `json:"results"`
}

func SearchHandler(apiClient api.CinemetaClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        query := strings.TrimSpace(c.Query("q"))
        if query == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "query é obrigatória"})
            return
        }

        results := make([]models.CatalogMeta, 0)
        var hadError bool

        // Buscar diretamente na API do Cinemeta usando o parâmetro search=
        for _, catalogType := range []string{"movie", "series"} {
            extraArgs := "search=" + url.QueryEscape(query)
            catalog, err := apiClient.GetCatalogWithFilters(catalogType, "top", extraArgs)
            if err != nil {
                hadError = true
                continue
            }
            results = append(results, catalog.Metas...)
        }

        if len(results) == 0 {
            if hadError {
                c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Cinemeta search indisponível"})
                return
            }
            c.JSON(http.StatusOK, SearchResponse{Results: results})
            return
        }

        c.JSON(http.StatusOK, SearchResponse{Results: results})
    }
}