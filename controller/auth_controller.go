
package controller

import (
	"milkly-member/entity"
	"milkly-member/service"

	"github.com/gofiber/fiber/v2"
	"github.com/go-playground/validator/v10"
)

type AuthController struct {
	authService service.AuthService
	validator   *validator.Validate
}

func NewAuthController(authService service.AuthService) *AuthController {
	return &AuthController{
		authService: authService,
		validator:   validator.New(),
	}
}

// Login godoc
// @Summary Login member
// @Description Authenticate a member and return a JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body entity.LoginRequest true "Login credentials"
// @Success 200 {object} entity.AuthResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /auth/login [post]
func (c *AuthController) Login(ctx *fiber.Ctx) error {
	var req entity.LoginRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := c.validator.Struct(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	response, err := c.authService.Login(ctx.Context(), &req)
	if err != nil {
		if err.Error() == "invalid credentials" {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return ctx.JSON(response)
}
