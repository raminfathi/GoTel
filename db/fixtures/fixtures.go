package fixtures

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/types"
	"golang.org/x/crypto/bcrypt"

	"go.mongodb.org/mongo-driver/v2/bson"
)

func AddBooking(store *db.Store, uid, rid bson.ObjectID, from, till time.Time) *types.Booking {
	booking := &types.Booking{
		UserID:   uid,
		RoomID:   rid,
		FromDate: from,
		TillDate: till,
	}
	insertedBooking, err := store.Booking.InsertBooking(context.Background(), booking)
	if err != nil {
		log.Fatal(err)
	}
	return insertedBooking
}

func AddRoom(store *db.Store, roomType types.RoomType, basePrice float64, hid bson.ObjectID) *types.Room {
	room := &types.Room{
		Type:      roomType,
		BasePrice: basePrice,
		Price:     basePrice,
		HotelID:   hid,
	}
	insertedRoom, err := store.Room.InsertRoom(context.Background(), room)
	if err != nil {
		log.Fatal(err)
	}
	return insertedRoom
}

func AddHotel(store *db.Store, name string, loc string, rating int, rooms []bson.ObjectID) *types.Hotel {
	var roomIDS = rooms
	if rooms == nil {
		roomIDS = []bson.ObjectID{}
	}
	hotel := types.Hotel{
		Name:     name,
		Location: loc,
		Rooms:    roomIDS,
		Rating:   rating,
	}
	insertedHotel, err := store.Hotel.InsertHotel(context.TODO(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	return insertedHotel
}

func AddUser(store *db.Store, fn, ln string, admin bool) *types.User {

	pw := fmt.Sprintf("%s_%s", fn, ln)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal(err)
	}

	user := types.User{
		FirstName:         fn,
		LastName:          ln,
		Email:             fmt.Sprintf("%s@%s.com", fn, ln),
		IsAdmin:           admin,
		EncryptedPassword: string(hashedPassword),
		CreatedAt:         time.Now(),
	}

	insertedUser, err := store.User.InsertUser(context.TODO(), &user)
	if err != nil {
		log.Fatal(err)
	}
	return insertedUser
}
