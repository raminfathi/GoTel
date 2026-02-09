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
	"go.mongodb.org/mongo-driver/v2/mongo" // Driver v2
	"golang.org/x/crypto/bcrypt"
)

// ---------------------------------------------------------
// تعریف استراکت‌ها (که باعث خطای undefined شده بودند)
// ---------------------------------------------------------

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User `json:"user"`
	Token string      `json:"token"`
}

// این همون استراکتی هست که توی فایل‌های دیگه (booking, room) استفاده شده
type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

// ---------------------------------------------------------
// سازنده (Constructor)
// ---------------------------------------------------------

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{
		userStore: userStore,
	}
}

// ---------------------------------------------------------
// هندلر اصلی لاگین (سازگار با Fiber v3)
// ---------------------------------------------------------

func (h *AuthHandler) HandleAuthenticate(c fiber.Ctx) error {
	var params AuthParams
	// در Fiber v3 برای خواندن بادی
	if err := c.Bind().Body(&params); err != nil {
		return err
	}

	// 1. پیدا کردن یوزر از دیتابیس
	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
		return err
	}

	// 2. مقایسه پسورد هش شده
	// user.EncryptedPassword: چیزی که در دیتابیس هست (Hash)
	// params.Password: چیزی که کاربر وارد کرده (Plaintext)
	err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password))
	if err != nil {
		fmt.Println("Password mismatch error:", err) // جهت دیباگ
		return invalidCredentials(c)
	}

	// 3. ساخت توکن و بازگشت پاسخ
	token := CreateTokenFromUser(user)
	return c.JSON(AuthResponse{
		User:  user,
		Token: token,
	})
}

// ---------------------------------------------------------
// توابع کمکی
// ---------------------------------------------------------

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
		"id":      user.ID, // در v2 دیگه hex لازم نیست، خود ID سریالایز میشه
		"email":   user.Email,
		"expires": validUntil,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// اینجا باید سکرت واقعی رو از env بخونی، فعلا برای بیلد شدن اینو میذاریم
	// os.Getenv("JWT_SECRET")
	secret := os.Getenv("JWT_SECRET")

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token with secret", err)
	}
	return tokenStr
}
