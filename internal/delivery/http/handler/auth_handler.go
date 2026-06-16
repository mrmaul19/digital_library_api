package handler

import (
    "github.com/gofiber/fiber/v2"
    "github.com/mrmaul19/digital-library/internal/domain"
    "github.com/mrmaul19/digital-library/pkg/response"
)

type AuthHandler struct {
    authUsecase domain.AuthUsecase
}

func NewAuthHandler(authUsecase domain.AuthUsecase) *AuthHandler {
    return &AuthHandler{authUsecase: authUsecase}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
    var req domain.RegisterRequest

    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "format request tidak valid")
    }

    // Validasi manual sederhana
    if req.Name == "" || req.Email == "" || req.Password == "" {
        return response.BadRequest(c, "name, email, dan password wajib diisi")
    }
    if len(req.Password) < 8 {
        return response.BadRequest(c, "password minimal 8 karakter")
    }

    result, err := h.authUsecase.Register(&req)
    if err != nil {
        return response.BadRequest(c, err.Error())
    }

    return response.Created(c, "registrasi berhasil", result)
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
    var req domain.LoginRequest

    if err := c.BodyParser(&req); err != nil {
        return response.BadRequest(c, "format request tidak valid")
    }

    if req.Email == "" || req.Password == "" {
        return response.BadRequest(c, "email dan password wajib diisi")
    }

    result, err := h.authUsecase.Login(&req)
    if err != nil {
        return response.Unauthorized(c, err.Error())
    }

    return response.OK(c, "login berhasil", result)
}