package handlers

import (
    "net/http"
    "strings"

    "topstrem/internal/api"
    "topstrem/internal/i18n"
    "topstrem/internal/templates"
)

func CatalogHandler(apiClient *api.Client) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 4 {
            http.NotFound(w, r)
            return
        }
        catalogType := pathParts[2]
        catalogID := pathParts[3]

        lang := i18n.DetectLanguage(r)

        catalog, err := apiClient.GetCatalog(catalogType, catalogID)
        if err != nil {
            http.Error(w, "Erro ao carregar catálogo", http.StatusInternalServerError)
            return
        }

        templates.CatalogPage(catalog.Metas, catalogType, catalogID, lang).Render(r.Context(), w)
    }
}