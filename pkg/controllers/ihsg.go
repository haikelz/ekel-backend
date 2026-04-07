package controllers

import (
	"net/http"

	"ekel-backend/pkg/models"
	"ekel-backend/pkg/services"

	"github.com/gofiber/fiber/v2"
)

type IHSGController struct {
	ihsgService *services.IHSGService
}

func NewIHSGController(ihsgService *services.IHSGService) *IHSGController {
	return &IHSGController{ihsgService: ihsgService}
}

func (ic *IHSGController) GetMarketsHandler(c *fiber.Ctx) error {
	markets, err := ic.ihsgService.GetMarkets()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(models.APIErrorResponse{
			StatusCode: http.StatusInternalServerError,
			Message:    "Failed to get market index data",
			Stack:      err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status_code": http.StatusOK,
		"message":     "Success get market index data",
		"data":        markets,
	})
}
