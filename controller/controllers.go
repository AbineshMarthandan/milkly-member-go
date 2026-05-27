package controller

import (
	"log"
	"milkly-member/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/gofiber/swagger"
	"go.uber.org/dig"
)

type Controllers struct {
	AuthController   *AuthController
	Config           *config.Config
}

// ControllerDeps holds all controller dependencies
type ControllerDeps struct {
	dig.In
	AuthController   *AuthController   `optional:"true"`
	Config           *config.Config
}

func NewControllers(deps ControllerDeps) *Controllers {
	return &Controllers{
		AuthController:   deps.AuthController,
		Config:           deps.Config,
	}
}

func (c *Controllers) Listen() {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	// Health check endpoint
	app.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.JSON(fiber.Map{
			"status":  "healthy",
			"service": "milkly-member",
		})
	})

	// Swagger endpoint
	app.Get("/swagger/*", fiberSwagger.HandlerDefault)

	// API routes
	api := app.Group("/api")
	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/login", c.AuthController.Login)
	auth.Post("/register", c.AuthController.Register)
	// Start server
	log.Fatal(app.Listen(":" + c.Config.Port))
}