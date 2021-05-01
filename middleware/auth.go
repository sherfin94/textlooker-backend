package middleware

import (
	"errors"
	"log"
	"os"
	"textlooker-backend/models"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type login struct {
	Email    string `json:"Email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

var identityKey string = "id"

func GenerateJWTAuthMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "\"auth zone\"",
		Key:         []byte(os.Getenv("JWT_SECRET_KEY")),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(models.User); ok {
				return jwt.MapClaims{
					identityKey: v.ID,
				}
			}
			return jwt.MapClaims{}
		},

		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			var user models.User
			models.Database.First(&user, claims[identityKey])
			return &user
		},

		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			var user models.User

			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			email := loginVals.Email
			password := loginVals.Password

			result := models.Database.Where("email = ?", email).First(&user)

			if errors.Is(result.Error, nil) {
				err := bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(password))
				if errors.Is(err, nil) {
					return user, nil
				}
			}

			return nil, jwt.ErrFailedAuthentication
		},

		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},

		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TimeFunc: time.Now,
	})

	if err != nil {
		log.Fatal("JWT Error:" + err.Error())
	}

	return authMiddleware
}
