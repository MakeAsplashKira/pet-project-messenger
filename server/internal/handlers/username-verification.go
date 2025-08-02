package handlers

import (
	"context"
	"encoding/json"
	"net/http"
)

type Username struct {
	Username string `json:"username"`
}

func (h *Handler) CheckUsernameAvailability(w http.ResponseWriter, r *http.Request) {
	var req Username
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	var exists bool
	err := h.db.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE username = $1)",
		req.Username,
	).Scan(&exists)

	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{
		"available": !exists,
	})
}
