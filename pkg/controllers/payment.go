package controllers

import (
	"net/http"
	"project/initializers"
	"project/pkg/models"

	"github.com/gin-gonic/gin"
)

func AddPaymentMethod(g *gin.Context) {

	var body struct {
		Paymentmethod string
	}
	if g.Bind(&body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}
	println("qwewqw", body.Paymentmethod)
	new := models.Payment{
		PaymentMethod: body.Paymentmethod,
	}
	result := initializers.DB.Create(&new)
	if result.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add payment method",
		})
		return
	}
	g.JSON(http.StatusOK, "payment added surcessfully")
}
