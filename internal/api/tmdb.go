package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "time"
)

type TMDBClient struct {
    httpClient *http.Client
    apiKey     string
}

// ==================== Estruturas de resposta ====================
type TMDBFindResponse struct {
    MovieResults []struct {
        ID int `json:"id"`
    } `json:"movie_results"`
    TVResults []struct {
        ID int `json:"id"`
    } `json:"tv_results"`
}

type TMDBMovieDetails struct {
    ID           int    `json:"id"`
    Title        string `json:"title"`
    Overview     string `json:"overview"`
    PosterPath   string `json:"poster_path"`
    BackdropPath string `json:"backdrop_path"`
    ReleaseDate  string `json:"release_date"`
    Runtime      int    `json:"runtime"`
    Genres       []struct {
        Name string `json:"name"`
    } `json:"genres"`
    Credits struct {
        Cast []struct {
            Name string `json:"name"`
        } `json:"cast"`
        Crew []struct {
            Name string `json:"name"`
            Job  string `json:"job"`
        } `json:"crew"`
    } `json:"credits"`
    Videos struct {
        Results []struct {
            Key  string `json:"key"`
            Site string `json:"site"`
            Type string `json:"type"`
        } `json:"results"`
    } `json:"videos"`
}

type TMDBTVDetails struct {
    ID              int    `json:"id"`
    Name            string `json:"name"`
    Overview        string `json:"overview"`
    PosterPath      string `json:"poster_path"`
    BackdropPath    string `json:"backdrop_path"`
    FirstAirDate    string `json:"first_air_date"`
    NumberOfSeasons int    `json:"number_of_seasons"`
    Genres          []struct {
        Name string `json:"name"`
    } `json:"genres"`
    Credits struct {
        Cast []struct {
            Name string `json:"name"`
        } `json:"cast"`
        Crew []struct {
            Name string `json:"name"`
            Job  string `json:"job"`
        } `json:"crew"`
    } `json:"credits"`
    Videos struct {
        Results []struct {
            Key  string `json:"key"`
            Site string `json:"site"`
            Type string `json:"type"`
        } `json:"results"`
    } `json:"videos"`
}

type TMDBSeasonDetails struct {
    ID           int          `json:"id"`
    SeasonNumber int          `json:"season_number"`
    Episodes     []TMDBEpisode `json:"episodes"`
}

type TMDBEpisode struct {
    ID            int    `json:"id"`
    EpisodeNumber int    `json:"episode_number"`
    Name          string `json:"name"`
    Overview      string `json:"overview"`
    StillPath     string `json:"still_path"`
    AirDate       string `json:"air_date"`
}

// Estruturas para catálogo (discover)
type TMDBDiscoverMoviesResponse struct {
    Page    int `json:"page"`
    Results []struct {
        ID           int     `json:"id"`
        Title        string  `json:"title"`
        Overview     string  `json:"overview"`
        PosterPath   string  `json:"poster_path"`
        ReleaseDate  string  `json:"release_date"`
        VoteAverage  float64 `json:"vote_average"`
    } `json:"results"`
    TotalPages int `json:"total_pages"`
}

type TMDBDiscoverTVResponse struct {
    Page    int `json:"page"`
    Results []struct {
        ID           int     `json:"id"`
        Name         string  `json:"name"`
        Overview     string  `json:"overview"`
        PosterPath   string  `json:"poster_path"`
        FirstAirDate string  `json:"first_air_date"`
        VoteAverage  float64 `json:"vote_average"`
    } `json:"results"`
    TotalPages int `json:"total_pages"`
}

// ==================== Estruturas para Watch Providers ====================
type TMDBWatchProvidersResponse struct {
    ID      int `json:"id"`
    Results map[string]struct {
        Link     string `json:"link"`
        Flatrate []struct {
            ProviderID   int    `json:"provider_id"`
            ProviderName string `json:"provider_name"`
            LogoPath     string `json:"logo_path"`
        } `json:"flatrate"`
        Rent []struct {
            ProviderID   int    `json:"provider_id"`
            ProviderName string `json:"provider_name"`
            LogoPath     string `json:"logo_path"`
        } `json:"rent"`
        Buy []struct {
            ProviderID   int    `json:"provider_id"`
            ProviderName string `json:"provider_name"`
            LogoPath     string `json:"logo_path"`
        } `json:"buy"`
    } `json:"results"`
}

// ==================== Construtor ====================
func NewTMDBClient() (*TMDBClient, error) {  
    apiKey := os.Getenv("TMDB_API_KEY")  
    if apiKey == "" {  
        return nil, fmt.Errorf("TMDB_API_KEY environment variable is required but not set")  
    }  
      
    return &TMDBClient{  
        httpClient: &http.Client{Timeout: 10 * time.Second},  
        apiKey:     apiKey,  
    }, nil  
}

