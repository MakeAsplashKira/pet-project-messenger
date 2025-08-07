package handlers

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"
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

func (h *EmailHandler) CheckEmailAvailability(w http.ResponseWriter, r *http.Request) {
	var req EmailCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}
	var exists bool
	err := h.DB.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)",
		strings.ToLower(req.Email),
	).Scan(&exists)

	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]bool{
		"available": !exists,
	})
}

func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func generateRandomCode() (string, error) {
	var n uint32
	err := binary.Read(rand.Reader, binary.BigEndian, &n)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n%1000000), nil
}

func (h *EmailHandler) SendVerificationCodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EmailCheckRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	code, err := generateRandomCode()
	if err != nil {
		http.Error(w, "Failed to generate verification code", http.StatusInternalServerError)
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
		http.Error(w, "Failed to delete old codes", http.StatusInternalServerError)
		return
	}
	_, err = h.DB.Exec(ctx,
		"INSERT INTO verification_codes (email, code, created_at, expires_at) VALUES ($1, $2, $3, $4)",
		req.Email, code, createdAt, expiresAt,
	)
	if err != nil {
		http.Error(w, "Failed to save verification code", http.StatusInternalServerError)
		return
	}

	subject := "Код подтверждения для MediaVerse"
	body := fmt.Sprintf(`
	    Здравствуйте!

	    Ваш код подтверждения: %s
	    Действителен до: %s

	    Никому не сообщайте этот код!`,
		code, expiresAt.Format("15:04 02.01.2006"))

	if err := h.sendEmail(req.Email, subject, body); err != nil {
		log.Println("Ошибка отправки  E-mail! ", err)
		http.Error(w, "Не удалось отправить код", http.StatusInternalServerError)
		return
	}
	fmt.Printf("Вот твой код бро: %v для E-mail: %v\n", code, req.Email)

	respondJSON(w, http.StatusOK, map[string]string{
		"status":  "success",
		"message": "Код отправлен на email",
	})
}

func (h *EmailHandler) sendEmail(to, subject, body string) error {
	// 1. Формируем письмо
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=UTF-8\r\n\r\n%s",
		h.smtpConfig.From, to, subject, body,
	)

	// 2. Настройка TLS
	tlsConfig := &tls.Config{
		ServerName: h.smtpConfig.Host,
	}

	// 3. Подключение
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", h.smtpConfig.Host, h.smtpConfig.Port), tlsConfig)
	if err != nil {
		return fmt.Errorf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, h.smtpConfig.Host)
	if err != nil {
		return fmt.Errorf("Ошибка SMTP-клиента: %v.", err)
	}
	defer client.Close()

	// 4. Аутентификация
	log.Printf(
		"Auth details: host=%s, user=%s, pass=%v",
		h.smtpConfig.Host,
		h.smtpConfig.Username,
		h.smtpConfig.Password != "",
	)
	if err = client.Auth(smtp.PlainAuth(
		"",
		h.smtpConfig.Username,
		h.smtpConfig.Password,
		h.smtpConfig.Host,
	)); err != nil {
		return fmt.Errorf("Ошибка аутентификации: %v.", err)
	}

	// 5. Отправка
	if err = client.Mail(h.smtpConfig.Username); err != nil {
		return fmt.Errorf("Ошибка отправителя: %v.", err)
	}
	if err = client.Rcpt(to); err != nil {
		return fmt.Errorf("Ошибка получателя: %v.", err)
	}

	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("Ошибка данных: %v.", err)
	}
	if _, err = w.Write([]byte(msg)); err != nil {
		return fmt.Errorf("Ошибка записи: %v.", err)
	}
	if err = w.Close(); err != nil {
		return fmt.Errorf("Ошибка закрытия: %v.", err)
	}

	return client.Quit()
}

func (h *EmailHandler) VerifyCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req VerificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
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
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Verification code not found or expired",
		})
		return
	}

	if dbCode.Code != req.Code {
		respondJSON(w, http.StatusBadRequest, map[string]string{
			"error": "Invalid verification code",
		})
		return
	}

	// Если код верный, помечаем email как подтвержденный
	_, err = h.DB.Exec(ctx,
		"UPDATE verification_codes SET is_used = true WHERE email = $1 AND code = $2",
		req.Email, req.Code,
	)
	if err != nil {
		http.Error(w, "Failed to verify email", http.StatusInternalServerError)
		return
	}

	respondJSON(w, http.StatusOK, map[string]string{
		"status": "success",
	})
}
