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

// HandlePutUser updates a user
// @Summary      Update a user
// @Description  Update user details by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id      path    string                true  "User ID"
// @Param        request body    types.UpdateUserParams true  "Update Data"
// @Param        X-Api-Token header string true "Token"
// @Success      200     {object}  map[string]string
// @Router       /user/{id} [put]
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

// HandleDeleteUser deletes a user
// @Summary      Delete a user
// @Description  Delete a user by ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        X-Api-Token header string true "Token"
// @Success      200  {object}  map[string]string
// @Router       /user/{id} [delete]
func (h *UserHandler) HandleDeleteUser(c fiber.Ctx) error {
	userId := c.Params("id")

	if err := h.userStore.DeleteUser(c.Context(), userId); err != nil {
		return types.ErrBadRequest()
	}

	return c.JSON(map[string]string{"message": "user deleted successfully", "id": userId})
}

// HandlePostUser creates a new user (Registration)
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        request body types.CreateUserParams true "User Data"
// @Success      200  {object}  types.User
// @Failure      400  {object}  map[string]string
// @Router       /user [post]
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

// HandleGetUser returns a user by ID
// @Summary      Get a user
// @Description  Get a user by their ID
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "User ID"
// @Param        X-Api-Token header string true "Token"
// @Success      200  {object}  types.User
// @Failure      404  {object}  map[string]string
// @Router       /user/{id} [get]
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

// HandleGetUsers returns all users (Admin only)
// @Summary      Get all users
// @Description  Get a list of all registered users (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        X-Api-Token header string true "Token"
// @Success      200  {array}  types.User
// @Failure      403  {object}  map[string]string
// @Router       /admin/user [get]
func (h *UserHandler) HandleGetUsers(c fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return types.ErrResourceNotFound("user")
	}
	return c.JSON(users)
}
