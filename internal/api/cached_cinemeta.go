// internal/api/cached_cinemeta.go  
package api  
  
import (  
    "context"  
	"fmt"
    "time"  
      
    "topstrem/internal/cache"  
    "topstrem/internal/models"  
)  
  
type CachedCinemetaClient struct {  
    client *Client  
    cache  *cache.RedisCache  
}  
  
func NewCachedCinemetaClient(client *Client, cache *cache.RedisCache) *CachedCinemetaClient {  
    return &CachedCinemetaClient{  
        client: client,  
        cache:  cache,  
    }  
}  
  
func (c *CachedCinemetaClient) GetCatalog(catalogType, catalogID string) (*models.CatalogResponse, error) {  
    ctx := context.Background()  
    key := fmt.Sprintf("cinemeta:catalog:%s:%s", catalogType, catalogID)  
      
    var result models.CatalogResponse  
    if err := c.cache.Get(ctx, key, &result); err == nil {  
        return &result, nil  
    }  
      
    data, err := c.client.GetCatalog(catalogType, catalogID)  
    if err != nil {  
        return nil, err  
    }  
      
    // Cache por 1 hora  
    c.cache.Set(ctx, key, data, time.Hour)  
    return data, nil  
}  
  
func (c *CachedCinemetaClient) GetMeta(mediaType, id string) (*models.Meta, error) {  
    ctx := context.Background()  
    key := fmt.Sprintf("cinemeta:meta:%s:%s", mediaType, id)  
      
    var result models.Meta  
    if err := c.cache.Get(ctx, key, &result); err == nil {  
        return &result, nil  
    }  
      
    data, err := c.client.GetMeta(mediaType, id)  
    if err != nil {  
        return nil, err  
    }  
      
    // Cache por 6 horas  
    c.cache.Set(ctx, key, data, 6*time.Hour)  
    return data, nil  
}