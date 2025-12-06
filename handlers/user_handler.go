package handlers

import (
	"net/http"

	"github.com/barzaevhalid/cloud_storage_backend/services"
	"github.com/barzaevhalid/cloud_storage_backend/utils"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	UserService *services.UserService
}

type registerReq struct {
	Email    string `json:"email" validate:"required,email"`
	FullName string `json:"fullname" validate:"required,min=3"`
	Password string `json:"password" validate:"required,min=3"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=3"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func NewUserHandler(s *services.UserService) *UserHandler {
	return &UserHandler{
		UserService: s,
	}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req registerReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid body"})
	}
	if err := utils.Validate.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	user, err := h.UserService.Register(req.Email, req.FullName, req.Password)

	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(user)
}

// POST /login
func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}

	if err := utils.Validate.Struct(req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	accessToken, refreshToken, err := h.UserService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

//POST refresh

func (h *UserHandler) Refresh(c *fiber.Ctx) error {
	var req RefreshRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request"})
	}
	newToken, err := h.UserService.RefreshToken(req.RefreshToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"access_token": newToken})
}

func (h *UserHandler) GetMe(c *fiber.Ctx) error {
	userId := c.Locals("user_id").(int64)

	u, err := h.UserService.GetMe(userId)

	if err != nil {
		return err
	}

	return c.JSON(u)
}
