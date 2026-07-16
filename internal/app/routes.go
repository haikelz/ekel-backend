package app

import (
	guestbookHTTP "ekel-backend/internal/domains/guestbook/delivery/http"
	guestbookRepository "ekel-backend/internal/domains/guestbook/repository"
	guestbookUsecase "ekel-backend/internal/domains/guestbook/usecase"
	ihsgHTTP "ekel-backend/internal/domains/ihsg/delivery/http"
	ihsgUsecase "ekel-backend/internal/domains/ihsg/usecase"
	reactionHTTP "ekel-backend/internal/domains/reaction/delivery/http"
	reactionRepository "ekel-backend/internal/domains/reaction/repository"
	reactionUsecase "ekel-backend/internal/domains/reaction/usecase"
	wakatimeHTTP "ekel-backend/internal/domains/wakatime/delivery/http"
	wakatimeUsecase "ekel-backend/internal/domains/wakatime/usecase"
	"ekel-backend/internal/env"

	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type RouteGroups struct {
	Admin    *echo.Group
	Customer *echo.Group
	Public   *echo.Group
}

func BuildRouteGroups(e *echo.Echo) RouteGroups {
	config := env.GetEnv()

	adminGroup := e.Group("/admin")
	adminGroup.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(config.SECRET_KEY_ADMIN)}))

	customerGroup := e.Group("/customer")
	customerGroup.Use(echojwt.WithConfig(echojwt.Config{SigningKey: []byte(config.SECRET_KEY_CUSTOMER)}))

	publicGroup := e.Group("/api/v1")

	return RouteGroups{
		Admin:    adminGroup,
		Customer: customerGroup,
		Public:   publicGroup,
	}
}

func RegisterModules(_ *echo.Echo, db *gorm.DB, routeGroups RouteGroups) {
	// Guestbook module
	guestbookRepo := guestbookRepository.NewGuestbookRepository(db)
	guestbookUc := guestbookUsecase.NewGuestbookUsecase(guestbookRepo)
	guestbookHTTP.NewGuestbookHandler(routeGroups.Public, guestbookUc)

	// IHSG module
	ihsgUc := ihsgUsecase.NewIhsgUsecase()
	ihsgHTTP.NewIhsgHandler(routeGroups.Public, ihsgUc)

	// Wakatime module
	wakatimeUc := wakatimeUsecase.NewWakatimeUsecase()
	wakatimeHTTP.NewWakatimeHandler(routeGroups.Public, wakatimeUc)

	// Reaction module
	reactionRepo := reactionRepository.NewReactionRepository(db)
	reactionUc := reactionUsecase.NewReactionUsecase(reactionRepo)
	reactionHTTP.NewReactionHandler(routeGroups.Public, reactionUc)
}
