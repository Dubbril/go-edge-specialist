package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go-edge-specailist/models"
	"go-edge-specailist/services"
	"net/http"
)

type SpecialistController struct {
	SpecialistService *services.SpecialistService
}

func NewSpecialistController(specialistService *services.SpecialistService) *SpecialistController {
	return &SpecialistController{SpecialistService: specialistService}
}

func (ctrl *SpecialistController) ReadSpecialist(c *gin.Context) {
	selectEnv := c.Query("selectEnv")
	fmt.Println(selectEnv)

	//file, err := c.FormFile("file")
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "Missing file parameter", "status": "Fail"})
	//	return
	//}

	err := ctrl.SpecialistService.ReadSpecialist(selectEnv)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "Fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Success"})
}

func (ctrl *SpecialistController) ExportSpecialist(c *gin.Context) {
	selectEnv := c.Query("selectEnv")
	fmt.Println(selectEnv)

	err := ctrl.SpecialistService.ExportSpecialist(selectEnv)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Error on export specialist with error " + err.Error(), "status": "Fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Success"})
}

func (ctrl *SpecialistController) SaveSpecialist(c *gin.Context) {
	// Bine @RequestBody @Valid AesRequest aesRequest in java
	var specialistReq models.SpecialistRequest
	if err := c.ShouldBindJSON(&specialistReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error(), "status": "Fail"})
		return
	}

	err := ctrl.SpecialistService.SaveSpecialist(specialistReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "status": "Fail"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "Success"})
}

func (ctrl *SpecialistController) FindByCustomerNo(c *gin.Context) {
	customerNoFilter := c.Query("customerNoFilter")
	if customerNoFilter == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "CustomerNo is empty", "RowNo": ""})
		return
	}

	specialistData, err := ctrl.SpecialistService.FilterByCustomerNo(customerNoFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "RowNo": ""})
		return
	}
	c.JSON(http.StatusOK, specialistData)
}

func (ctrl *SpecialistController) DeleteByIndex(c *gin.Context) {
	value := c.Query("rowNo")

	var err = ctrl.SpecialistService.DeleteByIndex(value)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err, "status": "Fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "Success"})
}
