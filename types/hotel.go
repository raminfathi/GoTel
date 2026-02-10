package types

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Hotel struct {
	ID       bson.ObjectID   `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string          `bson:"name" json:"name"`
	Location string          `bson:"location" json:"location"`
	Rooms    []bson.ObjectID `bson:"rooms" json:"rooms"`
	Rating   int             `bson:"rating" json:"rating"`
}

type HotelQueryParams struct {
	Rating int `query:"rating" json:"rating"`
}
type ResourceResp struct {
	Results int `json:"results"`
	Data    any `json:"data"`
	Page    int `json:"page"`
}

type CreateHotelParams struct {
	Name     string `json:"name" validate:"required,min=3,max=100"`
	Location string `json:"location" validate:"required,min=3,max=100"`
}

type UpdateHotelParams struct {
	Name     string `json:"name" validate:"omitempty,min=3,max=100"` // omitempty یعنی اگر نفرستاد گیر نده
	Location string `json:"location" validate:"omitempty,min=3,max=100"`
}

func (p CreateHotelParams) Validate() map[string]string {
	errors := map[string]string{}
	if len(p.Name) < 3 {
		errors["name"] = "name must be at least 3 characters"
	}
	if len(p.Location) < 3 {
		errors["location"] = "location must be at least 3 characters"
	}
	return errors
}
