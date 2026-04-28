package handlers

import (
    "net/http"
    "strings"
    "topstrem/internal/api"
    "topstrem/internal/i18n"
    "topstrem/internal/models"
    "topstrem/internal/templates"
)

func FavoritesHandler(apiClient *api.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        lang := i18n.DetectLanguage(r)
        catalogType := r.URL.Query().Get("type") // "movie" ou "series"
        idsParam := r.URL.Query().Get("ids")      // ex: "tt123,tt456"

        if catalogType != "movie" && catalogType != "series" {
            catalogType = "movie"
        }
        if idsParam == "" {
            // Sem IDs, exibe grade vazia
            templates.CatalogPage([]models.CatalogMeta{}, catalogType, "favorites", lang).Render(r.Context(), w)
            return
        }

        ids := strings.Split(idsParam, ",")
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
        templates.CatalogPage(metas, catalogType, "favorites", lang).Render(r.Context(), w)
    }
}