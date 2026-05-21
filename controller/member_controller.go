package controller

import (
	"net/http"
	"strconv"

	"milkly-member/entity"
	"milkly-member/service"

	"github.com/labstack/echo/v4"
	"github.com/go-playground/validator/v10"
)

type MemberController struct {
	memberService service.MemberService
	validator     *validator.Validate
}

func NewMemberController(memberService service.MemberService) *MemberController {
	return &MemberController{
		memberService: memberService,
		validator:     validator.New(),
	}
}

// CreateMember godoc
// @Summary Create a new member
// @Description Create a new member with the provided information
// @Tags members
// @Accept json
// @Produce json
// @Param member body entity.CreateMemberRequest true "Member data"
// @Success 201 {object} entity.MemberResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /members [post]
func (c *MemberController) CreateMember(ctx echo.Context) error {
	var req entity.CreateMemberRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	if err := c.validator.Struct(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	member, err := c.memberService.CreateMember(ctx.Request().Context(), &req)
	if err != nil {
		if err.Error() == "email already exists" {
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	return ctx.JSON(http.StatusCreated, member)
}

// GetMember godoc
// @Summary Get a member by ID
// @Description Get a member by their unique ID
// @Tags members
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Success 200 {object} entity.MemberResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /members/{id} [get]
func (c *MemberController) GetMember(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Member ID is required",
		})
	}

	member, err := c.memberService.GetMember(ctx.Request().Context(), id)
	if err != nil {
		if err.Error() == "member not found" {
			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	return ctx.JSON(http.StatusOK, member)
}

// UpdateMember godoc
// @Summary Update a member
// @Description Update a member's information
// @Tags members
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Param member body entity.UpdateMemberRequest true "Updated member data"
// @Success 200 {object} entity.MemberResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /members/{id} [put]
func (c *MemberController) UpdateMember(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Member ID is required",
		})
	}

	var req entity.UpdateMemberRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
	}

	if err := c.validator.Struct(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Validation failed",
			"details": err.Error(),
		})
	}

	member, err := c.memberService.UpdateMember(ctx.Request().Context(), id, &req)
	if err != nil {
		if err.Error() == "member not found" {
			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		if err.Error() == "email already exists" {
			return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	return ctx.JSON(http.StatusOK, member)
}

// DeleteMember godoc
// @Summary Delete a member
// @Description Delete a member by their ID
// @Tags members
// @Accept json
// @Produce json
// @Param id path string true "Member ID"
// @Success 204
// @Failure 400 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /members/{id} [delete]
func (c *MemberController) DeleteMember(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return ctx.JSON(http.StatusBadRequest, map[string]interface{}{
			"error": "Member ID is required",
		})
	}

	err := c.memberService.DeleteMember(ctx.Request().Context(), id)
	if err != nil {
		if err.Error() == "member not found" {
			return ctx.JSON(http.StatusNotFound, map[string]interface{}{
				"error": err.Error(),
			})
		}
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	return ctx.NoContent(http.StatusNoContent)
}

// ListMembers godoc
// @Summary List all members
// @Description Get a paginated list of all members
// @Tags members
// @Accept json
// @Produce json
// @Param limit query int false "Limit number of results" default(10)
// @Param offset query int false "Offset for pagination" default(0)
// @Success 200 {array} entity.MemberResponse
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /members [get]
func (c *MemberController) ListMembers(ctx echo.Context) error {
	limitStr := ctx.QueryParam("limit")
	offsetStr := ctx.QueryParam("offset")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	members, err := c.memberService.ListMembers(ctx.Request().Context(), limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
	}

	return ctx.JSON(http.StatusOK, map[string]interface{}{
		"data":   members,
		"limit":  limit,
		"offset": offset,
		"total":  len(members),
	})
}