package handlers

import (
    "net/http"
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

        // Negociar formato de resposta
        if IsJSONRequest(c.Request) {
            // Retornar JSON para aplicativos mobile/frontend
            c.JSON(http.StatusOK, CatalogDataResponse{
                Type:  catalogType,
                ID:    catalogID,
                Metas: catalog.Metas,
            })
            return
        }

        // Retornar HTML (template) para web
        templates.CatalogPage(catalog.Metas, catalogType, catalogID, lang).Render(c.Request.Context(), c.Writer)
    }
}