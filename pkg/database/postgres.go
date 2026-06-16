package database

import (
    "fmt"
    "log"

    "github.com/jmoiron/sqlx"
    _ "github.com/jackc/pgx/v5/stdlib"
    "github.com/mrmaul19/digital-library/config"
)

func NewPostgres(cfg *config.DatabaseConfig) *sqlx.DB {
    dsn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name,
    )

    db, err := sqlx.Connect("pgx", dsn)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(10)

    log.Println("✅ Database connected successfully")
    return db
}