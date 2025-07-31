package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"passoword"`
}

func (h *Handler) RegistrateNewUser(w http.ResponseWriter, r *http.Request) {
	var req User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	if len(req.Password) < 8 {
		http.Error(w, "Пароль слишком короткий (минимум 8 символов)", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	if err != nil {
		http.Error(w, "Ошибка при хешировании пароля: ", http.StatusInternalServerError)
		return
	}
	_, err = h.db.Exec(context.Background(),
		"INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, $4)",
		req.Username, req.Email, string(hashedPassword), time.Now(),
	)
	if err != nil {
		http.Error(w, "Failed to save user data: ", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
