package middleware

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"github.com/wafi04/chatting-app/services/shared/types"
)

var jwtSecretKey = []byte("jsjxakabxjaigisyqyg189")

type JWTClaims struct {
	UserID          string `json:"user_id"`
	Email           string `json:"email"`
	Name            string `json:"name"`
	Role            string `json:"role"`
	IsActive        bool   `json:"is_active"`
	IsEmailVerified bool   `json:"is_email_verified"`
	jwt.StandardClaims
}

func ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
func GenerateToken(user *types.UserInfo, day int64) (string, error) {
	claims := JWTClaims{
		UserID:          user.UserId,
		Email:           user.Email,
		Name:            user.Name,
		IsEmailVerified: user.IsEmailVerified,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(day) * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "wafiuddin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, nil
}

type contextKey string

const UserContextKey contextKey = "user"

func GetUserFromGinContext(c *gin.Context) (*types.UserInfo, error) {
	user, exists := c.Get(string(UserContextKey))
	if !exists {
		return nil, errors.New("user not found in context")
	}
	userInfo, ok := user.(*types.UserInfo)
	if !ok {
		return nil, errors.New("invalid user type in context")
	}
	return userInfo, nil
}
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Step 1: Cek cookie untuk access_token
		accessToken, err := c.Cookie("access_token")
		if err == nil && accessToken != "" {
			// Validasi token dari cookie
			if claims, err := ValidateToken(accessToken); err == nil {
				user := &types.UserInfo{
					UserId:          claims.UserID,
					Email:           claims.Email,
					Name:            claims.Name,
					IsEmailVerified: claims.IsEmailVerified,
				}
				fmt.Printf("token : %s", err)
				// Set user di context
				c.Set(string(UserContextKey), user)
				c.Next()
				return
			}
		}

		// Step 2: Jika tidak ada access_token valid, cek refresh_token dari cookie
		refreshToken, err := c.Cookie("refresh_token")
		if err == nil && refreshToken != "" {
			// Validasi refresh token
			if claims, err := ValidateToken(refreshToken); err == nil {
				user := &types.UserInfo{
					UserId:          claims.UserID,
					Email:           claims.Email,
					Name:            claims.Name,
					IsEmailVerified: claims.IsEmailVerified,
				}
				// Generate new access token
				newAccessToken, err := GenerateToken(user, 24) // Token baru berlaku 24 jam
				if err != nil {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"error": "Failed to generate new access token",
					})
					return
				}
				// Set new access token sebagai cookie
				SetAccressTokenCookie(c, newAccessToken)
				// Set user di context
				c.Set(string(UserContextKey), user)
				c.Next()
				return
			}
		}

		// Step 3: Jika tidak ada token valid, return unauthorized
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"error": "No valid tokens found",
		})
	}
}
func SetRefreshTokenCookie(c *gin.Context, token string) {
	log.Printf("addedd called")

	c.SetCookie(
		"refresh_token",
		token,
		int(168*3600),
		"/",
		"192.168.100.9",
		false,
		true,
	)
}

func SetAccressTokenCookie(c *gin.Context, token string) {
	log.Printf("addedd called")

	c.SetCookie(
		"access_token",
		token,
		int(24*3600),
		"/",
		"192.168.100.9",
		false,
		true,
	)
}

func SetSessionCookie(c *gin.Context, sessionID string) {
	log.Printf("addedd called")
	c.SetCookie(
		"auth_session",
		sessionID,
		int(168*3600),
		"/",
		"192.168.100.9",
		false,
		true,
	)
}

func ClearTokens(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true,
		true,
	)

	c.Header("Authorization", "")
}

func ResponseTime(r *gin.Engine) {
	log.Printf("calledd")
	r.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		c.Header("X-Response-Time", duration.String())
	})
}
