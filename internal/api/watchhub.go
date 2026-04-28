package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "time"

    "topstrem/internal/models"
)

const watchHubBaseURL = "https://watchhub.strem.io"

type WatchHubClient struct {
    httpClient *http.Client
}

func NewWatchHubClient() *WatchHubClient {
    return &WatchHubClient{
        httpClient: &http.Client{
            Timeout: 10 * time.Second,
        },
    }
}

func (c *WatchHubClient) GetStreams(mediaType, id string) (*models.StreamResponse, error) {
    url := fmt.Sprintf("%s/stream/%s/%s.json", watchHubBaseURL, mediaType, id)
    resp, err := c.httpClient.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API retornou status %d", resp.StatusCode)
    }

    var streamResp models.StreamResponse
    if err := json.NewDecoder(resp.Body).Decode(&streamResp); err != nil {
        return nil, err
    }
    return &streamResp, nil
}