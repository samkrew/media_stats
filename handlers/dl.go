package handlers

import (
	"github.com/labstack/echo"
	"net/http"
	"github.com/asaskevich/govalidator"
	"github.com/samkrew/media_stats/api"
	"github.com/samkrew/media_stats/grabber"
	"github.com/samkrew/media_stats/logger"
	"github.com/samkrew/media_stats/db"
)

func GrabFileStats(c echo.Context) error {
	var validationErrors []api.Error

	url := c.QueryParam("url")
	hash := c.QueryParam("md5")

	if !govalidator.IsURL(url) {
		validationErrors = append(validationErrors, api.ErrorUrlInvalid)
	}
	if !govalidator.IsHexadecimal(hash) || len(hash) != 32 {
		validationErrors = append(validationErrors, api.ErrorHashInvalid)
	}

	if len(validationErrors) != 0 {
		return c.JSON(http.StatusBadRequest, api.ErrorResponse(validationErrors))
	}

	id, err := grabber.AddTask(url, hash)

	if err != nil {
		logger.L.Errorf("New grabber task add error: %v", err)
		return c.JSON(http.StatusInternalServerError, api.ErrorResponse(nil))
	}

	err = db.NewTaskStats(id, url, hash)
	if err != nil {
		logger.L.Errorf("[Task #%v] Stats create: %v", err)
	}

	logger.L.Debugf("[Task #%v] %v", id, url)


	return c.JSON(http.StatusOK, api.SuccessResponse(api.EmptyPayload))
}