// ==================== Métodos públicos ====================
func (c *TMDBClient) FindByIMDBID(imdbID string) (int, string, error) {
    url := fmt.Sprintf("https://api.themoviedb.org/3/find/%s?api_key=%s&external_source=imdb_id", imdbID, c.apiKey)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return 0, "", err
    }
    defer resp.Body.Close()

    var findResp TMDBFindResponse
    if err := json.NewDecoder(resp.Body).Decode(&findResp); err != nil {
        return 0, "", err
    }
    if len(findResp.MovieResults) > 0 {
        return findResp.MovieResults[0].ID, "movie", nil
    }
    if len(findResp.TVResults) > 0 {
        return findResp.TVResults[0].ID, "series", nil
    }
    return 0, "", fmt.Errorf("not found")
}

func (c *TMDBClient) GetMovieDetails(tmdbID int, lang string) (*TMDBMovieDetails, error) {
    language := mapLangToTMDB(lang)
    url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d?api_key=%s&language=%s&append_to_response=credits,videos", tmdbID, c.apiKey, language)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var details TMDBMovieDetails
    if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
        return nil, err
    }
    return &details, nil
}

func (c *TMDBClient) GetTVDetails(tmdbID int, lang string) (*TMDBTVDetails, error) {
    language := mapLangToTMDB(lang)
    url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d?api_key=%s&language=%s&append_to_response=credits,videos", tmdbID, c.apiKey, language)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var details TMDBTVDetails
    if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
        return nil, err
    }
    return &details, nil
}

func (c *TMDBClient) GetTVSeriesByIMDB(imdbID string, lang string) (*TMDBTVDetails, error) {
    tmdbID, _, err := c.FindByIMDBID(imdbID)
    if err != nil {
        return nil, err
    }
    return c.GetTVDetails(tmdbID, lang)
}

func (c *TMDBClient) GetTVSeason(tmdbID int, seasonNumber int, lang string) (*TMDBSeasonDetails, error) {
    language := mapLangToTMDB(lang)
    url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/season/%d?api_key=%s&language=%s", tmdbID, seasonNumber, c.apiKey, language)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var season TMDBSeasonDetails
    if err := json.NewDecoder(resp.Body).Decode(&season); err != nil {
        return nil, err
    }
    return &season, nil
}

// DiscoverMovies retorna lista de filmes populares (para /catalog/movie/top)
func (c *TMDBClient) DiscoverMovies(lang string, page int) (*TMDBDiscoverMoviesResponse, error) {
    language := mapLangToTMDB(lang)
    url := fmt.Sprintf("https://api.themoviedb.org/3/discover/movie?api_key=%s&language=%s&sort_by=popularity.desc&page=%d", c.apiKey, language, page)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var result TMDBDiscoverMoviesResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}

// DiscoverTV retorna lista de séries populares (para /catalog/series/top)
func (c *TMDBClient) DiscoverTV(lang string, page int) (*TMDBDiscoverTVResponse, error) {
    language := mapLangToTMDB(lang)
    url := fmt.Sprintf("https://api.themoviedb.org/3/discover/tv?api_key=%s&language=%s&sort_by=popularity.desc&page=%d", c.apiKey, language, page)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var result TMDBDiscoverTVResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }
    return &result, nil
}

// GetMovieWatchProviders obtém provedores de streaming para um filme
func (c *TMDBClient) GetMovieWatchProviders(tmdbID int) (*TMDBWatchProvidersResponse, error) {
    url := fmt.Sprintf("https://api.themoviedb.org/3/movie/%d/watch/providers?api_key=%s", tmdbID, c.apiKey)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var providers TMDBWatchProvidersResponse
    if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
        return nil, err
    }
    return &providers, nil
}

// GetTVWatchProviders obtém provedores de streaming para uma série
func (c *TMDBClient) GetTVWatchProviders(tmdbID int) (*TMDBWatchProvidersResponse, error) {
    url := fmt.Sprintf("https://api.themoviedb.org/3/tv/%d/watch/providers?api_key=%s", tmdbID, c.apiKey)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    var providers TMDBWatchProvidersResponse
    if err := json.NewDecoder(resp.Body).Decode(&providers); err != nil {
        return nil, err
    }
    return &providers, nil
}

// ==================== Auxiliares ====================
func mapLangToTMDB(lang string) string {
    switch lang {
    case "pt":
        return "pt-BR"
    case "en":
        return "en-US"
    case "es":
        return "es-ES"
    case "fr":
        return "fr-FR"
    case "de":
        return "de-DE"
    case "it":
        return "it-IT"
    case "ja":
        return "ja-JP"
    case "zh":
        return "zh-CN"
    case "ru":
        return "ru-RU"
    case "ar":
        return "ar-SA"
    case "hi":
        return "hi-IN"
    case "ko":
        return "ko-KR"
    default:
        return "pt-BR"
    }
}