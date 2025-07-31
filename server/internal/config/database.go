package config

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	host         = "localhost"
	port         = 5432
	user         = "postgres"
	password     = "123"
	dbname       = "messenger"
	poolMaxConns = 10
)

func InitDB() *pgxpool.Pool {
	connStr := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname,
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
