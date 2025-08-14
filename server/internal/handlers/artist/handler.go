package artist

import "messenger/internal/handlers"

type ArtistHandler struct {
	*handlers.Handler
}

func NewArtistHandler(h *handlers.Handler) *ArtistHandler {
	return &ArtistHandler{
		Handler: h,
	}
}
