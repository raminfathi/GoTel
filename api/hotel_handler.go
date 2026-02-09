package api

import (
	"fmt"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"

	"github.com/gofiber/fiber/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type HotelHandler struct {
	store *db.Store
}

func NewHotelHandler(store *db.Store) *HotelHandler {
	return &HotelHandler{
		store: store,
	}
}
func (h *HotelHandler) HandleGetRooms(c fiber.Ctx) error {
	id := c.Params("id")
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID()
	}
	filter := bson.M{"HotelID": oid}

	rooms, err := h.store.Room.GetRooms(c.Context(), filter)
	if err != nil {
		return types.ErrResourceNotFound("hotel")
	}
	return c.JSON(rooms)
}
func (h *HotelHandler) HandleGetHotel(c fiber.Ctx) error {
	fmt.Println(c.Params)

	id := c.Params("id")
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	fmt.Println(hotel)
	if err != nil {
		return types.ErrResourceNotFound("hotel")
	}
	return c.JSON(hotel)
}

type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.Bind().Query(&params); err != nil {
		return types.ErrBadRequest()
	}
	filter := db.Map{}
	if params.Rating > 0 {
		filter["rating"] = params.Rating
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return types.ErrResourceNotFound("hotels")
	}
	resp := ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(params.Page),
	}
	return c.JSON(resp)
}
