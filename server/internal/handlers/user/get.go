package user

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Profile struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Avatar    string `json:"image_avatar"`
}

func (h *UserHandler) GetUserProfileDataHandler(c *gin.Context) {
	currentUserID := c.MustGet("user_id")
	if currentUserID == nil {
		c.JSON(http.StatusNotAcceptable, gin.H{"error": "ты пидарас куда лезишь"})
	}

	requestUserID := c.Param("user_id")

	var profile Profile

	err := h.DB.QueryRow(context.Background(),
		`SELECT user_id, username, first_name, last_name, image_avatar FROM users WHERE user_id = $1`,
		requestUserID).Scan(&profile.ID, &profile.Username, &profile.FirstName, &profile.LastName, &profile.Avatar)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": profile})
}
