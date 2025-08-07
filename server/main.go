package main

import (
	"log"
	"messenger/internal/config"
	"messenger/internal/handlers"
	"messenger/internal/handlers/music"
	"messenger/internal/middleware"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Загрузка .env
	if err := godotenv.Load("data.env"); err != nil {
		log.Fatal("Не удалось загрузить data.env: ", err)
	}
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET не найден в .env")
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

	// Инициализация обработчиков
	auth := handlers.NewHandler(db, jwtSecret)
	authEmail := handlers.NewEmailHandler(auth, smtpConfig)
	musicHandler := music.NewMusicHandler(auth)

	// Инициализация Gin
	router := gin.Default()

	// CORS Middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Публичные маршруты
	api := router.Group("/api")
	{
		api.POST("/check-email", authEmail.CheckEmailAvailability)
		api.POST("/send-verif-code", authEmail.SendVerificationCodeHandler)
		api.POST("/check-verif-code", authEmail.VerifyCode)
		api.POST("/check-username", auth.CheckUsernameAvailability)
		api.POST("/reg-user", auth.RegistrateNewUser)
		api.POST("/refresh-tokens", auth.RefreshTokens)
		api.POST("/login", auth.Login)

		api.POST("/upload-track", musicHandler.UploadTrackHandler)
		api.GET("/stream-track/:id", musicHandler.StreamTrackHandler)
	}

	// Защищённые маршруты
	protected := api.Group("/")
	protected.Use(middleware.AuthMiddleware(jwtSecret))
	{
		protected.POST("/logout", auth.Logout)
	}

	// Запуск сервера
	log.Println("Сервер запущен на :8080")
	log.Fatal(router.Run(":8080"))
}
