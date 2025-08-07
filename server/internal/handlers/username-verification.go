package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Username struct {
	Username string `json:"username"`
}

func (h *Handler) CheckUsernameAvailability(c *gin.Context) {
	var req Username
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	var exists bool
	err := h.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
		req.Username,
	).Scan(&exists)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка базы данных"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"available": !exists,
	})
}
