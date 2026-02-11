package api

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

// HandleAuthenticate authenticates a user
// @Summary      User Login
// @Description  Login with email and password to get a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body types.AuthParams true "Login Credentials"
// @Success      200  {object}  map[string]string
// @Failure      400  {object}  map[string]string
// @Failure      401  {object}  map[string]string
// @Router       /auth [post]
func (h *AuthHandler) HandleAuthenticate(c fiber.Ctx) error {
	fmt.Println("\n🔥🔥🔥 HANDLER REACHED! 🔥🔥🔥")
	rawBody := c.Body()
	fmt.Printf("📦 Raw Body from Swagger: %s\n", string(rawBody))
	var params types.AuthParams
	if err := c.Bind().Body(&params); err != nil {
		return types.ErrBadRequest()
	}
	fmt.Println("--- DEBUG LOGIN ---")
	fmt.Printf("Input Params: %+v\n", params)
	fmt.Printf("🔍 RECEIVED PARAMS: Email='%s' | Password='%s'\n", params.Email, params.Password)

	fmt.Println("-------------------")

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		// احتمالا اینجا داره ارور میده
		fmt.Println("DB Error:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	if errors := ValidateRequest(params); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	user, err = h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
	if err != nil {
		fmt.Println("Password mismatch error:", err)
		return invalidCredentials(c)
	}

	token := CreateTokenFromUser(user)
	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

func invalidCredentials(c fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func CreateTokenFromUser(user *types.User) string {
	now := time.Now()
	validUntil := now.Add(time.Hour * 4).Unix()

	claims := jwt.MapClaims{
		"id":      user.ID,
		"email":   user.Email,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
