package main
// @title           Digital Library API
// @version         1.0
// @description     REST API untuk Digital Library dengan fitur autentikasi JWT dan manajemen buku.

// @contact.name    Your Name
// @contact.email   your@email.com

// @host            localhost:3000
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Masukkan token dengan format: Bearer {token}

import (
    "log"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/cors"
    fiberLogger "github.com/gofiber/fiber/v2/middleware/logger"

    "github.com/mrmaul19/digital-library/config"
    "github.com/mrmaul19/digital-library/internal/delivery/http"
    "github.com/mrmaul19/digital-library/internal/delivery/http/handler"
    "github.com/mrmaul19/digital-library/internal/delivery/http/middleware"
    "github.com/mrmaul19/digital-library/internal/repository"
    "github.com/mrmaul19/digital-library/internal/usecase"
    "github.com/mrmaul19/digital-library/pkg/database"
    jwtpkg "github.com/mrmaul19/digital-library/pkg/jwt"
    "github.com/mrmaul19/digital-library/pkg/upload"
)

func main() {
    // 1. Load config
    cfg := config.Load()

    // 2. Connect database
    db := database.NewPostgres(&cfg.Database)
    defer db.Close()

    // 3. Init dependencies (Dependency Injection manual)
    jwtService  := jwtpkg.NewJWTService(cfg.JWT.Secret, cfg.JWT.ExpireHours)
    uploader    := upload.NewUploader(cfg.Upload.Dir, cfg.Upload.MaxFileSize)

    // Repository
    userRepo := repository.NewUserRepository(db)
    bookRepo := repository.NewBookRepository(db)

    // Usecase
    authUsecase := usecase.NewAuthUsecase(userRepo, jwtService)
    bookUsecase := usecase.NewBookUsecase(bookRepo)

    // Handler
    authHandler := handler.NewAuthHandler(authUsecase)
    bookHandler := handler.NewBookHandler(bookUsecase, uploader)

    // Middleware
    authMiddleware := middleware.NewAuthMiddleware(jwtService)

    // 4. Setup Fiber
    app := fiber.New(fiber.Config{
        AppName: "Digital Library API v1",
    })

    app.Use(fiberLogger.New())
    app.Use(cors.New(cors.Config{
        AllowOrigins: "*",
        AllowHeaders: "Origin, Content-Type, Accept, Authorization",
        AllowMethods: "GET, POST, PUT, DELETE",
    }))

    // Serve static files untuk uploads
    app.Static("/uploads", cfg.Upload.Dir)

    // 5. Setup routes
    router := http.NewRouter(app, authMiddleware, authHandler, bookHandler)
    router.Setup()

    // 6. Start server
    log.Printf("🚀 Server running on port %s", cfg.App.Port)
    log.Fatal(app.Listen(":" + cfg.App.Port))
}