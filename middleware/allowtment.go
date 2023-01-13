package middleware

import (
	"fmt"
	"net/http"
	"os"
	"project/initializers"
	"project/pkg/models"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

// var Jwtkey = []byte("secret_key")

func RequireAuth(g *gin.Context) {

	// controllers.Refreshtoken(g)
	println("middlewere running......")
	// geting cookie
	tokenString, err := g.Cookie("coookie")
	if err != nil {
		g.JSON(400, gin.H{
			"error": "login or signup to acess",
		})

		g.AbortWithStatus(http.StatusUnauthorized)
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Check the expm
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			g.JSON(401, gin.H{
				"error": "login to acess",
			})
			// g.AbortWithStatus(http.StatusUnauthorized)
		}
		// Find the user with token sub
		var user models.Userdata
		// initializers.DB.First(&user, claims["sub"])
		initializers.DB.First(&user, "email = ?", claims["sub"])

		if user.ID == 0 {
			g.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}

func AdminAuth(g *gin.Context) {
	println("adminmiddlewere running......")
	// geting cookie
	tokenString, err := g.Cookie("admincoookie")
	if err != nil {
		println("🥺🥺")
		g.AbortWithStatus(http.StatusUnauthorized)
	}

	// Parse takes the token string and a function for looking up the key. The latter is especially

	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(os.Getenv("SECRET")), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {

		// Check the expm
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			g.JSON(401, gin.H{
				"error": "you are not an admin",
			})
			g.AbortWithStatus(http.StatusUnauthorized)
		}
		// Find the user with token sub
		var user models.Userdata
		// initializers.DB.First(&user, claims["sub"])
		initializers.DB.First(&user, "email = ?", claims["sub"])

		if user.ID == 0 {

			g.AbortWithStatus(http.StatusUnauthorized)
		}

	}
}
