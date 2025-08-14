package artist

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

type Artist struct {
	Name    string `json:"name"`
	OwnerID string `json:"owner_id"`
	Avatar  string `json:"avatar"`
	Cover   string `json:"cover"`
}

func (h *ArtistHandler) UploadArtist(c *gin.Context) {
	var req Artist

	// 1. Парсим данные с формы (без файлов)
	req.Name = c.PostForm("name")
	req.OwnerID = c.PostForm("owner_id")

	// 2. Парсим файлы с формы
	coverFile, coverHeader, err := c.Request.FormFile("cover")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нету файла."})
		return
	}
	defer coverFile.Close()

	avatarFile, avatarHeader, err := c.Request.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нету файла."})
		return
	}
	defer avatarFile.Close()

	// 3. Создаем директорию загрузки
	uploadCoverDir := "uploads/artist/cover"
	uploadAvatarDir := "uploads/artist/avatar"
	if err := os.MkdirAll(uploadCoverDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать директорию для обложки"})
		return
	}
	if err := os.MkdirAll(uploadAvatarDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать директорию для аватара"})
	}
	// 4. Создаем имена файлов
	coverExt := filepath.Ext(coverHeader.Filename)
	avatarExt := filepath.Ext(avatarHeader.Filename)

	coverlName := fmt.Sprintf("cover_%d%s", time.Now().UnixNano(), coverExt)
	avatarName := fmt.Sprintf("avatar_%d%s", time.Now().UnixNano(), avatarExt)

	coverPath := filepath.Join(uploadCoverDir, coverlName)
	avatarPath := filepath.Join(uploadAvatarDir, avatarName)

	// 5. Сохраняем файлы
	if err := c.SaveUploadedFile(coverHeader, coverPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения обложки"})
		return
	}
	if err := c.SaveUploadedFile(avatarHeader, avatarPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения аватара"})
		return
	}
	//6. Сохраняем все в базу данных
	_, err = h.DB.Exec(context.Background(),
		`INSERT INTO artists (name, avatar_url, cover_url, owner_id, created_at) VALUES ($1, $2, $3, $4, $5)`,
		req.Name, avatarPath, coverPath, req.OwnerID, time.Now())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения данных в БД"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
	})
}
