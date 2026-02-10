package db

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/raminfathi/GoTel/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BookingStore interface {
	InsertBooking(context.Context, *types.Booking) (*types.Booking, error)
	GetBookings(context.Context, bson.M) ([]*types.Booking, error)
	GetBookingByID(context.Context, string) (*types.Booking, error)
	UpdateBooking(context.Context, string, bson.M) error
	IsRoomAvailable(context.Context, bson.ObjectID, time.Time, time.Time) (bool, error)
}
type MongoBookingStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	if dbname == "" {
		dbname = "hotel_db"
	}

	return &MongoBookingStore{
		client: client,
		coll:   client.Database(dbname).Collection("bookings"),
	}

}
func (s *MongoBookingStore) IsRoomAvailable(ctx context.Context, roomID bson.ObjectID, from, till time.Time) (bool, error) {
	filter := bson.M{
		"roomID": roomID,
		"fromDate": bson.M{
			"$lt": till,
		},
		"tillDate": bson.M{
			"$gt": from,
		},
		"canceled": false,
	}

	count, err := s.coll.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	// If count is 0, it means no overlapping bookings found -> Room is available
	return count == 0, nil
}
func (s *MongoBookingStore) UpdateBooking(ctx context.Context, id string, update bson.M) error {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	m := bson.M{
		"$set": update,
	}
	_, err = s.coll.UpdateByID(ctx, oid, m)
	return err
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, id string) (*types.Booking, error) {
	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var booking types.Booking
	if err := s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&booking); err != nil {
		return nil, err

	}
	return &booking, nil
}
func (s *MongoBookingStore) GetBookings(ctx context.Context, filter bson.M) ([]*types.Booking, error) {
	curr, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var booking []*types.Booking
	if err := curr.All(ctx, &booking); err != nil {
		return nil, err
	}
	fmt.Println(booking)
	return booking, nil
}
func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	resp, err := s.coll.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.ID = resp.InsertedID.(bson.ObjectID)
	return booking, nil

}
