package main

import (
	"log"
	"messenger/internal/config"
	"messenger/internal/handlers"
	"messenger/internal/middleware"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("data.env"); err != nil {
		log.Fatal("Не удалось загрузить data.env: ", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не находится в .env файле!")
	}

	db := config.InitDB()
	defer db.Close()

	// Роутинг
	h := handlers.NewHandler(db, jwtSecret)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/check-email", h.CheckEmailAvailability)
	mux.HandleFunc("/api/check-username", h.CheckUsernameAvailability)
	mux.HandleFunc("/api/send-verif-code", h.SendVerificationCodeHandler)
	mux.HandleFunc("/api/check-verif-code", h.VerifyCode)
	mux.HandleFunc("/api/reg-user", h.RegistrateNewUser)
	mux.HandleFunc("/api/refresh-tokens", h.RefreshTokens)
	mux.HandleFunc("/api/login", h.Login)
	mux.HandleFunc("/api/upload-track", h.UploadTrackHandler)
	protected := http.NewServeMux()
	protected.HandleFunc("/api/logout", h.Logout)

	mux.Handle("/", middleware.AuthMiddleware(jwtSecret, protected))

	handler := enableCORS(mux)

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}

func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Max-Age", "3600")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
