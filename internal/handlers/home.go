package handlers

import (
    "net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    // Redireciona para o catálogo de filmes populares
    http.Redirect(w, r, "/catalog/movie/top", http.StatusSeeOther)
}