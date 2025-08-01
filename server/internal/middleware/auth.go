package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(jwtSecret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Пропускаем OPTIONS запросы (для CORS)
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		// 2. Проверяем наличие заголовка
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// 3. Проверяем формат Bearer токена
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}
		tokenString := authHeader[7:]

		// 4. Парсим токен с проверкой подписи
		claims := &jwt.RegisteredClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
			// Проверяем алгоритм подписи
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(jwtSecret), nil
		})

		// 5. Проверяем ошибки парсинга
		if err != nil {
			http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// 6. Проверяем валидность токена
		if !token.Valid {
			http.Error(w, "Token is invalid", http.StatusUnauthorized)
			return
		}

		// 7. Проверяем срок действия (если указан)
		if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
			http.Error(w, "Token expired", http.StatusUnauthorized)
			return
		}

		// 8. Добавляем user_id в контекст
		ctx := context.WithValue(r.Context(), "user_id", claims.Subject)

		// 9. Передаем запрос дальше
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
