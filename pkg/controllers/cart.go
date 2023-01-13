package controllers

import (
	"net/http"
	"project/initializers"
	"project/pkg/function"
	"project/pkg/models"
	"project/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddToCart(g *gin.Context) {
	userID := function.GetUserId(g)
	var cartid int
	initializers.DB.Raw("SELECT ID FROM carts WHERE user_id=?", userID).Scan(&cartid)

	var body struct {
		ProductItemID int
		Quantity      int
	}

	if g.Bind(&body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	currentQuantity := function.ProductQuantity(body.ProductItemID)
	println("qq", currentQuantity)

	if body.Quantity > currentQuantity {
		g.JSON(400, gin.H{"error": "out of stock"})
		return
	}

	var product models.Product
	initializers.DB.Raw("SELECT total, product_name FROM products WHERE ID=?", body.ProductItemID).Scan(&product)
	total := product.Total * body.Quantity
	println("total", total)
	add_cart_items := models.ShoppingCartItem{
		CartId:        cartid,
		ProductItemID: body.ProductItemID,
		ProductName:   product.ProductName,
		Quantity:      body.Quantity,
		Total:         total,
	}
	add_to_cart_result := initializers.DB.Create(&add_cart_items)
	if add_to_cart_result.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add product",
		})
		return
	}

	// respond
	g.JSON(http.StatusOK, gin.H{"update": "product added to cart surcessfully"})
}

func Cart(g *gin.Context) {
	userID := function.GetUserId(g)
	println(userID, "user id")

	var carttotal int
	var cart models.Cart
	initializers.DB.Raw("select sum(total) from shopping_cart_items where cart_id=? ", userID).Scan(&carttotal)

	initializers.DB.Raw("update carts set total=? where id=?", carttotal, userID).Scan(&cart)

	var balance int
	initializers.DB.Raw("select balance from walets where id=?", userID).Scan(&balance)

	var total int
	initializers.DB.Raw("select total from carts where user_id=?", userID).Scan(&total)
	println("db total cart ", total)
	cartdisplay := response.Cart{
		Total: total,
		Walet: balance,
	}

	g.JSON(http.StatusOK, cartdisplay)

	var id []int
	initializers.DB.Raw("SELECT product_item_id from shopping_cart_items WHERE cart_id=?", userID).Scan(&id)

	var cartproduct []response.CartProduct
	for _, v := range id {
		println(v, "id......")
		initializers.DB.Raw("select *from products where id=?", v).Scan(&cartproduct)
		g.JSON(http.StatusOK, cartproduct)
	}

}

func UpdateCartItemQuantity(g *gin.Context) {
	userID := function.GetUserId(g)
	var body struct {
		ProductItemID int
		Quantity      int
	}
	if g.Bind(&body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	currentproductQuantity := function.ProductQuantity(body.ProductItemID)
	if body.Quantity > currentproductQuantity {
		g.JSON(400, gin.H{"error": "out of stock"})
		return
	}
	var price int
	initializers.DB.Raw("SELECT total FROM products WHERE ID=?", body.ProductItemID).Scan(&price)
	total := price * body.Quantity

	var UpdateCartItemQuantity models.ShoppingCartItem
	initializers.DB.Raw("update shopping_cart_items SET quantity=? WHERE cart_id=? AND product_item_id=?", body.Quantity, userID, body.ProductItemID).Scan(&UpdateCartItemQuantity)
	initializers.DB.Raw("update shopping_cart_items SET total=? WHERE cart_id=? AND product_item_id=?", total, userID, body.ProductItemID).Scan(&UpdateCartItemQuantity)
	g.JSON(200, gin.H{"update": "quantity updated"})

}
func CartItemDelete(g *gin.Context) {
	userID := function.GetUserId(g)
	params := g.Query("id")
	productId, _ := strconv.Atoi(params)
	var product models.Product
	initializers.DB.Raw(" DELETE FROM shopping_cart_items WHERE cart_id=? AND product_item_id=?", userID, productId).Scan(&product)

	g.JSON(http.StatusOK, gin.H{
		"message": "Product is deleted",
	})
}
