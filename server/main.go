package main

import (
	"log"
	"messenger/internal/config"
	"messenger/internal/handlers"
	"messenger/internal/handlers/music"
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

	smtpConfig := handlers.SMTPConfig{
		Host:     "smtp.yandex.ru",
		Port:     465,
		Username: "sanyasatana@yandex.ru",
		Password: os.Getenv("YANDEX_SMTP_PASSWORD"),
		From:     "Pet-project <sanyasatana@yandex.ru>",
	}

	// Роутинг
	auth := handlers.NewHandler(db, jwtSecret)
	authEmail := handlers.NewEmailHandler(auth, smtpConfig)
	musicHandler := music.NewMusicHandler(auth)
	mux := http.NewServeMux()

	//Авторизация
	mux.HandleFunc("/api/check-email", authEmail.CheckEmailAvailability)
	mux.HandleFunc("/api/send-verif-code", authEmail.SendVerificationCodeHandler)
	mux.HandleFunc("/api/check-verif-code", authEmail.VerifyCode)
	mux.HandleFunc("/api/check-username", auth.CheckUsernameAvailability)
	mux.HandleFunc("/api/reg-user", auth.RegistrateNewUser)
	mux.HandleFunc("/api/refresh-tokens", auth.RefreshTokens)
	mux.HandleFunc("/api/login", auth.Login)
	//Музыка
	mux.HandleFunc("/api/upload-track", musicHandler.UploadTrackHandler)
	mux.HandleFunc("/api/stream", musicHandler.StreamTrack)

	protected := http.NewServeMux()
	protected.HandleFunc("/api/logout", auth.Logout)

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
