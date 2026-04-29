package handlers

import (
    "context"
    "encoding/json"
    "errors"
    "net/http"
    "strings"
    "time"

    "golang.org/x/sync/errgroup"
    "topstrem/internal/api"
    "topstrem/internal/i18n"
    "topstrem/internal/models"
    "topstrem/internal/templates"
)

// DetailDataResponse contém os dados brutos de detalhe para JSON
type DetailDataResponse struct {
    MediaType string       `json:"mediaType"`
    ID        string       `json:"id"`
    Meta      models.Meta  `json:"meta"`
}

// getDetailData extrai os dados enriquecidos de detalhe
// Isso reutiliza a mesma lógica para HTML e JSON
func getDetailData(ctx context.Context, mediaType, id string, apiClient api.CinemetaClient, tmdbClient api.TMDBClientInterface) (*models.Meta, error) {
    // Estruturas para resultados concorrentes
    var (
        meta        *models.Meta
        cinemataErr error
        tmdbData    struct {
            found       bool
            tmdbID      int
            tmdbType    string
            details     interface{} // pode ser *TMDBMovieDetails ou *TMDBSeriesDetails
            trailers    []models.Trailer
            name, overview string
        }
    )

    g, ctx := errgroup.WithContext(ctx)

    // Goroutine 1: Buscar dados do Cinemata (metadados base)
    g.Go(func() error {
        m, err := apiClient.GetMeta(mediaType, id)
        if err != nil {
            cinemataErr = err
            return err
        }
        meta = m
        return nil
    })

    // Goroutine 2: Buscar dados complementares do TMDB (se cliente existir)
    if tmdbClient != nil {
        g.Go(func() error {
            // 2.1 Encontrar ID no TMDB
            tmdbID, tmdbMediaType, err := tmdbClient.FindByIMDBID(id)
            if err != nil {
                // Não falha a página inteira se TMDB falhar
                return nil
            }
            tmdbData.found = true
            tmdbData.tmdbID = tmdbID
            tmdbData.tmdbType = tmdbMediaType

            // 2.2 Buscar detalhes (filme ou série)
            if tmdbMediaType == "movie" {
                details, err := tmdbClient.GetMovieDetails(tmdbID, "")
                if err != nil {
                    return nil
                }
                tmdbData.details = details
                if details.Title != "" {
                    tmdbData.name = details.Title
                }
                if details.Overview != "" {
                    tmdbData.overview = details.Overview
                }
                // Extrair trailers
                for _, video := range details.Videos.Results {
                    if video.Site == "YouTube" && video.Type == "Trailer" {
                        tmdbData.trailers = append(tmdbData.trailers, models.Trailer{
                            Source: video.Key,
                            Type:   video.Type,
                        })
                    }
                }
            } else if tmdbMediaType == "series" {
                details, err := tmdbClient.GetTVDetails(tmdbID, "")
                if err != nil {
                    return nil
                }
                tmdbData.details = details
                if details.Name != "" {
                    tmdbData.name = details.Name
                }
                if details.Overview != "" {
                    tmdbData.overview = details.Overview
                }
                for _, video := range details.Videos.Results {
                    if video.Site == "YouTube" && video.Type == "Trailer" {
                        tmdbData.trailers = append(tmdbData.trailers, models.Trailer{
                            Source: video.Key,
                            Type:   video.Type,
                        })
                    }
                }
            }
            return nil
        })
    }

    // Aguarda conclusão de ambas goroutines (ou timeout)
    if err := g.Wait(); err != nil {
        // Se o Cinemata falhou, retorna erro (TMDB opcional)
        if cinemataErr != nil {
            return nil, cinemataErr
        }
        // Se apenas o TMDB falhou, segue com os dados do Cinemata
    }

    // Se o TMDB retornou dados melhores, enriquece o meta
    if tmdbData.found && meta != nil {
        // Prioriza nome e sinopse do TMDB se disponíveis
        if tmdbData.name != "" {
            meta.Name = tmdbData.name
        }
        if tmdbData.overview != "" {
            meta.Description = tmdbData.overview
        }
        // Insere trailers do TMDB no início (prioridade)
        if len(tmdbData.trailers) > 0 {
            meta.Trailers = append(tmdbData.trailers, meta.Trailers...)
        }
    }

    // Se mesmo após paralelismo o meta for nulo (caso extremo)
    if meta == nil {
        return nil, errors.New("conteúdo não encontrado")
    }

    return meta, nil
}

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

        // Contexto com timeout global (ex: 5 segundos)
        ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
        defer cancel()

        // Obter dados enriquecidos (reutilizado para HTML e JSON)
        meta, err := getDetailData(ctx, mediaType, id, apiClient, tmdbClient)
        if err != nil {
            http.Error(w, "Título não encontrado", http.StatusNotFound)
            return
        }

        // Negociar formato de resposta
        if IsJSONRequest(r) {
            // Retornar JSON para aplicativos mobile/frontend
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(DetailDataResponse{
                MediaType: mediaType,
                ID:        id,
                Meta:      *meta,
            })
            return
        }

        // Retornar HTML (template) para web
        templates.DetailPage(*meta, lang).Render(r.Context(), w)
    }
}