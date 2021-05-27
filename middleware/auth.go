package middleware

import (
	"errors"
	"log"
	"net/http"
	"textlooker-backend/database"
	"textlooker-backend/deployment"
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

var identityKey string = "user"

func GenerateJWTAuthMiddleware() *jwt.GinJWTMiddleware {
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:          "\"auth zone\"",
		Key:            []byte(deployment.GetEnv("JWT_SECRET_KEY")),
		Timeout:        30 * time.Minute,
		MaxRefresh:     30 * time.Minute,
		IdentityKey:    identityKey,
		SendCookie:     true,
		SecureCookie:   true,    //non HTTPS dev environments
		CookieHTTPOnly: true,    // JS can't modify
		CookieName:     "token", // default jwt
		TokenLookup:    "cookie:token",
		CookieSameSite: http.SameSiteNoneMode, //SameSiteDefaultMode, SameSiteLaxMode, SameSiteStrictMode, SameSiteNoneMode

		LoginResponse: func(context *gin.Context, status int, s string, t time.Time) {
			context.AbortWithStatus(status)
		},

		RefreshResponse: func(context *gin.Context, status int, s string, t time.Time) {
			context.AbortWithStatus(status)
		},

		LogoutResponse: func(context *gin.Context, status int) {
			context.AbortWithStatus(status)
		},

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
			database.Database.First(&user, claims[identityKey])
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

			result := database.Database.Where("email = ?", email).First(&user)

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
		log.Println("JWT Error:" + err.Error())
	}

	return authMiddleware
}
