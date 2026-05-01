// internal/api/cached_tmdb.go  
package api  
  
import (  
    "context"  
    "fmt"  
    "time"  
      
    "topstrem/internal/cache"  
	"topstrem/internal/models"
)  
  
type CachedTMDBClient struct {  
    client *TMDBClient  
    cache  *cache.RedisCache  
}  
  
func NewCachedTMDBClient(client *TMDBClient, cache *cache.RedisCache) *CachedTMDBClient {  
    return &CachedTMDBClient{  
        client: client,  
        cache:  cache,  
    }  
}  
  
func (c *CachedTMDBClient) FindByIMDBID(imdbID string) (int, string, error) {  
    ctx := context.Background()  
    key := fmt.Sprintf("tmdb:find:%s", imdbID)  
      
    type FindResult struct {  
        ID   int    `json:"id"`  
        Type string `json:"type"`  
    }  
      
    var result FindResult  
    if err := c.cache.Get(ctx, key, &result); err == nil {  
        return result.ID, result.Type, nil  
    }  
      
    id, mediaType, err := c.client.FindByIMDBID(imdbID)  
    if err != nil {  
        return 0, "", err  
    }  
      
    // Cache por 24 horas (IDs não mudam)  
    c.cache.Set(ctx, key, FindResult{ID: id, Type: mediaType}, 24*time.Hour)  
    return id, mediaType, nil  
}  
  
func (c *CachedTMDBClient) GetMovieDetails(tmdbID int, lang string) (*TMDBMovieDetails, error) {  
    ctx := context.Background()  
    key := fmt.Sprintf("tmdb:movie:%d:%s", tmdbID, lang)  
      
    var result TMDBMovieDetails  
    if err := c.cache.Get(ctx, key, &result); err == nil {  
        return &result, nil  
    }  
      
    data, err := c.client.GetMovieDetails(tmdbID, lang)  
    if err != nil {  
        return nil, err  
    }  
      
    // Cache por 12 horas  
    c.cache.Set(ctx, key, data, 12*time.Hour)  
    return data, nil  
}  
  
func (c *CachedTMDBClient) GetTVDetails(tmdbID int, lang string) (*TMDBTVDetails, error) {  
    ctx := context.Background()  
    key := fmt.Sprintf("tmdb:tv:%d:%s", tmdbID, lang)  
      
    var result TMDBTVDetails  
    if err := c.cache.Get(ctx, key, &result); err == nil {  
        return &result, nil  
    }  
      
    data, err := c.client.GetTVDetails(tmdbID, lang)  
    if err != nil {  
        return nil, err  
    }  
      
    // Cache por 12 horas  
    c.cache.Set(ctx, key, data, 12*time.Hour)  
    return data, nil  
}  
  
func (c *CachedTMDBClient) GetTVSeason(tmdbID int, seasonNumber int, lang string) (*TMDBSeasonDetails, error) {  
    ctx := context.Background()  
    key := fmt.Sprintf("tmdb:season:%d:%d:%s", tmdbID, seasonNumber, lang)  
      
    var result TMDBSeasonDetails  
    if err := c.cache.Get(ctx, key, &result); err == nil {  
        return &result, nil  
    }  
      
    data, err := c.client.GetTVSeason(tmdbID, seasonNumber, lang)  
    if err != nil {  
        return nil, err  
    }  
      
    // Cache por 24 horas (episódios não mudam)  
    c.cache.Set(ctx, key, data, 24*time.Hour)  
    return data, nil  
}

// GetTVSeriesByIMDB - implementa TMDBClientInterface
func (c *CachedTMDBClient) GetTVSeriesByIMDB(imdbID string, lang string) (*TMDBTVDetails, error) {
    tmdbID, mediaType, err := c.FindByIMDBID(imdbID)
    if err != nil {
        return nil, err
    }
    if mediaType != "series" {
        return nil, fmt.Errorf("ID %s não é uma série", imdbID)
    }
    return c.GetTVDetails(tmdbID, lang)
}

func (c *CachedTMDBClient) GetStreamsFromTMDB(imdbID string, mediaType string) ([]models.Stream, error) {
    ctx := context.Background()
    key := fmt.Sprintf("tmdb:streams:%s:%s", imdbID, mediaType)

    var result []models.Stream
    if err := c.cache.Get(ctx, key, &result); err == nil {
        return result, nil
    }

    data, err := c.client.GetStreamsFromTMDB(imdbID, mediaType)
    if err != nil {
        return nil, err
    }

    c.cache.Set(ctx, key, data, 1*time.Hour)
    return data, nil
}