package middleware

import (
    "strings"

    "github.com/gofiber/fiber/v2"
    "github.com/mrmaul19/digital-library/internal/domain"
    jwtpkg "github.com/mrmaul19/digital-library/pkg/jwt"
    "github.com/mrmaul19/digital-library/pkg/response"
)

type AuthMiddleware struct {
    jwtService *jwtpkg.JWTService
}

func NewAuthMiddleware(jwtService *jwtpkg.JWTService) *AuthMiddleware {
    return &AuthMiddleware{jwtService: jwtService}
}

// Protect — wajib login
func (m *AuthMiddleware) Protect() fiber.Handler {
    return func(c *fiber.Ctx) error {
        authHeader := c.Get("Authorization")
        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            return response.Unauthorized(c, "token tidak ditemukan")
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := m.jwtService.ValidateToken(tokenStr)
        if err != nil {
            return response.Unauthorized(c, "token tidak valid atau sudah expired")
        }

        // Simpan claims ke context agar bisa dipakai di handler
        c.Locals("userID", claims.UserID)
        c.Locals("userEmail", claims.Email)
        c.Locals("userRole", claims.Role)

        return c.Next()
    }
}

// AdminOnly — hanya untuk admin
func (m *AuthMiddleware) AdminOnly() fiber.Handler {
    return func(c *fiber.Ctx) error {
        role, ok := c.Locals("userRole").(domain.UserRole)
        if !ok || role != domain.RoleAdmin {
            return response.Unauthorized(c, "akses ditolak, hanya untuk admin")
        }
        return c.Next()
    }
}