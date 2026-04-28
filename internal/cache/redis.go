package cache  
  
import (  
    "context"  
    "encoding/json"  
    "time"  
      
    "github.com/redis/go-redis/v9"  
)  
  
type RedisCache struct {  
    client *redis.Client  
}  
  
func NewRedisCache(addr string) *RedisCache {  
    rdb := redis.NewClient(&redis.Options{  
        Addr:     addr,  
        Password: "", // sem senha por padrão  
        DB:       0,  // DB padrão  
    })  
      
    return &RedisCache{client: rdb}  
}  
  
func (c *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {  
    val, err := c.client.Get(ctx, key).Result()  
    if err != nil {  
        return err  
    }  
    return json.Unmarshal([]byte(val), dest)  
}  
  
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {  
    data, err := json.Marshal(value)  
    if err != nil {  
        return err  
    }  
    return c.client.Set(ctx, key, data, ttl).Err()  
}  
  
func (c *RedisCache) Delete(ctx context.Context, key string) error {  
    return c.client.Del(ctx, key).Err()  
}