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

func PlaceOrder(g *gin.Context) {

	userId := function.GetUserId(g)

	pay := g.Query("payment")
	payment, _ := strconv.Atoi(pay)
	add := g.Query("address")
	address, _ := strconv.Atoi(add)

	var itemId []int
	initializers.DB.Raw("SELECT ID FROM shopping_cart_items WHERE cart_id=?", userId).Scan(&itemId)

	for _, r := range itemId {
		var cartItems models.ShoppingCartItem
		var Payment string
		initializers.DB.Raw("SELECT product_item_id, quantity, total FROM shopping_cart_items WHERE ID=?", r).Scan(&cartItems)

		var balance int
		initializers.DB.Raw("select balance from walets where id=?", userId).Scan(&balance)

		var cart models.Cart
		var walet models.Walet

		newtotal := cartItems.Total - balance

		if newtotal < 0 {
			println("11111111s")
			bal := 0 - newtotal
			initializers.DB.Raw("update walets set balance=? where id=?", bal, userId).Scan(&walet)
			initializers.DB.Raw("update carts set total=? where id=?", 0, userId).Scan(&walet)
			cartItems.Total = 0
		} else {
			println("22222222")
			initializers.DB.Raw("update carts set total=? where id=?", newtotal, userId).Scan(&cart)
			balance = 0
			initializers.DB.Raw("update walets set balance=? where id=?", balance, userId).Scan(&walet)
			cartItems.Total = newtotal
		}
		println(payment, ".......payment")
		if payment == 2 {
			StripePayment(g)
		}

		initializers.DB.Raw("SELECT payment_method FROM payments WHERE ID=?", payment).Scan(&Payment)
		order := models.Order{
			UserID:    userId,
			ProductId: cartItems.ProductItemID,
			Quantity:  cartItems.Quantity,
			Price:     cartItems.Total,
			Address:   address,
			Payment:   Payment,
			Status:    "shiped",
		}
		result := initializers.DB.Create(&order)
		if result.Error != nil {
			g.JSON(400, gin.H{"error": "failed to place order"})
			return
		}
		var currentQuantity int
		var productTable models.Product
		initializers.DB.Raw("SELECT quantity FROM products WHERE id=?", cartItems.ProductItemID).Scan(&currentQuantity)
		newQuantity := currentQuantity - cartItems.Quantity
		initializers.DB.Raw("update products SET quantity=? WHERE id=?", newQuantity, cartItems.ProductItemID).Scan(&productTable)
		initializers.DB.Raw(" DELETE FROM shopping_cart_items WHERE ID=?", r).Scan(&Payment)
	}
	// respond
	g.JSON(http.StatusOK, gin.H{"update": "order placed surcessfully"})
}
func ListOrder(g *gin.Context) {
	userID := function.GetUserId(g)
	var orders []response.Update
	initializers.DB.Raw("SELECT *FROM orders WHERE user_id=?", userID).Scan(&orders)
	g.JSON(http.StatusOK, orders)
}

func OrderCancelation(g *gin.Context) {

	params := g.Query("id")
	orderId, _ := strconv.Atoi(params)
	var order models.Order
	var product response.Update
	var currentQuantity int
	var productTable models.Product
	// stock update
	initializers.DB.Raw("SELECT product_id, quantity FROM orders WHERE id=?", orderId).Scan(&product)
	initializers.DB.Raw("SELECT quantity FROM products WHERE id=?", product.ProductId).Scan(&currentQuantity)
	newQuantity := currentQuantity + product.Quantity
	initializers.DB.Raw("update products SET quantity=? WHERE id=?", newQuantity, product.ProductId).Scan(&productTable)
	// delete
	initializers.DB.Raw(" DELETE FROM orders WHERE ID=?", orderId).Scan(&order)
	g.JSON(http.StatusOK, gin.H{
		"message": "order cancelled",
	})
}

func RetutnOrder(g *gin.Context) {
	params := g.Query("id")
	userId := function.GetUserId(g)
	orderId, _ := strconv.Atoi(params)

	println(orderId, "orderid")
	var product models.Order
	initializers.DB.Raw("SELECT product_id, quantity,price FROM orders WHERE id=?", orderId).Scan(&product)
	initializers.DB.Raw("update orders set status=? where product=? and user")
	order := models.Return{
		UserID:    userId,
		ProductID: product.ProductId,
		Quantity:  product.Quantity,
		Price:     product.Price,
		Status:    "Processing",
	}
	g.JSON(http.StatusOK, order)

	result := initializers.DB.Create(&order)
	if result.Error != nil {
		g.JSON(400, gin.H{"error": "failed to return order"})
		return
	}
	g.JSON(http.StatusOK, gin.H{
		"message": "order returned",
	})
	var walet models.Walet
	var balance int
	initializers.DB.Raw("select balance from walets where id=?", userId).Scan(&balance)
	newbalance := balance + product.Price
	initializers.DB.Raw("UPDATE walets SET balance=? where id=?", newbalance, userId).Scan(&walet)
}
