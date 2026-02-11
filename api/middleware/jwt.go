package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c fiber.Ctx) error {
		// ---------------------------------------------------------
		// 1. SKIP LOGIC (Allow Public Routes)
		// ---------------------------------------------------------
		// We explicitly check if the request is for Registration.
		// If it is POST /api/v1/user, we skip authentication.
		path := c.Path()
		method := c.Method()
		fmt.Printf("🛡️ Middleware Check -> Method: %s | Path: %s\n", method, path)
		if method == "OPTIONS" {
			return c.Next()
		}
		if strings.Contains(path, "/auth") {
			return c.Next()
		}
		if method == "POST" && strings.Contains(path, "/user") {
			return c.Next()
		}

		if strings.Contains(path, "/swagger") {
			return c.Next()
		}

		if len(path) > 8 && path[:8] == "/swagger" {
			return c.Next()
		}

		// ---------------------------------------------------------
		// 2. AUTHENTICATION LOGIC
		// ---------------------------------------------------------
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
