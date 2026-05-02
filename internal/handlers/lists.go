package handlers

import (
    "net/http"
    "strings"

    "github.com/gin-gonic/gin"
    "topstrem/internal/api"
    "topstrem/internal/auth"
	"topstrem/internal/crypto"
    "topstrem/internal/i18n"
    "topstrem/internal/models"
    "topstrem/internal/storage"
    "topstrem/internal/templates"
)

// ListsPageHandler renderiza a visualização de uma lista (grade de itens)
func ListsHandler(apiClient api.CinemetaClient, store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        lang := i18n.DetectLanguage(c.Request)
        catalogType := c.Query("type") // "movie" ou "series"
        idsParam := c.Query("ids")     // ex: "tt123,tt456"
		listName := c.Query("list")
		userToken := c.Query("user")
		isOwner := false

        if catalogType != "movie" && catalogType != "series" {
            catalogType = "movie"
        }
        if idsParam == "" {
            templates.ListPage([]models.CatalogMeta{}, catalogType, listName, lang, isOwner).Render(c.Request.Context(), c.Writer)
            return
        }
		
        email, authErr := auth.GetEmailFromRequest(c.Request)
        if authErr == nil {
            if userToken != "" {
                // Descriptografa o token e compara com o email da sessão
                decrypted, err := crypto.Decrypt(userToken)
                if err == nil && decrypted == email {
                    isOwner = true
                }
            } else if listName != "" {
                // Fallback: verifica no banco se a lista pertence ao email logado
                _, err := store.GetList(email, listName)
                if err == nil {
                    isOwner = true
                }
            }
        }

        ids := strings.Split(idsParam, ",")
        for _, id := range ids {
            id = strings.TrimSpace(id)
            if !strings.HasPrefix(id, "tt") || len(id) < 3 {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID IMDb inválido: " + id})
                return
            }
        }

        metas := make([]models.CatalogMeta, 0, len(ids))
        for _, id := range ids {
            meta, err := apiClient.GetMeta(catalogType, id)
            if err != nil {
                continue
            }
            catalogMeta := models.CatalogMeta{
                ID:         meta.ID,
                Type:       meta.Type,
                Name:       meta.Name,
                Year:       meta.Year,
                ImdbRating: meta.ImdbRating,
                Runtime:    meta.Runtime,
                Poster:     meta.Poster,
                Genre:      meta.Genre,
            }
            metas = append(metas, catalogMeta)
        }

        templates.ListPage(metas, catalogType, listName, lang, isOwner).Render(c.Request.Context(), c.Writer)
    }
}

// ListsAPIHandler retorna as listas do usuário (GET)
func ListsAPIHandler(store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := auth.GetEmailFromRequest(c.Request)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
            return
        }

        lists, err := store.GetAllLists(email)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao buscar listas"})
            return
        }

        // Filtro opcional por tipo
        filterType := c.Query("type")
        if filterType == "movie" || filterType == "series" {
            filtered := make([]storage.ListInfo, 0, len(lists))
            for _, l := range lists {
                if l.Type == filterType {
                    filtered = append(filtered, l)
                }
            }
            lists = filtered
        }

        c.JSON(http.StatusOK, gin.H{"lists": lists})
    }
}

// UpdateListsAPIHandler gerencia criação/remoção de listas e itens (POST)
func UpdateListsAPIHandler(store *storage.Storage) gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := auth.GetEmailFromRequest(c.Request)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Não autenticado"})
            return
        }

        var payload struct {
            Action   string             `json:"action"`
            ListName string             `json:"listName"`
            ListType string             `json:"listType"` // usado em create
            Item     storage.FavoriteItem `json:"item"`   // usado em add_item/remove_item
        }

        if err := c.BindJSON(&payload); err != nil {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Payload inválido"})
            return
        }

        switch strings.ToLower(payload.Action) {
		case "create":
			// Verifica se o limite de 5 listas para o tipo foi atingido
			existingLists, err := store.GetAllLists(email)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao recuperar listas"})
				return
			}
			count := 0
			for _, lst := range existingLists {
				if lst.Type == payload.ListType {
					count++
				}
			}
			if count >= 5 {
				c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Limite máximo de 5 listas de " + payload.ListType + " atingido"})
				return
			}

			if err := store.CreateList(email, payload.ListName, payload.ListType); err != nil {
				c.AbortWithStatusJSON(http.StatusConflict, gin.H{"error": err.Error()})
				return
			}
        case "delete":
            if err := store.DeleteList(email, payload.ListName); err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
                return
            }
        case "add_item":
            if err := store.AddItemToList(email, payload.ListName, payload.Item); err != nil {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }
        case "remove_item":
            if err := store.RemoveItemFromList(email, payload.ListName, payload.Item.ID); err != nil {
                c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
                return
            }
        default:
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Ação inválida"})
            return
        }

        // Retorna as listas atualizadas
        lists, err := store.GetAllLists(email)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao atualizar listas"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"lists": lists})
    }
}