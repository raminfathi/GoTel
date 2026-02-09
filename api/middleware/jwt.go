package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c fiber.Ctx) error {
		token := c.Get("X-Api-Token")

		if token == "" {
			fmt.Println("token not present in the header")
			return types.ErrUnAuthorized()
		}

		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		if time.Now().Unix() > expires {
			return types.NewError(fiber.StatusUnauthorized, "token expired")
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return types.ErrUnAuthorized()
		}

		// Set the current authenticated user to the context.
		c.Locals("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signing method", token.Header["alg"])
			return nil, types.ErrUnAuthorized()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, types.ErrUnAuthorized()
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, types.ErrUnAuthorized()
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, types.ErrUnAuthorized()
	}
	return claims, nil
}
