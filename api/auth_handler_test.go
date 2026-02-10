package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v3"
	"github.com/raminfathi/GoTel/types"
)

func TestAuthenticate(t *testing.T) {
	// 1. Setup Test Database
	tdb := setup(t)
	defer tdb.teardown(t)

	// 2. Setup Fiber App & Handler
	app := fiber.New()
	// Note: AuthHandler needs UserStore, so we pass tdb.store.User
	authHandler := NewAuthHandler(tdb.store.User)
	app.Post("/auth", authHandler.HandleAuthenticate)

	// 3. Insert a test user into the database
	// We use NewUserFromParams to ensure password gets hashed correctly
	userParams := types.CreateUserParams{
		Email:     "cool@user.com",
		FirstName: "Cool",
		LastName:  "User",
		Password:  "supersecret123",
	}
	user, err := types.NewUserFromParams(userParams)
	if err != nil {
		t.Fatal(err)
	}
	_, err = tdb.store.User.InsertUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}

	// ---------------------------------------------------------
	// Test Case 1: Successful Login
	// ---------------------------------------------------------
	loginParams := types.AuthParams{
		Email:    "cool@user.com",
		Password: "supersecret123",
	}
	body, _ := json.Marshal(loginParams)
	req := httptest.NewRequest("POST", "/auth", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected http status 200 but got %d", resp.StatusCode)
	}

	// Decode response to check for token
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatal(err)
	}

	if _, ok := response["token"]; !ok {
		t.Error("expected response to contain 'token' but it did not")
	}

	// ---------------------------------------------------------
	// Test Case 2: Wrong Password
	// ---------------------------------------------------------
	wrongParams := types.AuthParams{
		Email:    "cool@user.com",
		Password: "wrongpassword",
	}
	body, _ = json.Marshal(wrongParams)
	req = httptest.NewRequest("POST", "/auth", bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")

	resp, err = app.Test(req)
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("expected http status 400 (or 401) but got %d", resp.StatusCode)
	}
}
