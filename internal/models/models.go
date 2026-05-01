package models

// MetaResponse representa a resposta completa da API para um título (meta)
type MetaResponse struct {
    Meta Meta `json:"meta"`
}

// Meta contém todos os detalhes de um filme ou série
type Meta struct {
    ID           string   `json:"id"`
    Type         string   `json:"type"`
    Name         string   `json:"name"`
    Year         string   `json:"year"`
    Poster       string   `json:"poster"`
    Background   string   `json:"background"`
    ImdbRating   string   `json:"imdbRating"`
    Description  string   `json:"description"`
    Cast         []string `json:"cast"`
    Director     []string `json:"director"`
    Genre        []string `json:"genre"`
    Runtime      string   `json:"runtime"`
    Country      string   `json:"country"`
    Awards       string   `json:"awards"`
    Trailers     []Trailer `json:"trailers"`
	Videos       []Video  `json:"videos"`
	Logo 		 string   `json:"logo"`
    // Adicione outros campos conforme necessário
}

// Trailer representa um trailer do YouTube
type Trailer struct {
    Source string `json:"source"` // ID do vídeo no YouTube
    Type   string `json:"type"`
}

// CatalogResponse representa a resposta do catálogo
type CatalogResponse struct {
    Metas []CatalogMeta `json:"metas"`
}

// CatalogMeta contém informações resumidas para exibição em listas
type CatalogMeta struct {
    ID     string   `json:"id"`
    Type   string   `json:"type"`
    Name   string   `json:"name"`
    Poster string   `json:"poster"`
    Year   string   `json:"year"`
    Genre  []string `json:"genre"`
}

// Video representa um episódio
type Video struct {
    ID        string `json:"id"`
    Name      string `json:"name"`
    Season    int    `json:"season"`
    Episode   int    `json:"episode"`
    Released  string `json:"released"`
	Thumbnail string `json:"thumbnail"`
}

// SeasonGroup agrupa episódios por temporada
type SeasonGroup struct {
    SeasonNumber int     `json:"season"`
    Episodes     []Video `json:"episodes"`
}

// ManifestResponse representa a resposta do manifest
type ManifestResponse struct {
    ID          string                 `json:"id"`
    Version     string                 `json:"version"`
    Name        string                 `json:"name"`
    Description string                 `json:"description"`
    Catalogs    []ManifestCatalog      `json:"catalogs"`
}

// ManifestCatalog representa um catálogo no manifest
type ManifestCatalog struct {
    Type  string `json:"type"`
    ID    string `json:"id"`
    Name  string `json:"name"`
    Extra map[string]interface{} `json:"extra,omitempty"`
}