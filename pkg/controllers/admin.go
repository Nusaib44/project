package controllers

import (
	"fmt"
	"net/http"
	"project/initializers"
	"project/pkg/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!user manangement!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!
func ListUser(g *gin.Context) {
	var user []models.Userdata
	initializers.DB.Find(&user)
	g.JSON(http.StatusOK, user)
}

func BlockUser(g *gin.Context) {

	params := g.Param("id")
	fmt.Println(params)
	page, _ := strconv.Atoi(params)
	var users models.Userdata
	initializers.DB.Raw("update Userdata SET Status=false WHERE id=?", page).Scan(&users)
	g.JSON(http.StatusOK, gin.H{"": "user boceked surcessfully", "id": page})
}
func Unblock(g *gin.Context) {
	params := g.Param("id")
	fmt.Printf("%T", params)
	page, _ := strconv.Atoi(params)
	var users models.Userdata
	initializers.DB.Raw("update Userdata SET Status=true WHERE id=?", page).Scan(&users)
	g.JSON(http.StatusOK, gin.H{"": "user unblocked surcessfully", "id": page})
}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!Product manangement!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

func AddProduct(g *gin.Context) {

	// data storage of user from body
	var PP struct {
		ProductName string
		Category    int
		Brand       string
		Price       int
		Quantity    int
		Image       string
		Description string
	}

	// binding json to go response writter
	if g.Bind(&PP) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	// createing product
	add_product := models.Product{
		ProductName: PP.ProductName,
		Category:    PP.Category,
		Brand:       PP.Brand,
		Quantity:    PP.Quantity,
		Price:       PP.Price,
		Total:       PP.Price,
		Image:       PP.Image,
		Description: PP.Description,
	}

	product_result := initializers.DB.Create(&add_product)

	if product_result.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add product",
		})
		return
	}
	// respond
	g.JSON(http.StatusOK, gin.H{"update": "product added surcessfully"})
}
func EditProduct(g *gin.Context) {
	var PP struct {
		ProductName string
		Category    int
		Brand       string
		Price       int
		Quantity    int
		Image       string
		Description string
	}

	params := g.Query("id")
	page, _ := strconv.Atoi(params)
	var product models.Product

	if err := g.Bind(&PP); err != nil {
		println("errrrr......not binded", err.Error())

		g.JSON(400, err.Error())
		return
	}
	initializers.DB.First(&product, page)
	if product.ID < 1 {
		g.JSON(400, gin.H{"message": "product not found"})
		return
	}
	initializers.DB.Model(&product).Updates(models.Product{
		ProductName: PP.ProductName,
		Category:    PP.Category,
		Brand:       PP.Brand,
		Price:       PP.Price,
		Quantity:    PP.Quantity,
		Image:       PP.Image,
		Description: PP.Description,
	})

	g.JSON(http.StatusOK, gin.H{"message": "product updated surcessfully"})
}
func DeleteProduct(g *gin.Context) {
	params := g.Query("id")
	page, _ := strconv.Atoi(params)
	var product models.Product

	initializers.DB.Raw(" DELETE FROM products WHERE id=?", page).Scan(&product)
	g.JSON(http.StatusOK, gin.H{"message": "Product is deleted"})

}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!category manangement!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

func AddCategory(g *gin.Context) {

	var CC struct{ Categoryname string }

	if g.Bind(&CC) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	AddCategory := models.Category{Categoryname: CC.Categoryname}
	category_result := initializers.DB.Create(&AddCategory)
	if category_result.Error != nil {
		g.JSON(400, gin.H{
			"error": "failed to add product",
		})
		return
	}
	// respond
	g.JSON(http.StatusOK, gin.H{"update": "category added surcessfully"})
}
func EditCategory(g *gin.Context) {
	params := g.Query("id")
	fmt.Printf("%T", params)
	page, _ := strconv.Atoi(params)

	var CC struct {
		Categoryname string
	}
	if g.Bind(&CC) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}
	// checking category
	var check models.Category
	initializers.DB.First(&check, "id= ?", page)
	if check.ID < 1 {
		g.JSON(400, gin.H{"error": "category not found"})
		return
	}

	var EditCategory models.Category
	if CC.Categoryname != "" {
		initializers.DB.Raw("update categories SET Categoryname=? WHERE id=?", CC.Categoryname, page).Scan(&EditCategory)
		g.JSON(http.StatusOK, gin.H{"update": "category name edited surcessfully"})
	}

}
func DeleteCategory(g *gin.Context) {
	params := g.Query("id")
	fmt.Printf("%T", params)
	page, _ := strconv.Atoi(params)
	var product models.Product
	initializers.DB.Raw(" DELETE FROM categories WHERE id=?", page).Scan(&product)

	g.JSON(http.StatusOK, gin.H{
		"message": "Category is deleted",
	})
}

// !!!!!!!!!!!!!!!!!!!!!!!!!!!!!!coupon manangement!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!!

func EditCoupon(g *gin.Context) {

}
