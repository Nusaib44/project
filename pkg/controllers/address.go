package controllers

import (
	"project/initializers"
	"project/pkg/function"
	"project/pkg/models"

	"github.com/gin-gonic/gin"
)

func AddNewAddress(g *gin.Context) {
	userid := function.GetUserId(g)
	var body struct {
		HouseName   string
		Street      string
		AddressLine string
		City        string
		State       string
		Pincode     int
		Country     string
		IsDefault   bool
	}
	if g.Bind(&body) != nil {
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	new := models.Address{
		UserId:      userid,
		HouseName:   body.HouseName,
		Street:      body.Street,
		AddressLine: body.AddressLine,
		City:        body.City,
		State:       body.State,
		Pincode:     body.Pincode,
		Country:     body.State,
		IsDefault:   body.IsDefault,
	}
	result := initializers.DB.Create(&new)
	if result.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add address... try again",
		})
		return
	}
}
