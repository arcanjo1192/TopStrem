package handlers

import (
    "encoding/xml"
    "net/http"

    "github.com/gin-gonic/gin"
    "topstrem/internal/api"
)

// URLSet represents the sitemap XML structure
type URLSet struct {
    XMLName xml.Name `xml:"urlset"`
    XMLNS   string   `xml:"xmlns,attr"`
    URLs    []URL    `xml:"url"`
}

// URL represents a single URL entry in the sitemap
type URL struct {
    Loc        string `xml:"loc"`
    LastMod    string `xml:"lastmod,omitempty"`
    ChangeFreq string `xml:"changefreq,omitempty"`
    Priority   string `xml:"priority,omitempty"`
}

func SitemapHandler(apiClient api.CinemetaClient) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Collect all URLs
        var urls []URL

        // Add static pages
        scheme := "http"
        if c.Request.TLS != nil || c.Request.Header.Get("X-Forwarded-Proto") == "https" {
            scheme = "https"
        }
        baseURL := scheme + "://" + c.Request.Host
        urls = append(urls, URL{
            Loc:        baseURL + "/",
            ChangeFreq: "daily",
            Priority:   "1.0",
        })
        urls = append(urls, URL{
            Loc:        baseURL + "/catalog/movie/top",
            ChangeFreq: "daily",
            Priority:   "0.8",
        })
        urls = append(urls, URL{
            Loc:        baseURL + "/catalog/series/top",
            ChangeFreq: "daily",
            Priority:   "0.8",
        })

        // Fetch movie catalog
        movieCatalog, err := apiClient.GetCatalog("movie", "top")
        if err == nil && movieCatalog != nil {
            for _, meta := range movieCatalog.Metas {
                urls = append(urls, URL{
                    Loc:        baseURL + "/detail/" + meta.Type + "/" + meta.ID,
                    ChangeFreq: "weekly",
                    Priority:   "0.6",
                })
            }
        }

        // Fetch series catalog
        seriesCatalog, err := apiClient.GetCatalog("series", "top")
        if err == nil && seriesCatalog != nil {
            for _, meta := range seriesCatalog.Metas {
                urls = append(urls, URL{
                    Loc:        baseURL + "/detail/" + meta.Type + "/" + meta.ID,
                    ChangeFreq: "weekly",
                    Priority:   "0.6",
                })
            }
        }

        // Create URLSet
        urlSet := URLSet{
            XMLNS: "http://www.sitemaps.org/schemas/sitemap/0.9",
            URLs:  urls,
        }

        // Set headers
        c.Header("Content-Type", "application/xml")
        c.Header("Cache-Control", "public, max-age=86400") // Cache for 1 day

        // Render XML
        c.XML(http.StatusOK, urlSet)
    }
}