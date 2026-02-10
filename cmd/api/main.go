package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/raminfathi/GoTel/api"
	"github.com/raminfathi/GoTel/api/middleware"
	"github.com/raminfathi/GoTel/db"
	"github.com/redis/go-redis/v9"

	"github.com/gofiber/fiber/v3"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var config = fiber.Config{
	ErrorHandler: api.ErrorHandler,
}

func main() {
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	redisAddr := os.Getenv("REDIS_URL")
	redisPw := os.Getenv("REDIS_PASSWORD")
	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPw,
		DB:       0,
	})

	client, err := mongo.Connect(options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()
	var (
		hotelStore   = db.NewMongoHotelStore(client)
		roomStore    = db.NewMongoRoomStore(client, hotelStore)
		userStore    = db.NewMongoUserStore(client)
		bookingStore = db.NewMongoBookingStore(client)
		cacheStore   = db.NewRedisCacheStore(redisClient)
		store        = &db.Store{
			Hotel:   hotelStore,
			Room:    roomStore,
			User:    userStore,
			Booking: bookingStore,
			Cache:   cacheStore,
		}
		hotelHandler   = api.NewHotelHandler(store)
		authHandler    = api.NewAuthHandler(userStore)
		userHandler    = api.NewUserHandler(userStore)
		roomHandler    = api.NewRoomHandler(store)
		bookingHandler = api.NewBookingHandler(store)

		app  = fiber.New(config)
		auth = app.Group("/api/")

		apiv1 = app.Group("/api/v1", middleware.JWTAuthentication(userStore))
		admin = apiv1.Group("/admin", api.AdminAuth)
	)
	fmt.Println("Redis client initialized:", redisClient)
	// auth
	auth.Post("auth", authHandler.HandleAuthenticate)

	//Versioned API routes
	// user handlers
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)
	apiv1.Post("/user", userHandler.HandlePostUser)
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)

	// hotel handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)
	admin.Post("/hotel", hotelHandler.HandlePostHotel)
	admin.Put("/hotel/:id", hotelHandler.HandlePutHotel)
	// room handler
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)
	admin.Post("/room", roomHandler.HandlePostRoom)
	// booking handler
	apiv1.Get("/booking", bookingHandler.HandleGetMyBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)
	//admin handler
	admin.Get("/booking", bookingHandler.HandleGetBookings)
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	log.Fatal(app.Listen(listenAddr))
}
func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}
}
