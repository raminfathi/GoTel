package api

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/raminfathi/GoTel/db"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type testdb struct {
	client *mongo.Client
	*db.Store
}

func (tdb *testdb) teardown(t *testing.T) {
	dbname := os.Getenv(db.MongoDBNameEnvName)

	if err := tdb.client.Database(dbname).Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}


}
func setup(t *testing.T) *testdb {
	_ = godotenv.Load("../.env")

	dburi := os.Getenv("MONGO_DB_URL_TEST")
	client, err := mongo.Connect(options.Client().ApplyURI(dburi))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	hotelStore := db.NewMongoHotelStore(client)
	return &testdb{
		client: client,
		Store: &db.Store{
			Hotel:   hotelStore,
			User:    db.NewMongoUserStore(client),
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
