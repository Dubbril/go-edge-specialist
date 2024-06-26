package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type HomeController struct{}

func NewHomeController() *HomeController {
	return &HomeController{}
}

func (h *HomeController) Index(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func (h *HomeController) FaviconHandler(c *gin.Context) {
	c.File("./favicon.ico")
}
