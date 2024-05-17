package rates

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type Handler struct {
	nbuRatesService *NbuRatesService
}

func NewHandler(nbuRatesService *NbuRatesService) *Handler {
	return &Handler{nbuRatesService: nbuRatesService}
}

func (h *Handler) GetNbuRate(c *fiber.Ctx) error {
	rate, err := h.nbuRatesService.GetRate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	_, err = c.Status(200).WriteString(strconv.FormatFloat(rate, 'f', 10, 64))
	return err
}
