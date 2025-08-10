package handlers

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"fmt"
	"net/smtp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type EmailHandler struct {
	*Handler
	smtpConfig SMTPConfig
}

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

func NewEmailHandler(h *Handler, smtpConfig SMTPConfig) *EmailHandler {
	return &EmailHandler{
		Handler:    h,
		smtpConfig: smtpConfig,
	}
}

type EmailCheckRequest struct {
	Email string `json:"email"`
}

type VerificationRequest struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

type VerificationCode struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (h *EmailHandler) CheckEmailAvailability(c *gin.Context) {
	var req EmailCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	var exists bool
	err := h.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		strings.ToLower(req.Email),
	).Scan(&exists)

	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка базы данных"})
		return
	}

	c.JSON(200, gin.H{
		"available": !exists,
	})
}

func generateRandomCode() (string, error) {
	var n uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n%1000000), nil
}

func (h *EmailHandler) SendVerificationCodeHandler(c *gin.Context) {
	var req EmailCheckRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	code, err := generateRandomCode()
	if err != nil {
		c.JSON(500, gin.H{"error": "Не удалось сгенерировать код"})
		return
	}

	expiresAt := time.Now().Add(15 * time.Minute)
	createdAt := time.Now()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err = h.DB.Exec(ctx,
		"DELETE FROM verification_codes WHERE email = $1",
		req.Email)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при удалении старых кодов"})
		return
	}

	_, err = h.DB.Exec(ctx,
		"INSERT INTO verification_codes (email, code, created_at, expires_at) VALUES ($1, $2, $3, $4)",
		req.Email, code, createdAt, expiresAt,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при сохранении кода"})
		return
	}

	// subject := "Код подтверждения для MediaVerse"
	// body := fmt.Sprintf(`
	// Здравствуйте!

	// Ваш код подтверждения: %s
	// Действителен до: %s

	// Никому не сообщайте этот код!`,
	// 	code, expiresAt.Format("15:04 02.01.2006"))

	// if err := h.sendEmail(req.Email, subject, body); err != nil {
	// 	log.Println("Ошибка отправки E-mail: ", err)
	// 	c.JSON(500, gin.H{"error": "Не удалось отправить код"})
	// 	return
	// }

	fmt.Printf("Вот твой код бро: %v для E-mail: %v\n", code, req.Email)

	c.JSON(200, gin.H{
		"status":  "success",
		"message": "Код отправлен на email",
	})
}

func (h *EmailHandler) sendEmail(to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		h.smtpConfig.From, to, subject, body,
	)

	tlsConfig := &tls.Config{
		ServerName: h.smtpConfig.Host,
	}

	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", h.smtpConfig.Host, h.smtpConfig.Port), tlsConfig)
	if err != nil {
		return fmt.Errorf("ошибка подключения: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, h.smtpConfig.Host)
	if err != nil {
		return fmt.Errorf("ошибка SMTP-клиента: %v", err)
	}
	defer client.Close()

	if err = client.Auth(smtp.PlainAuth(
		"",
		h.smtpConfig.Username,
		h.smtpConfig.Password,
		h.smtpConfig.Host,
	)); err != nil {
		return fmt.Errorf("ошибка аутентификации: %v", err)
	}

	if err = client.Mail(h.smtpConfig.Username); err != nil {
		return fmt.Errorf("ошибка отправителя: %v", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("ошибка получателя: %v", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("ошибка данных: %v", err)
	}
	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("ошибка записи: %v", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("ошибка закрытия: %v", err)
	}

	return client.Quit()
}

func (h *EmailHandler) VerifyCode(c *gin.Context) {
	var req VerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Неверный формат запроса"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var dbCode struct {
		Code      string
		ExpiresAt time.Time
	}

	err := h.DB.QueryRow(ctx,
		`SELECT code, expires_at FROM verification_codes 
		WHERE email = $1 AND expires_at > NOW() 
		ORDER BY created_at DESC LIMIT 1`,
		req.Email,
	).Scan(&dbCode.Code, &dbCode.ExpiresAt)

	if err != nil {
		c.JSON(400, gin.H{"error": "Код не найден или истёк"})
		return
	}

	if dbCode.Code != req.Code {
		c.JSON(400, gin.H{"error": "Неверный код"})
		return
	}

	_, err = h.DB.Exec(ctx,
		"UPDATE verification_codes SET is_used = true WHERE email = $1 AND code = $2",
		req.Email, req.Code,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Ошибка при подтверждении email"})
		return
	}

	c.JSON(200, gin.H{"status": "success"})
}
