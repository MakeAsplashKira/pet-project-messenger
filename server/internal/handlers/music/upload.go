package music

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/faiface/beep/mp3"
)

func (h *MusicHandler) UploadTrackHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseMultipartForm(50 << 20) // увеличил лимит, если треки могут быть крупнее
	if err != nil {
		http.Error(w, "Could not parse multipart form", http.StatusBadRequest)
		return
	}

	// --- Читаем трек в память ---
	file, header, err := r.FormFile("track_file")
	if err != nil {
		http.Error(w, "Track file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	var trackBuf bytes.Buffer
	if _, err := io.Copy(&trackBuf, file); err != nil {
		http.Error(w, "Failed to read track file", http.StatusInternalServerError)
		return
	}

	// --- Читаем обложку в память (если есть) ---
	var coverBuf bytes.Buffer
	var coverFileName string
	coverFile, coverHeader, err := r.FormFile("cover_image")
	if err == nil {
		defer coverFile.Close()
		if _, err := io.Copy(&coverBuf, coverFile); err != nil {
			http.Error(w, "Failed to read cover file", http.StatusInternalServerError)
			return
		}
		coverFileName = fmt.Sprintf("uploaded_covers/%d_%s", time.Now().UnixNano(), coverHeader.Filename)
	}

	title := r.FormValue("title")
	albumID := r.FormValue("album_id")
	artistID := r.FormValue("artist_id")
	uploaderID := r.FormValue("uploader_id")

	// --- Записываем трек временно в память (чтобы получить длительность) ---
	tmpTrackPath := fmt.Sprintf("tmp/%d_%s", time.Now().UnixNano(), header.Filename)
	err = os.WriteFile(tmpTrackPath, trackBuf.Bytes(), 0644)
	if err != nil {
		http.Error(w, "Failed to write temp track file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpTrackPath) // Удалим временный файл после

	duration, err := GetMP3Duration(tmpTrackPath)
	if err != nil {
		http.Error(w, "Failed to get track duration: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// --- Вставляем в базу, генерируем имена файлов для окончательного хранения ---
	trackFileName := fmt.Sprintf("uploaded_tracks/%d_%s", time.Now().UnixNano(), header.Filename)

	_, err = h.DB.Exec(context.Background(),
		`INSERT INTO tracks (album_id, artist_id, title, file_path, file_size, plays_count, uploader_id, created_at, cover_path, duration)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
		albumID, artistID, title, trackFileName, header.Size, 0, uploaderID, time.Now(), coverFileName, duration)
	if err != nil {
		http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// --- После успешной вставки записываем файлы на диск ---
	err = os.WriteFile(trackFileName, trackBuf.Bytes(), 0644)
	if err != nil {
		http.Error(w, "Failed to save track file", http.StatusInternalServerError)
		return
	}

	if coverFileName != "" {
		err = os.WriteFile(coverFileName, coverBuf.Bytes(), 0644)
		if err != nil {
			http.Error(w, "Failed to save cover file", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Track uploaded successfully")
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
