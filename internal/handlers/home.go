package handlers

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

func HomeHandler(c *gin.Context) {
    // Redireciona para o catálogo de filmes populares
    c.Redirect(http.StatusSeeOther, "/catalog/movie/top")
}