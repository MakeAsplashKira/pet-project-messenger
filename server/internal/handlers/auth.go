package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"messenger/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenEXP  = 15 * time.Minute   //15 Минут
	refreshTokenExp = 7 * 24 * time.Hour //1 Неделя
)

var revokedTokens sync.Map

type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

	var userID string
	err = h.db.QueryRow(context.Background(),
		"INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING user_id",
		req.Username, req.Email, string(hashedPassword), time.Now(),
	).Scan(&userID)
	if err != nil {
		http.Error(w, "Failed to save user data: ", http.StatusInternalServerError)
		return
	}

	tokens, err := h.generateTokens(userID)
	if err != nil {
		http.Error(w, "Ошибка при генерации токенов ", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hashToken(tokens.RefreshToken), time.Now().Add(refreshTokenExp))
	if err != nil {
		http.Error(w, "Ошибка при сохранении refresh токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   int(refreshTokenExp.Seconds()),
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status":       "success",
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	var user User
	err := h.db.QueryRow(context.Background(),
		"SELECT user_id, email, password_hash from users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		http.Error(w, "Пользователь не найден: ", http.StatusUnauthorized)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		http.Error(w, "Неверный пароль", http.StatusUnauthorized)
		return
	}

	tokens, err := h.generateTokens(user.ID)
	if err != nil {
		http.Error(w, "Ошибка при генерации токена: ", http.StatusInternalServerError)

		return
	}

	_, err = h.db.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		user.ID, hashToken(tokens.RefreshToken), time.Now().Add(refreshTokenExp))
	if err != nil {
		http.Error(w, "Ошибка при сохранении refresh токена", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   int(refreshTokenExp.Seconds()),
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokens.AccessToken,
		"status":       "success",
	})
}

func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	refreshCookies, err := r.Cookie("refresh_token")
	if err == nil {
		_, _ = h.db.Exec(context.Background(),
			`DELETE FROM refresh_tokens WHERE token_hash = $1`,
			hashToken(refreshCookies.Value))
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := authHeader[7:]
		revokedTokens.Store(tokenString, true)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		MaxAge:   -1,
	})

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

func (h *Handler) generateTokens(userID string) (*models.Tokens, error) {
	var jwtSecret = os.Getenv("JWT_SECRET")
	accessClaims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenEXP)),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	access, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return nil, err
	}

	refreshClaims := jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenExp)),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refresh, err := refreshToken.SignedString([]byte(jwtSecret + "refresh_salt"))
	if err != nil {
		return nil, err
	}

	return &models.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	}, nil
}

func hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (h *Handler) RefreshTokens(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Необходим refresh токен", http.StatusUnauthorized)
		return
	}

	var userID string
	err = h.db.QueryRow(context.Background(),
		`DELETE FROM refresh_tokens WHERE token_hash = $1 RETURNING user_id`, hashToken(cookie.Value)).Scan(&userID)
	if err != nil {
		http.Error(w, "Неверный refresh токен", http.StatusUnauthorized)
		return
	}

	tokens, err := h.generateTokens(userID)
	if err != nil {
		http.Error(w, "Ошибка генерации токена", http.StatusInternalServerError)
		return
	}

	_, err = h.db.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hashToken(tokens.RefreshToken), time.Now().Add(refreshTokenExp))
	if err != nil {
		http.Error(w, "Ошибка сохранения токена", http.StatusUnauthorized)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    tokens.RefreshToken,
		HttpOnly: true,
		Secure:   true,
		Path:     "/",
		MaxAge:   int(refreshTokenExp.Seconds()),
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"access_token": tokens.AccessToken,
		"status":       "success",
	})
}
