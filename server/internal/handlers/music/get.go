package music

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h *MusicHandler) StreamTrackHandler(c *gin.Context) {
	trackID := c.Param("id")
	if trackID == "" {
		c.String(http.StatusBadRequest, "Track ID is required")
		return
	}

	var filePath string
	err := h.DB.QueryRow(context.Background(),
		`SELECT file_path FROM tracks WHERE track_id = $1`, trackID).Scan(&filePath)
	if err != nil {
		c.String(http.StatusNotFound, "Track not found")
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		c.String(http.StatusNotFound, "File not found")
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		c.String(http.StatusInternalServerError, "File error")
		return
	}

	fileSize := fileInfo.Size()
	contentType := "audio/mpeg"
	rangeHeader := c.GetHeader("Range")

	var start, end int64

	if rangeHeader == "" {
		// Отдаем только первые 512 КБ
		start = 0
		const defaultChunkSize int64 = 512 * 1024
		end = start + defaultChunkSize - 1
		if end >= fileSize {
			end = fileSize - 1
		}
	} else {
		// Парсим заголовок Range
		ranges, err := parseRange(rangeHeader, fileSize)
		if err != nil || len(ranges) == 0 {
			c.String(http.StatusBadRequest, "Invalid Range header")
			return
		}
		if len(ranges) > 1 {
			c.String(http.StatusRequestedRangeNotSatisfiable, "Multiple ranges not supported")
			return
		}
		start = ranges[0].start
		end = ranges[0].end
		if end >= fileSize {
			end = fileSize - 1
		}
	}

	// Позиционируемся в файле
	_, err = file.Seek(start, io.SeekStart)
	if err != nil {
		c.String(http.StatusInternalServerError, "Seek error")
		return
	}

	// Устанавливаем заголовки
	c.Header("Content-Type", contentType)
	c.Header("Content-Length", strconv.FormatInt(end-start+1, 10))
	c.Header("Content-Range", fmt.Sprintf("bytes %d-%d/%d", start, end, fileSize))
	c.Header("Accept-Ranges", "bytes")
	c.Status(http.StatusPartialContent)

	// Передаём диапазон байт
	io.CopyN(c.Writer, file, end-start+1)
}

type httpRange struct {
	start, end int64
}

func parseRange(s string, size int64) ([]httpRange, error) {
	if s == "" {
		return nil, nil
	}

	const b = "bytes="
	if !strings.HasPrefix(s, b) {
		return nil, errors.New("invalid range")
	}

	var ranges []httpRange
	for _, ra := range strings.Split(s[len(b):], ",") {
		ra = strings.TrimSpace(ra)
		if ra == "" {
			continue
		}

		i := strings.Index(ra, "-")
		if i < 0 {
			return nil, errors.New("invalid range format")
		}

		startStr, endStr := strings.TrimSpace(ra[:i]), strings.TrimSpace(ra[i+1:])
		var r httpRange

		if startStr == "" {
			// Пример: bytes=-500
			length, err := strconv.ParseInt(endStr, 10, 64)
			if err != nil {
				return nil, errors.New("invalid range end")
			}
			if length > size {
				length = size
			}
			r.start = size - length
			r.end = size - 1
		} else {
			// Пример: bytes=500-999
			start, err := strconv.ParseInt(startStr, 10, 64)
			if err != nil || start < 0 {
				return nil, errors.New("invalid range start")
			}
			r.start = start

			if endStr == "" {
				r.end = size - 1
			} else {
				end, err := strconv.ParseInt(endStr, 10, 64)
				if err != nil || end < start {
					return nil, errors.New("invalid range end")
				}
				if end >= size {
					end = size - 1
				}
				r.end = end
			}
		}

		ranges = append(ranges, r)
	}

	if len(ranges) == 0 {
		return nil, errors.New("no valid ranges found")
	}
	return ranges, nil
}
