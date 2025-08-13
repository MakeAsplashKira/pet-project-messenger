package user

import (
	"context"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
)

func (h *UserHandler) ServeAvatar(c *gin.Context) {

	requestUserID := c.Param("user_id")

	var filePath string
	err := h.DB.QueryRow(context.Background(),
		`SELECT image_avatar FROM users WHERE user_id = $1`, requestUserID).Scan(&filePath)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar not found"})
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{"error": "Avatar file not found"})
	}

	c.Header("Cache-Control", "public, max-age=86400")
	c.File(filePath)
}
