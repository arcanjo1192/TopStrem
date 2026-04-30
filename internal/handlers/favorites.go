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

func FavoritesHandler(apiClient api.CinemetaClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        lang := i18n.DetectLanguage(c.Request)
        catalogType := c.Query("type") // "movie" ou "series"
        idsParam := c.Query("ids")      // ex: "tt123,tt456"

        if catalogType != "movie" && catalogType != "series" {
            catalogType = "movie"
        }
        if idsParam == "" {
            // Sem IDs, exibe grade vazia
            templates.CatalogPage([]models.CatalogMeta{}, catalogType, "favorites", lang).Render(c.Request.Context(), c.Writer)
            return
        }

        ids := strings.Split(idsParam, ",")
		
		for _, id := range ids {  
			id = strings.TrimSpace(id)  
			if !strings.HasPrefix(id, "tt") || len(id) < 3 {  
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID IMDb inválido nos favoritos: " + id})
				return  
			}  
		}
				
        metas := make([]models.CatalogMeta, 0, len(ids))

        for _, id := range ids {
            meta, err := apiClient.GetMeta(catalogType, id)
            if err != nil {
                continue
            }
            // Converte Meta para CatalogMeta
            catalogMeta := models.CatalogMeta{
                ID:     meta.ID,
                Type:   meta.Type,
                Name:   meta.Name,
                Year:   meta.Year,
                Poster: meta.Poster,
                Genre:  meta.Genre,
            }
            metas = append(metas, catalogMeta)
        }

        // Renderiza a página com os favoritos
        templates.CatalogPage(metas, catalogType, "favorites", lang).Render(c.Request.Context(), c.Writer)
    }
}