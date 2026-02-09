package types

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

type Booking struct {
	ID         bson.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	UserID     bson.ObjectID `bson:"userID,omitempty" json:"userID,omitempty"`
	RoomID     bson.ObjectID `bson:"roomID,omitempty" json:"roomID,omitempty"`
	NumPersons int           `bson:"numPersons,omitempty" json:"numPersons,omitempty"`
	FromDate   time.Time     `bson:"fromDate,omitempty" json:"fromDate,omitempty"`
	TillDate   time.Time     `bson:"tillDate,omitempty" json:"tillDate,omitempty"`
	Canceled   bool          `bson:"canceled" json:"canceled"`
}
type BookRoomParams struct {
	FromDate   time.Time `json:"fromDate"`
	TillDate   time.Time `json:"tillDate"`
	NumPersons int       `json:"numPersons"`
}

func (p BookRoomParams) Validate() map[string]string {
	errors := map[string]string{}
	now := time.Now()

	if p.FromDate.Before(now) {
		errors["fromDate"] = "date must be in the future"
	}
	if p.TillDate.Before(p.FromDate) {
		errors["tillDate"] = "till date must be after from date"
	}
	if p.NumPersons <= 0 {
		errors["numPersons"] = "must be at least 1 person"
	}
	return errors
}
