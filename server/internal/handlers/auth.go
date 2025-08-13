package handlers

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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
	ID            string `json:"id"`
	Username      string `json:"username"`
	Email         string `json:"email"`
	Password      string `json:"password"`
	ImageOriginal string `json:"image_original"`
	ImageAvatar   string `json:"image_avatar"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Country       string `json:"country"`
	City          string `json:"city"`
	Gender        string `json:"gender"`
}

func (h *Handler) RegistrateNewUser(c *gin.Context) {
	// 1. Парсим данные (не файлы)
	var req User
	req.Username = c.PostForm("username")
	req.Email = c.PostForm("email")
	req.Password = c.PostForm("password")
	req.FirstName = c.PostForm("first_name")
	req.LastName = c.PostForm("last_name")
	req.Country = c.PostForm("country")
	req.City = c.PostForm("city")
	req.Gender = c.PostForm("gender")

	if req.Username == "" || req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Необходимо заполнить все обязательные поля"})
		return
	}

	// 2. Получаем файлы формы
	originalFile, originalHeader, err := c.Request.FormFile("image_original")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нету файла."})
		return
	}
	defer originalFile.Close()

	avatarFile, avatarHeader, err := c.Request.FormFile("image_avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Нету файла."})
		return
	}
	defer avatarFile.Close()

	// 3. Создаем папку для загрузке, если её нету.

	uploadDir := "uploads/avatars/"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Не удалось создать директорию"})
		return
	}

	// 4. Генерируем уникальный имена файлов
	originalExt := filepath.Ext(originalHeader.Filename)
	avatarExt := filepath.Ext(avatarHeader.Filename)

	originalName := fmt.Sprintf("original_%d%s", time.Now().UnixNano(), originalExt)
	avatarName := fmt.Sprintf("avatar_%d%s", time.Now().UnixNano(), avatarExt)

	originalPath := filepath.Join(uploadDir, originalName)
	avatarPath := filepath.Join(uploadDir, avatarName)

	// 5. Сохраняем файлы
	if err := c.SaveUploadedFile(originalHeader, originalPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения original файла"})
		return
	}
	if err := c.SaveUploadedFile(avatarHeader, avatarPath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сохранения avatar файла"})
		return
	}

	// 6. Валидация пароля
	if len(req.Password) < 8 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Пароль слишком короткий (минимум 8 символов)"})
		return
	}
	// 7. Хеширование пароля
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка при хешировании пароля"})
		return
	}

	var userID string
	err = h.DB.QueryRow(context.Background(),
		`INSERT INTO users (username, email, password_hash,
		image_original, image_avatar, first_name,
		last_name, country, city, gender, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING user_id`,
		req.Username, req.Email, string(hashedPassword), originalPath, avatarPath,
		req.FirstName, req.LastName, req.Country, req.City, req.Gender, time.Now(),
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
		"success":      true,
		"access_token": tokens.AccessToken,
		"user_id":      userID,
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
		"user_id":      user.ID,
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
