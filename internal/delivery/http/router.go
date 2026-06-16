package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	_ "github.com/mrmaul19/digital-library/docs"
	"github.com/mrmaul19/digital-library/internal/delivery/http/handler"
	"github.com/mrmaul19/digital-library/internal/delivery/http/middleware"
)

type Router struct {
	app      *fiber.App
	authMid  *middleware.AuthMiddleware
	authHand *handler.AuthHandler
	bookHand *handler.BookHandler
}

func NewRouter(
	app *fiber.App,
	authMid *middleware.AuthMiddleware,
	authHand *handler.AuthHandler,
	bookHand *handler.BookHandler,
) *Router {
	return &Router{app, authMid, authHand, bookHand}
}

func (r *Router) Setup() {
	r.app.Get("/swagger/*", swagger.HandlerDefault)

	api := r.app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok"})
	})

	auth := api.Group("/auth")
	auth.Post("/register", r.authHand.Register)
	auth.Post("/login", r.authHand.Login)

	books := api.Group("/books")
	books.Get("/", r.bookHand.GetAll)
	books.Get("/:id", r.bookHand.GetByID)
	books.Post("/", r.authMid.Protect(), r.bookHand.Create)
	books.Put("/:id", r.authMid.Protect(), r.bookHand.Update)
	books.Delete("/:id", r.authMid.Protect(), r.authMid.AdminOnly(), r.bookHand.Delete)
}