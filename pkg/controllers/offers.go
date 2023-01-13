package controllers

import (
	"net/http"
	"project/initializers"
	"project/pkg/models"

	"github.com/gin-gonic/gin"
)

func AddProductOffer(g *gin.Context) {
	var body struct {
		Product int
		Offer   int
	}
	if g.Bind(&body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	// add offer
	productOffer := models.ProductOffers{
		Product:    body.Product,
		OfferPrice: body.Offer,
	}
	OfferResult := initializers.DB.Create(&productOffer)
	if OfferResult.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add Offer",
		})
		return
	}
	var price int
	var products models.Product
	initializers.DB.Raw("SELECT total from products where id=?", body.Product).Scan(&price)
	newprice := price - body.Offer
	println(price, "price")
	println("new ", newprice)
	initializers.DB.Raw("update products set total=?,product_offer=? where id=?", newprice, body.Offer, body.Product).Scan(&products)
	g.JSON(http.StatusOK, gin.H{
		"message": "product order added",
	})
}

func AddCategoryOffer(g *gin.Context) {

	var body struct {
		Category int
		Offer    int
	}
	if g.Bind(&body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	Offer := models.CategoryOffers{
		Category:   body.Category,
		OfferPrice: body.Offer,
	}
	OfferResult := initializers.DB.Create(&Offer)
	if OfferResult.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add Offer",
		})
		return
	}
	var id []int
	initializers.DB.Raw("SELECT id from products where category=?", body.Category).Scan(&id)
	for _, v := range id {
		var price int
		var products models.Product
		initializers.DB.Raw("SELECT total from products where id=?", v).Scan(&price)
		newprice := price - ((price * body.Offer) / 100)
		initializers.DB.Raw("update products set total=?,category_offer=? where id=?", newprice, body.Offer, v).Scan(&products)
	}
	g.JSON(http.StatusOK, gin.H{
		"message": "category offer added",
	})

}
