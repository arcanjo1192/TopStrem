package handlers

import (
    "encoding/json"
    "net/http"
    "sort"
    "strings"

    "topstrem/internal/api"
    "topstrem/internal/models"
    "topstrem/internal/templates"
)

// EpisodesDataResponse contém os dados brutos de episódios para JSON
type EpisodesDataResponse struct {
    SeriesID string               `json:"seriesId"`
    Seasons  []models.SeasonGroup `json:"seasons"`
}

// getEpisodesData extrai os dados brutos de episódios
// Isso reutiliza a mesma lógica para HTML e JSON
func getEpisodesData(id string, lang string, apiClient api.CinemetaClient, tmdbClient api.TMDBClientInterface) ([]models.SeasonGroup, *models.Meta, error) {
    // 1. Obter sempre a estrutura base do Cinemeta
    meta, err := apiClient.GetMeta("series", id)
    if err != nil || meta == nil {
        return nil, nil, err
    }

    // 2. Agrupar episódios por temporada (dados originais do Cinemeta)
    seasonMap := make(map[int][]models.Video)
    for _, v := range meta.Videos {
        seasonMap[v.Season] = append(seasonMap[v.Season], v)
    }

    // 3. Tentar obter os títulos do TMDB (apenas se o cliente TMDB existir)
    if tmdbClient != nil {
        // 3.1 Descobrir o ID da série no TMDB
        tmdbID, _, err := tmdbClient.FindByIMDBID(id)
        if err == nil && tmdbID != 0 {
            // 3.2 Para cada temporada com episódios, buscar detalhes do TMDB
            tmdbSeasons := make(map[int]map[int]string) // temporada -> númeroEpisódio -> nome

            for seasonNum := range seasonMap {
                details, err := tmdbClient.GetTVSeason(tmdbID, seasonNum, lang)
                if err != nil {
                    continue
                }
                tmdbEpisodes := make(map[int]string)
                for _, ep := range details.Episodes {
                    tmdbEpisodes[ep.EpisodeNumber] = ep.Name
                }
                if len(tmdbEpisodes) > 0 {
                    tmdbSeasons[seasonNum] = tmdbEpisodes
                }
            }

            // 3.3 Substituir os nomes apenas onde o TMDB tiver informação
            for season, episodes := range seasonMap {
                if tmdbEps, ok := tmdbSeasons[season]; ok {
                    for i, ep := range episodes {
                        if name, exists := tmdbEps[ep.Episode]; exists && name != "" {
                            episodes[i].Name = name
                        }
                    }
                    seasonMap[season] = episodes
                }
            }
        }
    }

    // 4. Ordenar e montar resposta
    var seasons []models.SeasonGroup
    for num, eps := range seasonMap {
        sort.Slice(eps, func(i, j int) bool {
            return eps[i].Episode < eps[j].Episode
        })
        seasons = append(seasons, models.SeasonGroup{
            SeasonNumber: num,
            Episodes:     eps,
        })
    }
    sort.Slice(seasons, func(i, j int) bool {
        return seasons[i].SeasonNumber < seasons[j].SeasonNumber
    })

    return seasons, meta, nil
}

func EpisodesHandler(apiClient api.CinemetaClient, tmdbClient api.TMDBClientInterface) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 4 {
            http.Error(w, "ID não fornecido", http.StatusBadRequest)
            return
        }
        id := pathParts[3]
        if !strings.HasPrefix(id, "tt") || len(id) < 3 {
            http.Error(w, "ID IMDb inválido", http.StatusBadRequest)
            return
        }

        lang := r.URL.Query().Get("lang")
        if len(lang) > 0 && len(lang) < 2 {
            http.Error(w, "Código de idioma inválido", http.StatusBadRequest)
            return
        }
        if lang == "" {
            lang = "pt"
        }

        // Obter dados (reutilizado para HTML e JSON)
        seasons, meta, err := getEpisodesData(id, lang, apiClient, tmdbClient)
        if err != nil {
            http.Error(w, "Série não encontrada", http.StatusNotFound)
            return
        }

        // Negociar formato de resposta
        if IsJSONRequest(r) {
            // Retornar JSON para aplicativos mobile/frontend
            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(EpisodesDataResponse{
                SeriesID: id,
                Seasons:  seasons,
            })
            return
        }

        // Retornar HTML (template) para web
        // Passar meta e seasons para o template
        templates.EpisodesPage(meta, seasons, lang).Render(r.Context(), w)
    }
}