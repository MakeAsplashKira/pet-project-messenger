package music

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/faiface/beep/mp3"
	"github.com/gin-gonic/gin"
)

func (h *MusicHandler) UploadTrackHandler(c *gin.Context) {
	// --- Получаем трек ---
	trackHeader, err := c.FormFile("track_file")
	if err != nil {
		c.JSON(400, gin.H{"error": "Трек обязателен для загрузки"})
		return
	}
	trackFile, err := trackHeader.Open()
	if err != nil {
		c.JSON(500, gin.H{"error": "Не удалось открыть трек-файл"})
		return
	}
	defer trackFile.Close()

	var trackBuf bytes.Buffer
	if _, err := io.Copy(&trackBuf, trackFile); err != nil {
		c.JSON(500, gin.H{"error": "Не удалось прочитать трек"})
		return
	}

	// --- Получаем обложку (необязательно) ---
	var coverBuf bytes.Buffer
	var coverFileName string
	coverHeader, err := c.FormFile("cover_image")
	if err == nil {
		coverFile, err := coverHeader.Open()
		if err != nil {
			c.JSON(500, gin.H{"error": "Не удалось открыть обложку"})
			return
		}
		defer coverFile.Close()

		if _, err := io.Copy(&coverBuf, coverFile); err != nil {
			c.JSON(500, gin.H{"error": "Не удалось прочитать обложку"})
			return
		}
		coverFileName = fmt.Sprintf("uploaded_covers/%d_%s", time.Now().UnixNano(), coverHeader.Filename)
	}

	// --- Метаданные ---
	title := c.PostForm("title")
	albumID := c.PostForm("album_id")
	artistID := c.PostForm("artist_id")
	uploaderID := c.PostForm("uploader_id")

	// --- Временное сохранение трека ---
	tmpTrackPath := fmt.Sprintf("tmp/%d_%s", time.Now().UnixNano(), trackHeader.Filename)
	if err := os.WriteFile(tmpTrackPath, trackBuf.Bytes(), 0644); err != nil {
		c.JSON(500, gin.H{"error": "Не удалось записать временный файл"})
		return
	}
	defer os.Remove(tmpTrackPath)

	// --- Длительность трека ---
	duration, err := GetMP3Duration(tmpTrackPath)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при получении длительности: " + err.Error()})
		return
	}

	// --- Генерация имени и запись в базу ---
	trackFileName := fmt.Sprintf("uploaded_tracks/%d_%s", time.Now().UnixNano(), trackHeader.Filename)

	_, err = h.DB.Exec(context.Background(),
		`INSERT INTO tracks 
		(album_id, artist_id, title, file_path, file_size, plays_count, uploader_id, created_at, cover_path, duration)
		VALUES ($1, $2, $3, $4, $5, 0, $6, $7, $8, $9)`,
		albumID, artistID, title, trackFileName, trackHeader.Size, uploaderID, time.Now(), coverFileName, duration)

	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка базы данных: " + err.Error()})
		return
	}

	// --- Сохраняем файлы на диск ---
	if err := os.WriteFile(trackFileName, trackBuf.Bytes(), 0644); err != nil {
		c.JSON(500, gin.H{"error": "Не удалось сохранить трек"})
		return
	}

	if coverFileName != "" {
		if err := os.WriteFile(coverFileName, coverBuf.Bytes(), 0644); err != nil {
			c.JSON(500, gin.H{"error": "Не удалось сохранить обложку"})
			return
		}
	}

	c.JSON(201, gin.H{
		"status":  "success",
		"message": "Трек успешно загружен",
	})
}

func GetMP3Duration(filePath string) (int, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	if err != nil {
		return 0, err
	}
	defer streamer.Close()

	duration := time.Duration(streamer.Len()) * time.Second / time.Duration(format.SampleRate)
	return int(duration.Seconds()), nil
}
