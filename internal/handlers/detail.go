package handlers

import (
    "net/http"
    "strings"

    "topstrem/internal/api"
    "topstrem/internal/i18n"
    "topstrem/internal/models"
    "topstrem/internal/templates"
)

func DetailHandler(apiClient api.CinemetaClient, tmdbClient api.TMDBClientInterface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 4 {
            http.NotFound(w, r)
            return
        }
        mediaType := pathParts[2]
        id := pathParts[3] // IMDb ID, ex: tt0133093
		
		if mediaType != "movie" && mediaType != "series" {  
			http.Error(w, "Tipo de mídia inválido", http.StatusBadRequest)  
			return  
		}  
		if !strings.HasPrefix(id, "tt") || len(id) < 3 {  
			http.Error(w, "ID IMDb inválido", http.StatusBadRequest)  
			return  
		}

        lang := i18n.DetectLanguage(r)

        // Busca dados do Cinemata (inclui trailers originais)
        meta, err := apiClient.GetMeta(mediaType, id)
        if err != nil {
            http.Error(w, "Título não encontrado", http.StatusNotFound)
            return
        }

        // Tenta enriquecer com TMDB (título, sinopse e trailers)
        if tmdbClient != nil {
            tmdbID, tmdbMediaType, err := tmdbClient.FindByIMDBID(id)
            if err == nil && (tmdbMediaType == "movie" || tmdbMediaType == "series") {
                var tmdbTrailers []models.Trailer
                if tmdbMediaType == "movie" {
                    details, err := tmdbClient.GetMovieDetails(tmdbID, lang)
                    if err == nil && details != nil {
                        // Título e sinopse
                        if details.Title != "" {
                            meta.Name = details.Title
                        }
                        if details.Overview != "" {
                            meta.Description = details.Overview
                        }
                        // Trailers (YouTube, tipo "Trailer")
                        for _, video := range details.Videos.Results {
                            if video.Site == "YouTube" && video.Type == "Trailer" {
                                tmdbTrailers = append(tmdbTrailers, models.Trailer{
                                    Source: video.Key,
                                    Type:   video.Type,
                                })
                            }
                        }
                    }
                } else { // series
                    details, err := tmdbClient.GetTVDetails(tmdbID, lang)
                    if err == nil && details != nil {
                        if details.Name != "" {
                            meta.Name = details.Name
                        }
                        if details.Overview != "" {
                            meta.Description = details.Overview
                        }
                        for _, video := range details.Videos.Results {
                            if video.Site == "YouTube" && video.Type == "Trailer" {
                                tmdbTrailers = append(tmdbTrailers, models.Trailer{
                                    Source: video.Key,
                                    Type:   video.Type,
                                })
                            }
                        }
                    }
                }

                // Prioriza trailers do TMDB (insere no início da lista)
                if len(tmdbTrailers) > 0 {
                    meta.Trailers = append(tmdbTrailers, meta.Trailers...)
                }
            }
        }

        templates.DetailPage(*meta, lang).Render(r.Context(), w)
    }
}