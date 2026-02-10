package db

import (
	"context"
	"os"

	"github.com/raminfathi/GoTel/types"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type RoomStore interface {
	InsertRoom(context.Context, *types.Room) (*types.Room, error)
	GetRooms(context.Context, bson.M) ([]*types.Room, error)
}

type MongoRoomStore struct {
	client *mongo.Client
	coll   *mongo.Collection

	HotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, HotelStore HotelStore) *MongoRoomStore {
	dbname := os.Getenv(MongoDBNameEnvName)
	if dbname == "" {
		dbname = "hotel_db"
	}

	return &MongoRoomStore{
		client:     client,
		coll:       client.Database(dbname).Collection("rooms"),
		HotelStore: HotelStore,
	}
}
func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*types.Room, error) {
	resp, err := s.coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*types.Room
	if err := resp.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
func (s *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	// 1. اتاق رو اینسرت کن
	resp, err := s.coll.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.ID = resp.InsertedID.(bson.ObjectID)

	if err := s.HotelStore.UpdateHotelsRooms(ctx, room.HotelID, room.ID); err != nil {
		return nil, err
	}

	return room, nil
}
