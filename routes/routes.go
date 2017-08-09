package routes

import (
	"github.com/labstack/echo"
	"github.com/samkrew/media_stats/handlers"
	"github.com/labstack/echo/middleware"
)

func Set(e *echo.Echo) {
	e.HideBanner = true
	e.Pre(middleware.RemoveTrailingSlash())

	e.GET("/dl", handlers.GrabFileStats)
	e.GET("/st", handlers.ViewStats)
}
