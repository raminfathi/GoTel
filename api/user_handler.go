package api

import (
	"errors"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserHandler struct {
	userStore db.UserStore
}

func NewUserHandler(userStore db.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

func (h *UserHandler) HandlePutUser(c fiber.Ctx) error {
	var params types.UpdateUserParams
	userId := c.Params("id")

	if err := c.Bind().Body(&params); err != nil {
		return types.ErrBadRequest()
	}
	filter := db.Map{"_id": userId}
	if err := h.userStore.UpdateUser(c.Context(), filter, params); err != nil {
		return c.Status(400).JSON(map[string]string{"error": err.Error()})
	}

	return c.JSON(map[string]string{"message": "user updated successfully", "id": userId})
}

func (h *UserHandler) HandleDeleteUser(c fiber.Ctx) error {
	userId := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
		return types.ErrBadRequest()
	}

	return c.JSON(map[string]string{"message": "user deleted successfully", "id": userId})
}
func (h *UserHandler) HandlePostUser(c fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.Bind().Body(&params); err != nil {
		return err
	}
	if errors := ValidateRequest(params); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	_, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err == nil {
		// ارور nil یعنی یوزر پیدا شد -> پس تکراریه
		return c.Status(fiber.StatusBadRequest).JSON(map[string]string{
			"error": "email already exists",
		})
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *UserHandler) HandleGetUser(c fiber.Ctx) error {

	var (
		id = c.Params("id")
	)
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
		return err
	}

	return c.JSON(user)
}

func (h *UserHandler) HandleGetUsers(c fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return types.ErrResourceNotFound("user")
	}
	return c.JSON(users)
}
