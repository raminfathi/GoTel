package types

import "go.mongodb.org/mongo-driver/v2/bson"

type Hotel struct {
	ID       bson.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string          `bson:"name" json:"name"`
	Location string          `bson:"location" json:"location"`
	Rooms    []bson.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int             `bson:"rating" json:"rating"`
}

// RoomType uses an Enum for better type safety
type RoomType int

const (
	_ RoomType = iota
	SingleRoomType
	DoubleRoomType
	SeaViewRoomType
	DeluxeRoomType
)

type Room struct {
	ID        bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Type      RoomType      `bson:"type" json:"type"`
	BasePrice float64       `bson:"basePrice" json:"basePrice"`
	Price     float64       `bson:"price" json:"price"`
	HotelID   bson.ObjectID `bson:"hotelID" json:"hotelId"`
}

type HotelQueryParams struct {
	Rating int `query:"rating" json:"rating"`
}
type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}
