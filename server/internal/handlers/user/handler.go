package user

import "messenger/internal/handlers"

type UserHandler struct {
	*handlers.Handler
}

func NewUserHandler(h *handlers.Handler) *UserHandler {
	return &UserHandler{
		Handler: h,
	}
}
