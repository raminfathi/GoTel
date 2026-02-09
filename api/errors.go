package api

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
	"github.com/raminfathi/GoTel/types" // Import types here as well
)

func ErrorHandler(c fiber.Ctx, err error) error {
	if apiError, ok := err.(types.Error); ok {
		return c.Status(apiError.Code).JSON(apiError)
	}
	apiError := types.NewError(http.StatusInternalServerError, err.Error())
	return c.Status(apiError.Code).JSON(apiError)
}
