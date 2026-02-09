package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/raminfathi/GoTel/db/fixtures"
)

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	// یوزر ساخته میشه: ایمیل James@Foo.com و پسورد James_Foo
	fixtures.AddUser(tdb.Store, "James", "Foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "James@Foo.com",
		Password: "password_kamelan_ghalat", // <--- این باید غلط باشه!
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	// انتظار داریم ۴۰۰ بگیریم
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400 but got %d", resp.StatusCode)
	}
}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	insertedUser := fixtures.AddUser(tdb.Store, "James", "Foo", false)

	app := fiber.New()
	authHandler := NewAuthHandler(tdb.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "James@Foo.com",
		Password: "James_Foo", // <--- این باید درست باشه
	}
	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d", resp.StatusCode)
	}

	var authResp AuthResponse
	if err := json.NewDecoder(resp.Body).Decode(&authResp); err != nil {
		t.Fatal(err)
	}

	// پاک کردن پسورد هش شده برای مقایسه دقیق
	insertedUser.EncryptedPassword = ""
	// هندل کردن اختلاف زمان جزئی (اختیاری)
	insertedUser.CreatedAt = authResp.User.CreatedAt

	if !reflect.DeepEqual(insertedUser, authResp.User) {
		fmt.Printf("Expected: %+v\n", insertedUser)
		fmt.Printf("Got:      %+v\n", authResp.User)
		t.Fatalf("user mismatch")
	}
}
