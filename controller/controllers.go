package controller

import (
	"log"
	"milkly-member/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	fiberSwagger "github.com/gofiber/swagger"
)

type Controllers struct {
	MemberController *MemberController
	Config          *config.Config
}

func NewControllers(memberController *MemberController, cfg *config.Config) *Controllers {
	return &Controllers{
		MemberController: memberController,
		Config:          cfg,
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

	// Member routes
	member := api.Group("/member")
	member.Post("/", c.MemberController.CreateMember)
	member.Get("/", c.MemberController.ListMembers)
	member.Get("/:id", c.MemberController.GetMember)
	member.Put("/:id", c.MemberController.UpdateMember)
	member.Delete("/:id", c.MemberController.DeleteMember)

	// Start server
	log.Fatal(app.Listen(":" + c.Config.Port))
}