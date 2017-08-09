package handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/samkrew/media_stats/api"
	"github.com/samkrew/media_stats/db"
)

func ViewStats(c echo.Context) error {
	stats := db.GetAllStats()
	return c.JSONPretty(http.StatusOK, api.SuccessResponse(stats), " ")
}
