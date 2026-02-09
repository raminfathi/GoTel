package api

import (
	"fmt"

	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
)

func getAuthUser(c fiber.Ctx) (*types.User, error) {
	user, ok := c.Locals("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}
	return user, nil
}
