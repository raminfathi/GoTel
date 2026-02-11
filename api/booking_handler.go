package api

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"

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

// HandleGetMyBookings returns user bookings
// @Summary      Get my bookings
// @Description  Get all bookings for the logged-in user
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        X-Api-Token header string true "Token"
// @Success      200  {array}   types.Booking
// @Router       /booking [get]
func (h *BookingHandler) HandleGetMyBookings(c fiber.Ctx) error {
	user, err := getAuthUser(c)
	if err != nil {
		return types.ErrUnAuthorized()
	}

	filter := bson.M{"userID": user.ID}

	bookings, err := h.store.Booking.GetBookings(c.Context(), filter)
	if err != nil {
		return types.ErrResourceNotFound("bookings")
	}

	if bookings == nil {
		return c.JSON([]*types.Booking{})
	}

	return c.JSON(bookings)
}

// HandleCancelBooking cancels a booking
// @Summary      Cancel booking
// @Description  Cancel a booking (Method is POST/PUT for safety)
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Booking ID"
// @Param        X-Api-Token header string true "Token"
// @Success      200  {object}  map[string]string
// @Router       /booking/{id}/cancel [post]
func (h *BookingHandler) HandleCancelBooking(c fiber.Ctx) error {
	id := c.Params("id")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return types.ErrResourceNotFound("booking")
	}
	user, err := getAuthUser(c)
	if err != nil {
		return types.ErrUnAuthorized()
	}
	if booking.UserID != user.ID {
		return types.ErrUnAuthorized()
	}
	if booking.Canceled {
		return c.Status(fiber.StatusBadRequest).JSON(genericResp{
			Type: "error",
			Msg:  "booking already canceled",
		})
	}
	if err := h.store.Booking.UpdateBooking(c.Context(), id, bson.M{"canceled": true}); err != nil {
		return err
	}
	return c.JSON(genericResp{
		Type: "msg",
		Msg:  "updated",
	})
}

// HandleGetBookings returns all bookings (Admin only)
// @Summary      Get all bookings
// @Description  Get a list of all bookings in the system
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        X-Api-Token header string true "Token"
// @Success      200  {array}   types.Booking
// @Router       /admin/booking [get]
func (h *BookingHandler) HandleGetBookings(c fiber.Ctx) error {
	booking, err := h.store.Booking.GetBookings(c.Context(), bson.M{})
	if err != nil {
		return types.ErrResourceNotFound("bookings")
	}
	fmt.Println(booking)
	return c.JSON(booking)

}

// HandleGetBooking returns a specific booking
// @Summary      Get booking details
// @Description  Get a single booking by ID
// @Tags         booking
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Booking ID"
// @Param        X-Api-Token header string true "Token"
// @Success      200  {object}  types.Booking
// @Router       /booking/{id} [get]
func (h *BookingHandler) HandleGetBooking(c fiber.Ctx) error {
	id := c.Params("id")
	cacheKey := fmt.Sprintf("booking-%s", id)
	val, err := h.store.Cache.Get(c.Context(), cacheKey)
	if err == nil && val != "" {
		fmt.Println("--->> Serving from CACHE")

		var booking types.Booking
		if err := json.Unmarshal([]byte(val), &booking); err == nil {
			if err := h.checkBookingOwner(c, &booking); err != nil {
				return err
			}
			return c.JSON(booking)
		}
	}
	fmt.Println("--->> Serving from MongoDB")
	booking, err := h.store.Booking.GetBookingByID(c.Context(), id)
	if err != nil {
		return types.ErrResourceNotFound("booking")
	}
	if err := h.checkBookingOwner(c, booking); err != nil {
		return err
	}

	serialized, err := json.Marshal(booking)
	if err == nil {
		h.store.Cache.Set(c.Context(), cacheKey, serialized, time.Minute*1)
	}
	return c.JSON(booking)
}

func (h *BookingHandler) checkBookingOwner(c fiber.Ctx, booking *types.Booking) error {
	user, err := getAuthUser(c)
	if err != nil {
		return types.ErrUnAuthorized()
	}
	if booking.UserID != user.ID && !user.IsAdmin {
		return types.ErrUnAuthorized()
	}
	return nil

}
