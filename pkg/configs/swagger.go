package configs

import (
	"github.com/gofiber/contrib/swagger"
)

var SwgCfg = swagger.Config{
	BasePath:    "/",
	Title:       "Guestbook Backend API Docs",
	Path:        "docs",
	FileContent: []byte(SwaggerJSON),
	CacheAge:    60,
}
