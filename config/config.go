package config

import (
    "log"
    "os"
    "strconv"

    "github.com/joho/godotenv"
)

type Config struct {
    App      AppConfig
    Database DatabaseConfig
    JWT      JWTConfig
    Upload   UploadConfig
}

type AppConfig struct {
    Port string
    Env  string
}

type DatabaseConfig struct {
    Host     string
    Port     string
    User     string
    Password string
    Name     string
}

type JWTConfig struct {
    Secret      string
    ExpireHours int
}

type UploadConfig struct {
    Dir         string
    MaxFileSize int64
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, reading from environment")
    }

    expireHours, _ := strconv.Atoi(getEnv("JWT_EXPIRE_HOURS", "24"))
    maxSize, _ := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "5242880"), 10, 64)

    return &Config{
        App: AppConfig{
            Port: getEnv("APP_PORT", "3000"),
            Env:  getEnv("APP_ENV", "development"),
        },
        Database: DatabaseConfig{
            Host:     getEnv("DB_HOST", "localhost"),
            Port:     getEnv("DB_PORT", "5432"),
            User:     getEnv("DB_USER", "postgres"),
            Password: getEnv("DB_PASSWORD", ""),
            Name:     getEnv("DB_NAME", "digital_library"),
        },
        JWT: JWTConfig{
            Secret:      getEnv("JWT_SECRET", "secret"),
            ExpireHours: expireHours,
        },
        Upload: UploadConfig{
            Dir:         getEnv("UPLOAD_DIR", "./uploads"),
            MaxFileSize: maxSize,
        },
    }
}

func getEnv(key, defaultVal string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return defaultVal
}