package types

import "go.mongodb.org/mongo-driver/v2/bson"

type RoomType int

const (
	Single RoomType = iota + 1
	Double
	SeaView
	KingSuite
)

type Room struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType      `bson:"type" json:"type"`
	BasePrice float64       `bson:"basePrice" json:"basePrice"`
	Price     float64       `bson:"price" json:"price"`
	HotelID   bson.ObjectID `bson:"hotelID" json:"hotelId"`
}

// ---------------------------------------------
// ورودی کاربر (DTO)
// ---------------------------------------------
type CreateRoomParams struct {
	HotelID   string  `json:"hotelId"`
	Type      string  `json:"type"`
	BasePrice float64 `json:"basePrice"`
}

func (p CreateRoomParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(p.HotelID) == 0 {
		errors["hotelId"] = "hotelId is required"
	}
	if len(p.Type) < 2 {
		errors["type"] = "type must be at least 2 characters"
	}
	if p.BasePrice <= 0 {
		errors["basePrice"] = "price must be greater than 0"
	}
	return errors
}
