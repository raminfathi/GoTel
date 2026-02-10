package api

import (
	"encoding/json"
	"time"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/gofiber/fiber/v3"
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
	cacheKey := "hotel-rooms-" + c.OriginalURL()

	val, err := h.store.Cache.Get(c.Context(), cacheKey)
	if err == nil && val != "" {
		var rooms []*types.Room
		if err := json.Unmarshal([]byte(val), &rooms); err == nil {
			return c.JSON(rooms)
		}
	}

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID()
	}
	filter := bson.M{"hotelID": oid}
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
func (h *HotelHandler) HandleGetHotel(c fiber.Ctx) error {

	id := c.Params("id")

	cacheKey := "hotel-" + id

	val, err := h.store.Cache.Get(c.Context(), cacheKey)
	if err == nil && val != "" {
		var hotel types.Hotel
		if err := json.Unmarshal([]byte(val), &hotel); err == nil {
			return c.JSON(hotel)
		}
	}
	// oid, err := bson.ObjectIDFromHex(id)
	// if err != nil {
	// 	return types.ErrInvalidID()
	// }
	hotel, err := h.store.Hotel.GetHotelByID(c.Context(), id)
	if err != nil {
		return types.ErrResourceNotFound("hotel")
	}
	serialized, err := json.Marshal(hotel)
	if err == nil {
		h.store.Cache.Set(c.Context(), cacheKey, serialized, time.Minute*5)
	}

	return c.JSON(hotel)

}

// type ResourceResp struct {
// 	Results int `json:"results"`
// 	Data    any `json:"data"`
// 	Page    int `json:"page"`
// }

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

func (h *HotelHandler) HandleGetHotels(c fiber.Ctx) error {
	var params HotelQueryParams
	if err := c.Bind().Query(&params); err != nil {
		return types.ErrBadRequest()
	}
	cacheKey := "hotels-" + c.OriginalURL()
	val, err := h.store.Cache.Get(c.Context(), cacheKey)
	if err == nil && val != "" {
		var cachedResp types.ResourceResp
		if err := json.Unmarshal([]byte(val), &cachedResp); err == nil {
			return c.JSON(cachedResp)
		}

	}

	filter := db.Map{}
	if params.Rating > 0 {
		filter["rating"] = params.Rating
	}

	hotels, err := h.store.Hotel.GetHotels(c.Context(), filter, &params.Pagination)
	if err != nil {
		return types.ErrResourceNotFound("hotels")
	}

	resp := types.ResourceResp{
		Data:    hotels,
		Results: len(hotels),
		Page:    int(params.Page),
	}

	return c.JSON(resp)
	serialized, err := json.Marshal(resp)
	if err == nil {
		h.store.Cache.Set(c.Context(), cacheKey, serialized, time.Second*30)
	}
	return c.JSON(resp)
}
