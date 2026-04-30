package handlers

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "topstrem/internal/api"
    "topstrem/internal/auth"
    "topstrem/internal/i18n"
    "topstrem/internal/models"
    "topstrem/internal/storage"
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

type FavoriteUpdateRequest struct {
    Action string               `json:"action"`
    Item   storage.FavoriteItem `json:"item"`
}

func FavoritesAPIHandler(store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := auth.GetEmailFromRequest(c.Request)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
            return
        }

        favorites, err := store.GetFavorites(email)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar favoritos"})
            return
        }

        filterType := c.Query("type")
        if filterType == "movie" || filterType == "series" {
            filtered := make([]storage.FavoriteItem, 0, len(favorites))
            for _, item := range favorites {
                if item.Type == filterType {
                    filtered = append(filtered, item)
                }
            }
            favorites = filtered
        }

        c.JSON(http.StatusOK, gin.H{"favorites": favorites})
    }
}

func UpdateFavoritesAPIHandler(store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := auth.GetEmailFromRequest(c.Request)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
            return
        }

        var payload FavoriteUpdateRequest
        if err := c.BindJSON(&payload); err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
            return
        }

        if payload.Item.ID == "" {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID é obrigatório"})
            return
        }

        switch strings.ToLower(payload.Action) {
        case "add":
            if payload.Item.Type == "" {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Type é obrigatório para adicionar favorito"})
                return
            }
            if err := store.AddFavorite(email, payload.Item); err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao salvar favorito"})
                return
            }
        case "remove":
            if err := store.RemoveFavorite(email, payload.Item.ID); err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao remover favorito"})
                return
            }
        default:
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ação inválida"})
            return
        }

        favorites, err := store.GetFavorites(email)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar favoritos"})
            return
        }

        c.JSON(http.StatusOK, gin.H{"favorites": favorites})
    }
}
