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

func (c *CachedCinemetaClient) GetManifest() (*models.ManifestResponse, error) {
    ctx := context.Background()
    key := "cinemeta:manifest"

    var result models.ManifestResponse
    if err := c.cache.Get(ctx, key, &result); err == nil {
        return &result, nil
    }

    data, err := c.client.GetManifest()
    if err != nil {
        return nil, err
    }

    // Cache por 24 horas (manifest raramente muda)
    c.cache.Set(ctx, key, data, 24*time.Hour)
    return data, nil
}

func (c *CachedCinemetaClient) GetCatalogWithFilters(catalogType, catalogID, extraArgs string) (*models.CatalogResponse, error) {
    ctx := context.Background()
    key := fmt.Sprintf("cinemeta:catalog:%s:%s:%s", catalogType, catalogID, extraArgs)

    var result models.CatalogResponse
    if err := c.cache.Get(ctx, key, &result); err == nil {
        return &result, nil
    }

    data, err := c.client.GetCatalogWithFilters(catalogType, catalogID, extraArgs)
    if err != nil {
        return nil, err
    }

    // Cache por 1 hora (mesma duração do catálogo sem filtro)
    c.cache.Set(ctx, key, data, time.Hour)
    return data, nil
}