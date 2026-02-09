package types

import "go.mongodb.org/mongo-driver/v2/bson"

type Hotel struct {
	ID       bson.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string          `bson:"name" json:"name"`
	Location string          `bson:"location" json:"location"`
	Rooms    []bson.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int             `bson:"rating" json:"rating"`
}

type Room struct {
	ID      bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Size    string        `bson:"size" json:"size"` //small, normal, kingsize
	SeaSide bool          `bson:"seaside" json:"seaside"`
	Price   float64       `bson:"Price" json:"Price"`
	HotelID bson.ObjectID `bson:"HotelID" json:"HotelID"`
}
