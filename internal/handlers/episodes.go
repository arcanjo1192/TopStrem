package handlers

import (
    "net/http"
    "sort"
    "strings"

    "github.com/gin-gonic/gin"
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

func EpisodesHandler(apiClient api.CinemetaClient, tmdbClient api.TMDBClientInterface) gin.HandlerFunc {
    return func(c *gin.Context) {
        pathParts := strings.Split(c.Request.URL.Path, "/")
        if len(pathParts) < 4 {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID não fornecido"})
            return
        }
        id := pathParts[3]
        if !strings.HasPrefix(id, "tt") || len(id) < 3 {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "ID IMDb inválido"})
            return
        }

        lang := c.Query("lang")
        if len(lang) > 0 && len(lang) < 2 {
            c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "Código de idioma inválido"})
            return
        }
        if lang == "" {
            lang = "pt"
        }

        // Obter dados (reutilizado para HTML e JSON)
        seasons, meta, err := getEpisodesData(id, lang, apiClient, tmdbClient)
        if err != nil {
            c.AbortWithStatusJSON(http.StatusNotFound, gin.H{"error": "Série não encontrada"})
            return
        }

        // Negociar formato de resposta
        if IsJSONRequest(c.Request) {
            // Retornar JSON para aplicativos mobile/frontend
            c.JSON(http.StatusOK, EpisodesDataResponse{
                SeriesID: id,
                Seasons:  seasons,
            })
            return
        }

        // Retornar HTML (template) para web
        // Passar meta e seasons para o template
        templates.EpisodesPage(meta, seasons, lang).Render(c.Request.Context(), c.Writer)
    }
}