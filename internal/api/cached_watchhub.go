// internal/api/cached_watchhub.go  
package api  
  
import (  
    "context"  
    "fmt"  
    "time"  
      
    "topstrem/internal/cache"  
    "topstrem/internal/models"  
)  
  
type CachedWatchHubClient struct {  
    client *WatchHubClient  
    cache  *cache.RedisCache  
}  
  
func NewCachedWatchHubClient(client *WatchHubClient, cache *cache.RedisCache) *CachedWatchHubClient {  
    return &CachedWatchHubClient{  
        client: client,  
        cache:  cache,  
    }  
}  
  
func (c *CachedWatchHubClient) GetStreams(mediaType, id string) (*models.StreamResponse, error) {  
    ctx := context.Background()  
    key := fmt.Sprintf("watchhub:streams:%s:%s", mediaType, id)  
      
    var result models.StreamResponse  
    if err := c.cache.Get(ctx, key, &result); err == nil {  
        return &result, nil  
    }  
      
    data, err := c.client.GetStreams(mediaType, id)  
    if err != nil {  
        return nil, err  
    }  
      
    // Cache por 30 minutos (links de streaming mudam frequentemente)  
    c.cache.Set(ctx, key, data, 30*time.Minute)  
    return data, nil  
}