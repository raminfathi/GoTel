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

// HandleGetRooms returns rooms of a specific hotel
// @Summary      Get hotel rooms
// @Description  Get all rooms belonging to a specific hotel
// @Tags         hotel
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Hotel ID"
// @Param        X-Api-Token header string true "Token"
// @Success      200  {array}   types.Room
// @Router       /hotel/{id}/rooms [get]
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

// HandleGetHotel returns a single hotel
// @Summary      Get hotel by ID
// @Description  Get detailed information of a specific hotel
// @Tags         hotel
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Hotel ID"
// @Param        X-Api-Token header string true "Token"
// @Success      200  {object}  types.Hotel
// @Failure      404  {object}  map[string]string
// @Router       /hotel/{id} [get]
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

type HotelQueryParams struct {
	db.Pagination
	Rating int
}

// HandleGetHotels returns all hotels
// @Summary      Get all hotels
// @Description  Get a list of all hotels with filtering options
// @Tags         hotel
// @Accept       json
// @Produce      json
// @Param        X-Api-Token header string true "Token"
// @Success      200  {array}   types.Hotel
// @Router       /hotel [get]
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

// HandlePostHotel adds a new hotel (Admin only)
// @Summary      Add a hotel
// @Description  Add a new hotel
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request body types.CreateHotelParams true "Hotel Data"
// @Param        X-Api-Token header string true "Token"
// @Success      200  {object}  types.Hotel
// @Router       /admin/hotel [post]
func (h *HotelHandler) HandlePostHotel(c fiber.Ctx) error {
	var params types.CreateHotelParams
	if err := c.Bind().Body(&params); err != nil {
		return types.ErrBadRequest()
	}

	if errors := params.Validate(); len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	hotel := types.Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   0,
		Rooms:    []bson.ObjectID{},
	}

	insertedHotel, err := h.store.Hotel.InsertHotel(c.Context(), &hotel)
	if err != nil {
		return err
	}

	return c.JSON(insertedHotel)
}

// HandlePutHotel (Admin Only)
// HandlePutHotel updates a hotel (Admin only)
// @Summary      Update a hotel
// @Description  Update hotel details
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id      path    string                true  "Hotel ID"
// @Param        request body    types.UpdateHotelParams true  "Update Data"
// @Param        X-Api-Token header string true "Token"
// @Success      200     {object}  map[string]string
// @Router       /admin/hotel/{id} [put]
func (h *HotelHandler) HandlePutHotel(c fiber.Ctx) error {
	id := c.Params("id")

	var params types.UpdateHotelParams
	if err := c.Bind().Body(&params); err != nil {
		return types.ErrBadRequest()
	}

	updateData := db.Map{}
	if params.Name != "" {
		updateData["name"] = params.Name
	}
	if params.Location != "" {
		updateData["location"] = params.Location
	}

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return types.ErrInvalidID()
	}

	filter := db.Map{"_id": oid}

	if err := h.store.Hotel.UpdateHotel(c.Context(), filter, updateData); err != nil {
		return err
	}

	// h.store.Cache.Delete(c.Context(), "hotel-"+id)

	return c.JSON(db.Map{"msg": "updated successfully"})
}
