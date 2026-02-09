package api

import (
	"fmt"
	"net/http"

	"github.com/raminfathi/GoTel/db"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type BookingHandler struct {
	store *db.Store
}

func NewBookingHandler(store *db.Store) *BookingHandler {
	return &BookingHandler{
		store: store,
	}
}
func (h *BookingHandler) HandleCancelBooking(c fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrNotResourceNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return ErrUnAuthorized()
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})
}

// TODO: this needs to be admin authorized
func (h *BookingHandler) HandleGetBookings(c fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return ErrNotResourceNotFound("bookings")
	}
	fmt.Println(booking)
	return c.JSON(booking)

}

// TODO: this needs to be user authorized
func (h *BookingHandler) HandleGetBooking(c fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return ErrUnAuthorized()
	}
	user, err := getAuthUser(c)
	if err != nil {
		return ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return c.Status(http.StatusUnauthorized).JSON(genericResp{
			Type: "error",
			Msg:  "not authorized",
		})
	}
	return c.JSON(booking)

}
