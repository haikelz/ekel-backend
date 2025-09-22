package server

import (
	"guestbook-backend/pkg/configs"
	"guestbook-backend/pkg/controllers"
	"guestbook-backend/pkg/server/routes"
	"guestbook-backend/pkg/services"
)

func (s *FiberApp) RegisterFiberRoutes() {
	// routes.HomeRoute(s)
	// routes.SwaggerRoute(s)
	// routes.PrometheusRoute(s)

	db := configs.NewDB()

	guestbookService := services.NewGuestbookService(db)
	guestbookController := controllers.NewGuestbookController(guestbookService)

	// routes.CreateGuestbookRoute(s, guestbookController)
	routes.GetGuestbookRoute(s, guestbookController)
	// routes.DeleteGuestbookRoute(s, guestbookController)
	// routes.UpdateGuestbookRoute(s, guestbookController)
	// routes.LoginAdminRoute(s, guestbookController)
	// routes.ValidateCookieRoute(s, guestbookController)
	// routes.GetAllUserByEmailRoute(s, guestbookController)
	// routes.DeleteUserByEmailRoute(s, guestbookController)
}
