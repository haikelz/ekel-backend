package routes

import (
	"ekel-backend/pkg/controllers"

	"github.com/gofiber/fiber/v2"
)

func GetIHSGMarketsRoute(api fiber.Router, ihsgController *controllers.IHSGController) {
	api.Get("/api/v1/ihsg/markets", ihsgController.GetMarketsHandler)
}
