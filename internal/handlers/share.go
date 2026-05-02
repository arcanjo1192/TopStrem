package handlers

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "topstrem/internal/auth"
    "topstrem/internal/crypto"
)

func ShareTokenHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        email, err := auth.GetEmailFromRequest(c.Request)
        if err != nil {
            c.JSON(http.StatusUnauthorized, gin.H{"error": "não autenticado"})
            return
        }
        token, err := crypto.Encrypt(email)
        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "erro ao gerar token"})
            return
        }
        c.JSON(http.StatusOK, gin.H{"token": token})
    }
}