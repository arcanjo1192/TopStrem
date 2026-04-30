package handlers

import (
    "net/http"
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

        lowerQuery := strings.ToLower(query)
        results := make([]models.CatalogMeta, 0)
        hadCatalogError := false

        // Buscar nos catálogos populares de filmes e séries
        for _, catalogType := range []string{"movie", "series"} {
            catalog, err := apiClient.GetCatalog(catalogType, "top")
            if err != nil {
                hadCatalogError = true
                continue
            }
            for _, meta := range catalog.Metas {
                if strings.Contains(strings.ToLower(meta.Name), lowerQuery) {
                    results = append(results, meta)
                }
            }
        }

        if len(results) == 0 {
            if hadCatalogError {
                c.AbortWithStatusJSON(http.StatusBadGateway, gin.H{"error": "Cinemeta search indisponível"})
                return
            }
            c.JSON(http.StatusOK, SearchResponse{Results: results})
            return
        }

        c.JSON(http.StatusOK, SearchResponse{Results: results})
    }
}
