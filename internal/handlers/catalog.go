package handlers

import (
    "net/http"
    "sort"
    "strings"

    "github.com/gin-gonic/gin"
    "topstrem/internal/api"
    "topstrem/internal/i18n"
    "topstrem/internal/models"
    "topstrem/internal/templates"
)

// CatalogDataResponse contém os dados brutos do catálogo para JSON
type CatalogDataResponse struct {
    Type  string              `json:"type"`
    ID    string              `json:"id"`
    Metas []models.CatalogMeta `json:"metas"`
}

// getCatalogData extrai os dados brutos do catálogo
// Isso reutiliza a mesma lógica para HTML e JSON
func getCatalogData(apiClient api.CinemetaClient, catalogType, catalogID string) (*models.CatalogResponse, error) {
    return apiClient.GetCatalog(catalogType, catalogID)
}

// getUniqueGenres coleta gêneros únicos dos metas
func getUniqueGenres(metas []models.CatalogMeta) []string {
    genreSet := make(map[string]bool)
    for _, meta := range metas {
        for _, genre := range meta.Genre {
            genreSet[genre] = true
        }
    }
    var genres []string
    for genre := range genreSet {
        genres = append(genres, genre)
    }
    // Sort for consistency
    sort.Strings(genres)
    return genres
}

// filterMetasByGenre filtra metas por gênero
func filterMetasByGenre(metas []models.CatalogMeta, genre string) []models.CatalogMeta {
    var filtered []models.CatalogMeta
    for _, meta := range metas {
        for _, g := range meta.Genre {
            if g == genre {
                filtered = append(filtered, meta)
                break
            }
        }
    }
    return filtered
}

func CatalogHandler(apiClient api.CinemetaClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        pathParts := strings.Split(c.Request.URL.Path, "/")
        if len(pathParts) < 4 {
            c.AbortWithStatus(http.StatusNotFound)
            return
        }
        catalogType := pathParts[2]
        catalogID := pathParts[3]
		
		if catalogType != "movie" && catalogType != "series" {  
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Tipo de catálogo inválido"})
			return  
		}  
		if catalogID == "" {  
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID do catálogo não fornecido"})
			return  
		}

        lang := i18n.DetectLanguage(c.Request)

        // Obter dados (reutilizado para HTML e JSON)
        catalog, err := getCatalogData(apiClient, catalogType, catalogID)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar catálogo"})
            return
        }

        // Filtrar por categoria se fornecida
        category := c.Query("category")
        var filteredMetas []models.CatalogMeta
        if category != "" {
            filteredMetas = filterMetasByGenre(catalog.Metas, category)
        } else {
            filteredMetas = catalog.Metas
        }

        // Coletar gêneros únicos para o filtro
        genres := getUniqueGenres(catalog.Metas)

        // Negociar formato de resposta
        if IsJSONRequest(c.Request) {
            // Retornar JSON para aplicativos mobile/frontend
            c.JSON(http.StatusOK, CatalogDataResponse{
                Type:  catalogType,
                ID:    catalogID,
                Metas: filteredMetas,
            })
            return
        }

        // Retornar HTML (template) para web
        templates.CatalogPage(filteredMetas, catalogType, catalogID, category, genres, lang).Render(c.Request.Context(), c.Writer)
    }
}