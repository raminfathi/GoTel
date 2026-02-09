package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p *BookRoomParams) validate() error {
	if p.FromDate.IsZero() || p.TillDate.IsZero() {
		return fmt.Errorf("dates cannot be empty")
	}
	now := time.Now()
	if now.After(p.FromDate) || now.After(p.TillDate) {
		return fmt.Errorf("cannot book a room in the past")
	}
	return nil
}

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}
func (h *RoomHandler) HandleGetRooms(c fiber.Ctx) error {
	rooms, err := h.store.Room.GetRooms(c.Context(), bson.M{})
	if err != nil {
		return err
	}
	return c.JSON(rooms)
}
func (h *RoomHandler) HandleBookRoom(c fiber.Ctx) error {
	var params BookRoomParams
	if err := c.Bind().Body(&params); err != nil {
		return err
	}
	if err := params.validate(); err != nil {
		return err
	}
	roomID, err := bson.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return err
	}
	userCtx := c.Locals("user")
	// if !ok {
	// 	return c.Status(http.StatusInternalServerError).JSON(genericResp{
	// 		Type: "error",
	// 		Msg:  "internal server error",
	// 	})
	// }
	user, ok := userCtx.(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error (user not found in context)",
		})
	}
	ok, err = h.isRoomAvailableForBooking(c.Context(), roomID, params)
	if err != nil {
		return err
	}
	if !ok {
		return c.Status(http.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  fmt.Sprintf("room %s already booked", c.Params("id")),
		})
	}
	booking := types.Booking{
		UserID:     user.ID,
		RoomID:     roomID,
		FromDate:   params.FromDate,
		TillDate:   params.TillDate,
		NumPersons: params.NumPersons,
	}

	inserted, err := h.store.Booking.InsertBooking(c.Context(), &booking)
	if err != nil {
		return err
	}
	return c.JSON(inserted)
}

func (h *RoomHandler) isRoomAvailableForBooking(ctx context.Context, roomID bson.ObjectID, params BookRoomParams) (bool, error) {

	where := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$lt": params.TillDate,
		},
		"tillDate": bson.M{
			"$gt": params.FromDate,
		},
	}

	fmt.Printf("Searching for RoomID: %s\n", roomID.Hex())
	fmt.Printf("Filter: %+v\n", where)
	bookings, err := h.store.Booking.GetBookings(ctx, where)
	if err != nil {
		return false, err
	}
	ok := len(bookings) == 0
	return ok, nil

}
