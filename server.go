package main

import (
	"github.com/labstack/echo"
	"github.com/samkrew/media_stats/routes"
	"github.com/samkrew/media_stats/logger"
	"github.com/samkrew/media_stats/grabber"
)

const PORT = "3000"

func main() {
	go grabber.StartGrabber()

	e := echo.New()
	routes.Set(e)
	logger.L.Infof("Start video stats server on port: %s", PORT)
	logger.L.Fatal(e.Start(":" + PORT))
}