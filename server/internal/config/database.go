package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDB() *pgxpool.Pool {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	poolMaxConns := 10

	portInt, err := strconv.Atoi(port)
	if err != nil {
		log.Fatal("Неверный формат порта: ", err)
	}

	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, portInt, dbname,
	)

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatal("Ошибка парсинка конфига:", err)
	}
	config.MaxConns = int32(poolMaxConns)
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30

	db, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal("Ошибка подключения к БД: ", err)
	}
	if err := db.Ping(context.Background()); err != nil {
		log.Fatal("Ошибка ping к БД: ", err)
	}
	log.Println("Успешное подключение к PostgreSQL")
	return db
}
