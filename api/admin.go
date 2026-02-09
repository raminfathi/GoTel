package api

import (
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
)

func AdminAuth(c fiber.Ctx) error {
	user, ok := c.Locals("user").(*types.User)
	if !ok {
		return types.ErrUnAuthorized()
	}
	if !user.IsAdmin {
		return types.ErrUnAuthorized()

	}
	return c.Next()
}
