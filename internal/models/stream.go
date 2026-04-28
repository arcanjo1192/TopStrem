package models

type StreamResponse struct {
    Streams []Stream `json:"streams"`
}

type Stream struct {
    Name        string `json:"name"`
    ExternalUrl string `json:"externalUrl"`
    AndroidUrl  string `json:"androidUrl"`
    IosUrl      string `json:"iosUrl"`
}