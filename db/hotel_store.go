package db

import (
	"context"
	"fmt"
	"os"

	"github.com/raminfathi/GoTel/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type HotelStore interface {
	InsertHotel(context.Context, *types.Hotel) (*types.Hotel, error)
	UpdateHotel(context.Context, Map, Map) error
	GetHotels(context.Context, Map, *Pagination) ([]*types.Hotel, error)
	GetHotelByID(context.Context, string) (*types.Hotel, error)
	UpdateHotelsRooms(context.Context, bson.ObjectID, bson.ObjectID) error
}

type MongoHotelStore struct {
	client *mongo.Client
	coll   *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	dbname := os.Getenv(MongoDBNameEnvName)

	return &MongoHotelStore{
		client: client,
		coll:   client.Database(dbname).Collection("hotels"),
	}
}
func (s *MongoHotelStore) UpdateHotelsRooms(ctx context.Context, hotelID bson.ObjectID, roomID bson.ObjectID) error {
	filter := bson.M{"_id": hotelID}

	update := bson.M{"$push": bson.M{"rooms": roomID}}

	_, err := s.coll.UpdateOne(ctx, filter, update)
	return err
}
func (s *MongoHotelStore) GetHotelByID(ctx context.Context, id string) (*types.Hotel, error) {
	fmt.Println(id)

	oid, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var hotel *types.Hotel
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&hotel)
	if err != nil {
		return nil, err
	}
	fmt.Println(hotel)
	return hotel, nil
}

func (s *MongoHotelStore) GetHotels(ctx context.Context, filter Map, pag *Pagination) ([]*types.Hotel, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var hotels []*types.Hotel
	if err := resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}
func (s *MongoHotelStore) UpdateHotel(ctx context.Context, filter Map, update Map) error {
	doc := bson.M{"$set": update}
	_, err := s.coll.UpdateOne(ctx, filter, doc)
	return err
}

func (s *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	resp, err := s.coll.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.ID = resp.InsertedID.(bson.ObjectID)
	return hotel, nil
}
