package api

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/raminfathi/GoTel/db"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type testdb struct {
	client *mongo.Client
	store  *db.Store
}

func setup(t *testing.T) *testdb {
	if err := godotenv.Load("../.env"); err != nil {
		log.Println("Could not load .env file")
	}

	dbURI := os.Getenv("MONGO_DB_URL_TEST")
	if dbURI == "" {
		dbURI = "mongodb://localhost:27017"
	}

	testDBName := fmt.Sprintf("hotel_db_test_%d", time.Now().UnixNano())
	os.Setenv("DBNAME", testDBName)

	client, err := mongo.Connect(options.Client().ApplyURI(dbURI))
	if err != nil {
		t.Fatal(err)
	}

	client.Database(testDBName).Drop(context.TODO())

	hotelStore := db.NewMongoHotelStore(client)

	return &testdb{
		client: client,
		store: &db.Store{
			User:    db.NewMongoUserStore(client),
			Hotel:   hotelStore,
			Room:    db.NewMongoRoomStore(client, hotelStore),
			Booking: db.NewMongoBookingStore(client),
		},
	}
}
func (tdb *testdb) teardown(t *testing.T) {
	if err := tdb.client.Disconnect(context.TODO()); err != nil {
		t.Fatal(err)
	}
}
