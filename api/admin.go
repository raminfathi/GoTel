package api

import (
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
)

func AdminAuth(c fiber.Ctx) error {
	user, ok := c.Locals("user").(*types.User)
	if !ok {
		return ErrUnAuthorized()
	}
	if !user.IsAdmin {
		return ErrUnAuthorized()

	}
	return c.Next()
}
