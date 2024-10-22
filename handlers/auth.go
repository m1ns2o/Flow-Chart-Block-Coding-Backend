// handlers/auth.go
package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT secret key for signing tokens
var jwtSecret []byte

// Claims structure for JWT
type Claims struct {
	ClassID  uint   `json:"class_id"`
	Classnum string `json:"classnum"`
	jwt.RegisteredClaims
}

// SetJWTSecret sets the JWT secret key
func SetJWTSecret(secret []byte) {
	jwtSecret = secret
}

// generateToken creates a new JWT token
func generateToken(classID uint, classnum string) (string, error) {
	claims := Claims{
		ClassID:  classID,
		Classnum: classnum,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AuthMiddleware validates JWT tokens
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
			c.Abort()
			return
		}

		c.Set("class_id", claims.ClassID)
		c.Set("classnum", claims.Classnum)
		c.Next()
	}
}