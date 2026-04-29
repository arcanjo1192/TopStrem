package handlers

import (
    "encoding/json"
    "net/http"
    "strings"

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

func CatalogHandler(apiClient api.CinemetaClient) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 4 {
            http.NotFound(w, r)
            return
        }
        catalogType := pathParts[2]
        catalogID := pathParts[3]
		
		if catalogType != "movie" && catalogType != "series" {  
			http.Error(w, "Tipo de catálogo inválido", http.StatusBadRequest)  
			return  
		}  
		if catalogID == "" {  
			http.Error(w, "ID do catálogo não fornecido", http.StatusBadRequest)  
			return  
		}

        lang := i18n.DetectLanguage(r)

        // Obter dados (reutilizado para HTML e JSON)
        catalog, err := getCatalogData(apiClient, catalogType, catalogID)
        if err != nil {
            http.Error(w, "Erro ao carregar catálogo", http.StatusInternalServerError)
            return
        }

        // Negociar formato de resposta
        if IsJSONRequest(r) {
            // Retornar JSON para aplicativos mobile/frontend
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(CatalogDataResponse{
                Type:  catalogType,
                ID:    catalogID,
                Metas: catalog.Metas,
            })
            return
        }

        // Retornar HTML (template) para web
        templates.CatalogPage(catalog.Metas, catalogType, catalogID, lang).Render(r.Context(), w)
    }
}