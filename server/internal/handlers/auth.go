package handlers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"messenger/internal/models"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

const (
	accessTokenEXP  = 15 * time.Minute
	refreshTokenExp = 7 * 24 * time.Hour
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

func (h *Handler) RegistrateNewUser(c *gin.Context) {
	var req User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль слишком короткий (минимум 8 символов)"})
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	var userID string
	err = h.DB.QueryRow(context.Background(),
		"INSERT INTO users (username, email, password_hash, created_at) VALUES ($1, $2, $3, $4) RETURNING user_id",
		req.Username, req.Email, string(hashedPassword), time.Now(),
	).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save user data"})
		return
	}

	tokens, err := h.generateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токенов"})
		return
	}

	_, err = h.DB.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hashToken(tokens.RefreshToken), time.Now().Add(refreshTokenExp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении refresh токена"})
		return
	}

	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		int(refreshTokenExp.Seconds()),
		"/",
		"",
		true, // Secure
		true, // HttpOnly
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":       "success",
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверный формат запроса"})
		return
	}

	var user User
	err := h.DB.QueryRow(context.Background(),
		"SELECT user_id, email, password_hash from users WHERE email = $1",
		req.Email,
	).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Пользователь не найден"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный пароль"})
		return
	}

	tokens, err := h.generateTokens(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при генерации токена"})
		return
	}

	_, err = h.DB.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		user.ID, hashToken(tokens.RefreshToken), time.Now().Add(refreshTokenExp))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при сохранении refresh токена"})
		return
	}

	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		int(refreshTokenExp.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":       "success",
		"access_token": tokens.AccessToken,
	})
}

func (h *Handler) Logout(c *gin.Context) {
	refreshCookie, err := c.Cookie("refresh_token")
	if err == nil {
		_, _ = h.DB.Exec(context.Background(),
			`DELETE FROM refresh_tokens WHERE token_hash = $1`,
			hashToken(refreshCookie))
	}

	authHeader := c.GetHeader("Authorization")
	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
		tokenString := authHeader[7:]
		revokedTokens.Store(tokenString, true)
	}

	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged out successfully",
	})
}

func (h *Handler) RefreshTokens(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Необходим refresh токен"})
		return
	}

	var userID string
	err = h.DB.QueryRow(context.Background(),
		`DELETE FROM refresh_tokens WHERE token_hash = $1 RETURNING user_id`, hashToken(cookie)).Scan(&userID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Неверный refresh токен"})
		return
	}

	tokens, err := h.generateTokens(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка генерации токена"})
		return
	}

	_, err = h.DB.Exec(context.Background(),
		`INSERT INTO refresh_tokens (user_id, token_hash, expires_at) VALUES ($1, $2, $3)`,
		userID, hashToken(tokens.RefreshToken), time.Now().Add(refreshTokenExp))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Ошибка сохранения токена"})
		return
	}

	c.SetCookie(
		"refresh_token",
		tokens.RefreshToken,
		int(refreshTokenExp.Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusCreated, gin.H{
		"status":       "success",
		"access_token": tokens.AccessToken,
	})
}

// Служебные функции (без изменений)

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
