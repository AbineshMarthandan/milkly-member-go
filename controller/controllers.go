package controller

import (
	"milkly-member/config"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
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
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Health check endpoint
	e.GET("/health", func(ctx echo.Context) error {
		return ctx.JSON(200, map[string]interface{}{
			"status":  "healthy",
			"service": "milkly-member",
		})
	})

	// Swagger endpoint
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// API routes
	api := e.Group("/api")

	// Member routes
	member := api.Group("/member")
	member.POST("", c.MemberController.CreateMember)
	member.GET("", c.MemberController.ListMembers)
	member.GET("/:id", c.MemberController.GetMember)
	member.PUT("/:id", c.MemberController.UpdateMember)
	member.DELETE("/:id", c.MemberController.DeleteMember)

	// Start server
	e.Logger.Fatal(e.Start(":" + c.Config.Port))
}