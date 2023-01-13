package controllers

import (
	"net/http"
	"project/initializers"
	"project/pkg/function"
	"project/pkg/models"
	"time"

	"github.com/gin-gonic/gin"
)

func AddCoupen(g *gin.Context) {
	var body struct {
		Code   string
		Value  int
		Limit  int64
		IsUsed bool
	}

	// binding json to go response writter
	if g.Bind(&body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}
	coupen := models.Coupen{
		Code:   body.Code,
		Value:  body.Value,
		Limit:  body.Limit,
		IsUsed: false,
		Expire: time.Now().Add(time.Hour * 24 * 30).Unix(),
	}
	add_coupen := initializers.DB.Create(&coupen)

	if add_coupen.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add product",
		})
		return
	}
}

func ValidateCoupen(g *gin.Context) {

	userid := function.GetUserId(g)
	// coupen check
	coupen_code := g.Query("code")
	println("cccc", coupen_code)
	var coupen models.Coupen
	initializers.DB.Raw("SELECT *FROM coupens where code=?", coupen_code).Scan(&coupen)
	if len(coupen_code) < 5 || coupen.ID == 0 {
		println("1")
		g.JSON(400, gin.H{"error": "enter valid coupen code"})
		return
	}
	if coupen.IsUsed {
		println("2")
		g.JSON(400, gin.H{"error": "coupen already used"})
		return
	}
	var total int
	println("3")
	initializers.DB.Raw("select total from carts where id=? ", userid).Scan(&total)
	if !coupen.IsUsed && coupen.Limit < int64(total) {
		var walet models.Walet
		initializers.DB.Raw("select balance from walets where id=? ", userid).Scan(&walet)
		walet.Balance += coupen.Value
		initializers.DB.Raw("update walets set balance=? where id=?", walet.Balance, userid).Scan(&walet)

		initializers.DB.Raw("UPDATE coupens SET is_used=? WHERE id=?", true, coupen.ID).Scan(&coupen)
		g.JSON(http.StatusOK, "coupen added surcessfully")
	}

}
