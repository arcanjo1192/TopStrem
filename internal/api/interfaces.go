package api

import "topstrem/internal/models"

type CinemetaClient interface {
    GetCatalog(catalogType, catalogID string) (*models.CatalogResponse, error)
    GetCatalogWithFilters(catalogType, catalogID, extraArgs string) (*models.CatalogResponse, error)
    GetMeta(mediaType, id string) (*models.Meta, error)
    GetManifest() (*models.ManifestResponse, error)
}

type TMDBClientInterface interface {
    FindByIMDBID(imdbID string) (int, string, error)
    GetMovieDetails(tmdbID int, lang string) (*TMDBMovieDetails, error)
    GetTVDetails(tmdbID int, lang string) (*TMDBTVDetails, error)
    GetTVSeriesByIMDB(imdbID string, lang string) (*TMDBTVDetails, error)
    GetTVSeason(tmdbID int, seasonNumber int, lang string) (*TMDBSeasonDetails, error)
    GetStreamsFromTMDB(imdbID string, mediaType string) ([]models.Stream, error) // novo
}

type WatchHubClientInterface interface {
    GetStreams(mediaType, id string) (*models.StreamResponse, error)
}