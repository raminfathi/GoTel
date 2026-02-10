package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/raminfathi/GoTel/db/fixtures"
	"github.com/raminfathi/GoTel/types"
)

func TestBookRoom(t *testing.T) {
	tdb := setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	roomHandler := NewRoomHandler(tdb.store)

	user := fixtures.AddUser(tdb.store, "james", "bond", false)
	hotel := fixtures.AddHotel(tdb.store, "Grand Hotel", "London", 5, nil)
	room := fixtures.AddRoom(tdb.store, types.Double, 100.0, hotel.ID)

	app.Post("/room/:id/book", func(c fiber.Ctx) error {
		c.Locals("user", user)
		return c.Next()
	}, roomHandler.HandleBookRoom)

	bookingParams := types.BookRoomParams{
		FromDate:   time.Now(),
		TillDate:   time.Now().AddDate(0, 0, 5),
		NumPersons: 2,
	}

	body, _ := json.Marshal(bookingParams)
	route := fmt.Sprintf("/room/%s/book", room.ID.Hex())

	req, _ := http.NewRequest("POST", route, bytes.NewReader(body))
	req.Header.Add("Content-Type", "application/json")

	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200 but got %d", resp.StatusCode)
	}

	var booking types.Booking
	if err := json.NewDecoder(resp.Body).Decode(&booking); err != nil {
		t.Fatal(err)
	}

	if booking.UserID != user.ID {
		t.Errorf("expected booking user ID %s but got %s", user.ID, booking.UserID)
	}
	if booking.RoomID != room.ID {
		t.Errorf("expected booking room ID %s but got %s", room.ID, booking.RoomID)
	}

	if booking.NumPersons != 2 {
		t.Errorf("expected 2 persons but got %d", booking.NumPersons)
	}
}
