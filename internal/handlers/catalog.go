package handlers

import (
    "fmt"
    "net/http"
    "net/url"
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
    Type  string                `json:"type"`
    ID    string                `json:"id"`
    Metas []models.CatalogMeta  `json:"metas"`
}

// getUniqueGenres coleta gêneros únicos dos metas (fallback)
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
    sort.Strings(genres)
    return genres
}

// getGenresFromManifest extrai a lista de gêneros do manifesto para um catálogo específico
func getGenresFromManifest(client api.CinemetaClient, catalogType, catalogID string) ([]string, error) {
    manifest, err := client.GetManifest()
    if err != nil {
        return nil, err
    }
    for _, cat := range manifest.Catalogs {
        if cat.Type == catalogType && cat.ID == catalogID {
            return cat.Genres, nil
        }
    }
    return nil, fmt.Errorf("catálogo não encontrado no manifesto")
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

        // Obter gêneros disponíveis para o dropdown (do manifesto ou fallback)
        genres, err := getGenresFromManifest(apiClient, catalogType, catalogID)
        if err != nil {
            // fallback: busca o catálogo completo para extrair gêneros
            catalogFallback, fallbackErr := apiClient.GetCatalog(catalogType, catalogID)
            if fallbackErr == nil && catalogFallback != nil {
                genres = getUniqueGenres(catalogFallback.Metas)
            } else {
                genres = []string{}
            }
        }

        // Obter catálogo, com ou sem filtro de gênero (usando rota Stremio)
        category := c.Query("category")
        var metas []models.CatalogMeta
        if category != "" {
            extraArgs := "genre=" + url.QueryEscape(category)   // formato exigido pela API
            catalog, err := apiClient.GetCatalogWithFilters(catalogType, catalogID, extraArgs)
            if err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar catálogo"})
                return
            }
            metas = catalog.Metas
        } else {
            catalog, err := apiClient.GetCatalog(catalogType, catalogID)
            if err != nil {
                c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Erro ao carregar catálogo"})
                return
            }
            metas = catalog.Metas
        }

        // Limite de 50 títulos
        if len(metas) > 50 {
            metas = metas[:50]
        }

        // Negociar formato de resposta
        if IsJSONRequest(c.Request) {
            c.JSON(http.StatusOK, CatalogDataResponse{
                Type:  catalogType,
                ID:    catalogID,
                Metas: metas,
            })
            return
        }

        // Renderizar HTML
        templates.CatalogPage(metas, catalogType, catalogID, category, genres, lang).Render(c.Request.Context(), c.Writer)
    }
}