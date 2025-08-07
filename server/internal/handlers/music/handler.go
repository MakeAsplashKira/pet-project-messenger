package music

import "messenger/internal/handlers"

type MusicHandler struct {
	*handlers.Handler
}

func NewMusicHandler(h *handlers.Handler) *MusicHandler {
	return &MusicHandler{
		Handler: h,
	}
}
