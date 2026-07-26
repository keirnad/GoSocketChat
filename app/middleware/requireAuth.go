package middleware

import (
	"gochatroom/initializers"
	"gochatroom/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func RequireAuth(c *gin.Context) {
	tokenString, err := c.Cookie("Authorization")

	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
	} else {
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(os.Getenv("SECRET")), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			log.Fatal(err)
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if float64(time.Now().Unix()) > claims["exp"].(float64) {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			var user models.User
			initializers.DB.First(&user, claims["sub"])

			if user.ID == 0 {
				c.AbortWithStatus(http.StatusUnauthorized)
			}

			c.Set("user", user)
			c.Next()

		} else {
			c.AbortWithStatus(http.StatusUnauthorized)
		}
	}
}
