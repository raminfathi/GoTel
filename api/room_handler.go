package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type RoomHandler struct {
	store *db.Store
}

func NewRoomHandler(store *db.Store) *RoomHandler {
	return &RoomHandler{
		store: store,
	}
}
func (h *RoomHandler) HandleGetRooms(c fiber.Ctx) error {

	cacheKey := "rooms-" + c.OriginalURL()
	val, err := h.store.Cache.Get(c.Context(), cacheKey)
	if err == nil && val != "" {
		var rooms []*types.Room
		if err := json.Unmarshal([]byte(val), &rooms); err == nil {
			return c.JSON(rooms)
		}
	}
	filter := bson.M{}
	if hotelID := c.Query("hotelId"); hotelID != "" {
		oid, err := bson.ObjectIDFromHex(hotelID)
		if err != nil {
			return types.ErrInvalidID()
		}
		filter["hotelID"] = oid
	}
	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return types.ErrResourceNotFound("rooms")
	}
	serialized, err := json.Marshal(rooms)
	if err == nil {
		h.store.Cache.Set(c.Context(), cacheKey, serialized, time.Minute*1)
	}

	return c.JSON(rooms)
}
func (h *RoomHandler) HandleBookRoom(c fiber.Ctx) error {
	var params types.BookRoomParams
	if err := c.Bind().Body(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}
	roomID, err := bson.ObjectIDFromHex(c.Params("id"))
	if err != nil {
		return types.ErrInvalidID()
	}
	userCtx := c.Locals("user")

	user, ok := userCtx.(*types.User)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error (user not found in context)",
		})
	}
	ok, err = h.store.Booking.IsRoomAvailable(c.Context(), roomID, params.FromDate, params.TillDate)
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
