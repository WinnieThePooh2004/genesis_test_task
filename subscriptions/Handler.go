package subscriptions

import (
	"github.com/gofiber/fiber/v2"
	"regexp"
)

type Handler struct {
	service IService
}

var emailRegex *regexp.Regexp = regexp.MustCompile("^[a-z0-9!#$%&'*+/=?^_`{|}~-]+(?:\\.[a-z0-9!#$%&'*+/=?^_`{|}~-]+)*@(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?$")

func NewHandler(service IService) *Handler {
	return &Handler{service: service}
}

func (h *Handler) AddSubscription(c *fiber.Ctx) error {
	email := c.FormValue("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is required"})
	}

	if !emailRegex.MatchString(email) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email is invalid"})
	}

	exist, err := h.service.Add(email)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if exist {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "email already exists"})
	}

	return c.SendStatus(fiber.StatusCreated)
}
