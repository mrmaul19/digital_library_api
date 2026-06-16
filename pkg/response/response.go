package response

import "github.com/gofiber/fiber/v2"

type Response struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    Error   string      `json:"error,omitempty"`
}

type PaginatedResponse struct {
    Success bool        `json:"success"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
    Meta    Meta        `json:"meta"`
}

type Meta struct {
    Total      int `json:"total"`
    Page       int `json:"page"`
    Limit      int `json:"limit"`
    TotalPages int `json:"total_pages"`
}

func OK(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusOK).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func Created(c *fiber.Ctx, message string, data interface{}) error {
    return c.Status(fiber.StatusCreated).JSON(Response{
        Success: true,
        Message: message,
        Data:    data,
    })
}

func BadRequest(c *fiber.Ctx, message string) error {
    return c.Status(fiber.StatusBadRequest).JSON(Response{
        Success: false,
        Message: "Bad Request",
        Error:   message,
    })
}

func Unauthorized(c *fiber.Ctx, message string) error {
    return c.Status(fiber.StatusUnauthorized).JSON(Response{
        Success: false,
        Message: "Unauthorized",
        Error:   message,
    })
}

func NotFound(c *fiber.Ctx, message string) error {
    return c.Status(fiber.StatusNotFound).JSON(Response{
        Success: false,
        Message: "Not Found",
        Error:   message,
    })
}

func InternalError(c *fiber.Ctx, message string) error {
    return c.Status(fiber.StatusInternalServerError).JSON(Response{
        Success: false,
        Message: "Internal Server Error",
        Error:   message,
    })
}

func Paginated(c *fiber.Ctx, message string, data interface{}, meta Meta) error {
    return c.Status(fiber.StatusOK).JSON(PaginatedResponse{
        Success: true,
        Message: message,
        Data:    data,
        Meta:    meta,
    })
}