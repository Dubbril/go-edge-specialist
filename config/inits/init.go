package inits

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-edge-specailist/controllers"
	"os/exec"
	"runtime"
)

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot open browser")
	}

}

func InitHomePage(router *gin.Engine) {
	// Serve static files from the "static" directory
	router.LoadHTMLGlob("views/*")
	router.Static("/static", "./static")
	router.Static("/css", "static/css")
	router.Static("/js", "static/js")

	// Register Home Controller
	homeController := controllers.NewHomeController()
	router.GET("/", homeController.Index)
	router.GET("/favicon.ico", homeController.FaviconHandler)

	// Use the logger middleware
	//router.Use(middleware.LogHandler())

}
