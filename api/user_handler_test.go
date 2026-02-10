package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/raminfathi/GoTel/types"
)

func TestPostUser(t *testing.T) {
	// 1. Setup
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.store.User)
	app.Post("/", userHandler.HandlePostUser)

	uniqueEmail := fmt.Sprintf("test_%d@user.com", time.Now().UnixNano())

	params := types.CreateUserParams{
		Email:     uniqueEmail,
		FirstName: "Test",
		LastName:  "User",
		Password:  "supersecurepassword",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 200 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		t.Errorf("expected 200 but got %d. Response Body: %v", resp.StatusCode, errResp)
	}

	user, err := tdb.store.User.GetUserByEmail(context.TODO(), uniqueEmail)
	if err != nil {
		t.Error("User was not saved in DB:", err)
	}
	if user.Email != uniqueEmail {
		t.Errorf("expected email %s but got %s", uniqueEmail, user.Email)
	}
}
