package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"go-edge-specailist/config/inits"
	"go-edge-specailist/controllers"
	"go-edge-specailist/services"
)

func main() {

	// Set Gin to release mode to disable debug output
	gin.SetMode(gin.ReleaseMode)
	// Create a new Gin router
	router := gin.Default()
	inits.InitHomePage(router)

	// Handle Log Request & Response
	//router.Use(middleware.LogHandler())

	// Register Specialist Controller
	specialistService := services.NewSpecialistService()
	specialistController := controllers.NewSpecialistController(specialistService)
	router.GET("/api/v1/specialist/read", specialistController.ReadSpecialist)
	router.GET("/api/v1/specialist/query", specialistController.FindByCustomerNo)
	router.POST("/api/v1/specialist/save", specialistController.SaveSpecialist)
	router.DELETE("/api/v1/specialist/delete", specialistController.DeleteByIndex)
	router.GET("/api/v1/specialist/export", specialistController.ExportSpecialist)

	// Open browser on start
	//inits.OpenBrowser("http://localhost:8083")

	// Run the server on port 8080
	err := router.Run(":8083")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot start on port 8083")
		return
	}
}
