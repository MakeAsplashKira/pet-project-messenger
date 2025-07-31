package main

import (
	"log"
	"messenger/internal/config"
	"messenger/internal/handlers"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load("data.env"); err != nil {
		log.Fatal("Не удалось загрузить data.env: ", err)
	}
	db := config.InitDB()
	defer db.Close()

	// Роутинг
	h := handlers.NewHandler(db)
	mux := http.NewServeMux()

	mux.HandleFunc("/api/check-email", h.CheckEmailAvailability)
	mux.HandleFunc("/api/check-username", h.CheckUsernameAvailability)
	mux.HandleFunc("/api/send-verif-code", h.SendVerificationCodeHandler)
	mux.HandleFunc("/api/check-verif-code", h.VerifyCode)
	mux.HandleFunc("/api/reg-user", h.RegistrateNewUser)

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
