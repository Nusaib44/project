package controllers

import (
	"net/http"
	"project/initializers"
	"project/pkg/function"
	"project/pkg/models"
	"project/response"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserInfo(g *gin.Context) {

	var user models.Userdata
	var add []models.Address
	var order []models.Order
	userID := function.GetUserId(g)

	initializers.DB.Raw("SELECT *FROM orders WHERE user_id=?", userID).Scan(&order)
	initializers.DB.Raw("SELECT *FROM addresses WHERE user_id=?", userID).Scan(&add)
	initializers.DB.Raw("SELECT *FROM userdata WHERE user_id=?", userID).Scan(&user)

	userinfo := response.UserInfo{
		UserDetails: user,
		Address:     add,
		Order:       order,
	}
	g.JSON(http.StatusOK, userinfo)
}

func EditUserInfo(g *gin.Context) {

	var user models.Userdata
	id := function.GetUserId(g)
	var Body struct {
		Username    string
		Email       string
		PhoneNumber string
		Password    string
	}

	if g.Bind(&Body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	if Body.Username != "" {
		initializers.DB.Model(&models.Userdata{}).Where("id=?", id).Update("username", Body.Username)
		initializers.DB.Raw("update userdata SET username=? WHERE id=?", Body.Username, id).Scan(&user)
		g.JSON(http.StatusOK, gin.H{"update": "username updated surcessfully"})
	}
	if len(Body.PhoneNumber) == 10 {
		initializers.DB.Raw("update userdata SET phone_number=? WHERE id=?", Body.PhoneNumber, id).Scan(&user)
		g.JSON(http.StatusOK, gin.H{"update": "phone_number updated surcessfully"})
	}
	if Body.Email != "" {
		initializers.DB.Raw("update userdata SET email=? WHERE id=?", Body.Email, id).Scan(&user)
		g.JSON(http.StatusOK, gin.H{"update": "Email  updated surcessfully"})
	}

}

func ChangePassword(g *gin.Context) {
	var user models.Userdata
	var Body struct{ pass string }
	println("vkyvcjgfhvc", Body.pass)
	id := function.GetUserId(g)
	pass := g.Query("password")

	if g.Bind(&Body) != nil {
		println("errrrr......not binded")
		g.JSON(400, gin.H{"error": "failed to load"})
		return
	}

	if Body.pass != "" {
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(Body.pass))
		if err != nil {
			g.JSON(http.StatusBadRequest, gin.H{"error": "invalid password"})
		} else {
			hash, hash_err := bcrypt.GenerateFromPassword([]byte(pass), 10)
			if hash_err != nil {
				g.JSON(400, gin.H{"error": "failed to hash passsword"})
				return
			}
			initializers.DB.Raw("update userdata SET password=? WHERE id=?", hash, id).Scan(&user)
			g.JSON(http.StatusOK, gin.H{"update": "password  updated surcessfully"})

		}
	}

}
