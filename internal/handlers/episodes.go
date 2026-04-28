package handlers

import (
    "encoding/json"
    "net/http"
    "sort"
    "strconv"
    "strings"

    "topstrem/internal/api"
    "topstrem/internal/models"
)

type SeasonGroup struct {
    SeasonNumber int            `json:"season"`
    Episodes     []models.Video `json:"episodes"`
}

func EpisodesHandler(apiClient *api.Client, tmdbClient *api.TMDBClient) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        pathParts := strings.Split(r.URL.Path, "/")
        if len(pathParts) < 4 {
            http.Error(w, "ID não fornecido", http.StatusBadRequest)
            return
        }
        id := pathParts[3]

        lang := r.URL.Query().Get("lang")
        if lang == "" {
            lang = "pt"
        }

        // Tenta TMDB
        series, err := tmdbClient.GetTVSeriesByIMDB(id, lang)
        if err == nil && series != nil && series.ID != 0 {
            seasonsMap := make(map[int][]models.Video)
            for i := 1; i <= series.NumberOfSeasons; i++ {
                seasonDetails, err := tmdbClient.GetTVSeason(series.ID, i, lang)
                if err != nil {
                    continue
                }
                for _, ep := range seasonDetails.Episodes {
                    video := models.Video{
                        ID:        strconv.Itoa(ep.ID), // int to string
                        Name:      ep.Name,
                        Season:    seasonDetails.SeasonNumber,
                        Episode:   ep.EpisodeNumber,
                        Released:  ep.AirDate,
                        Thumbnail: "https://image.tmdb.org/t/p/w500" + ep.StillPath,
                    }
                    seasonsMap[seasonDetails.SeasonNumber] = append(seasonsMap[seasonDetails.SeasonNumber], video)
                }
            }

            var seasons []SeasonGroup
            for num, eps := range seasonsMap {
                sort.Slice(eps, func(i, j int) bool {
                    return eps[i].Episode < eps[j].Episode
                })
                seasons = append(seasons, SeasonGroup{SeasonNumber: num, Episodes: eps})
            }
            sort.Slice(seasons, func(i, j int) bool {
                return seasons[i].SeasonNumber < seasons[j].SeasonNumber
            })

            w.Header().Set("Content-Type", "application/json")
            json.NewEncoder(w).Encode(seasons)
            return
        }

        // Fallback para Cinemata
        meta, err := apiClient.GetMeta("series", id)
        if err != nil || meta == nil {
            http.Error(w, "Série não encontrada", http.StatusNotFound)
            return
        }

        seasonMap := make(map[int][]models.Video)
        for _, v := range meta.Videos {
            seasonMap[v.Season] = append(seasonMap[v.Season], v)
        }

        var seasons []SeasonGroup
        for num, eps := range seasonMap {
            sort.Slice(eps, func(i, j int) bool {
                return eps[i].Episode < eps[j].Episode
            })
            seasons = append(seasons, SeasonGroup{
                SeasonNumber: num,
                Episodes:     eps,
            })
        }
        sort.Slice(seasons, func(i, j int) bool {
            return seasons[i].SeasonNumber < seasons[j].SeasonNumber
        })

        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(seasons)
    }
}