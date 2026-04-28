package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "topstrem/internal/models"
)

const baseURL = "https://v3-cinemeta.strem.io"

type Client struct {
    httpClient *http.Client
}

func NewClient() *Client {
    return &Client{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

// GetCatalog busca o catálogo de filmes ou séries
func (c *Client) GetCatalog(catalogType, catalogID string) (*models.CatalogResponse, error) {
    url := fmt.Sprintf("%s/catalog/%s/%s.json", baseURL, catalogType, catalogID)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var catalog models.CatalogResponse
    if err := json.NewDecoder(resp.Body).Decode(&catalog); err != nil {
        return nil, err
    }
    return &catalog, nil
}

// GetMeta busca os detalhes completos de um título
func (c *Client) GetMeta(mediaType, id string) (*models.Meta, error) {
    url := fmt.Sprintf("%s/meta/%s/%s.json", baseURL, mediaType, id)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    var metaResp models.MetaResponse
    if err := json.NewDecoder(resp.Body).Decode(&metaResp); err != nil {
        return nil, err
    }
    return &metaResp.Meta, nil
}