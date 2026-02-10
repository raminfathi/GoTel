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
	// 1. Init Dependencies
	mongoEndpoint := os.Getenv("MONGO_DB_URL")
	redisAddr := os.Getenv("REDIS_URL")
	redisPw := os.Getenv("REDIS_PASSWORD")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPw,
		DB:       0,
	})
	fmt.Println("Redis client initialized:", redisClient)

	client, err := mongo.Connect(options.Client().ApplyURI(mongoEndpoint))
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			log.Fatal(err)
		}
	}()

	// 2. Init Stores
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userStore := db.NewMongoUserStore(client)
	bookingStore := db.NewMongoBookingStore(client)
	cacheStore := db.NewRedisCacheStore(redisClient)

	store := &db.Store{
		Hotel:   hotelStore,
		Room:    roomStore,
		User:    userStore,
		Booking: bookingStore,
		Cache:   cacheStore,
	}

	// 3. Init Handlers
	hotelHandler := api.NewHotelHandler(store)
	authHandler := api.NewAuthHandler(userStore)
	userHandler := api.NewUserHandler(userStore)
	roomHandler := api.NewRoomHandler(store)
	bookingHandler := api.NewBookingHandler(store)

	// 4. Setup Fiber & Routes
	app := fiber.New(config)

	// authGroup
	authGroup := app.Group("/api")

	v1 := app.Group("/api/v1")

	apiv1 := v1.Group("/", middleware.JWTAuthentication(userStore))

	admin := apiv1.Group("/admin", api.AdminAuth)

	// ===========================
	// 🔓 Public Routes
	// ===========================
	authGroup.Post("/auth", authHandler.HandleAuthenticate)
	v1.Post("/user", userHandler.HandlePostUser)

	// ===========================
	// 🔒 Private Routes
	// ===========================

	// User Handlers
	apiv1.Get("/user", userHandler.HandleGetUsers)
	apiv1.Get("/user/:id", userHandler.HandleGetUser)
	apiv1.Put("/user/:id", userHandler.HandlePutUser)
	apiv1.Delete("/user/:id", userHandler.HandleDeleteUser)

	// Hotel Handlers
	apiv1.Get("/hotel", hotelHandler.HandleGetHotels)
	apiv1.Get("/hotel/:id", hotelHandler.HandleGetHotel)
	apiv1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRooms)

	// Room Handlers
	apiv1.Get("/room", roomHandler.HandleGetRooms)
	apiv1.Post("/room/:id/book", roomHandler.HandleBookRoom)

	// Booking Handlers
	apiv1.Get("/booking", bookingHandler.HandleGetMyBookings)
	apiv1.Get("/booking/:id", bookingHandler.HandleGetBooking)
	apiv1.Get("/booking/:id/cancel", bookingHandler.HandleCancelBooking)

	// ===========================
	// 👮 Admin Routes
	// ===========================
	admin.Post("/hotel", hotelHandler.HandlePostHotel)
	admin.Put("/hotel/:id", hotelHandler.HandlePutHotel)
	admin.Post("/room", roomHandler.HandlePostRoom)
	admin.Get("/booking", bookingHandler.HandleGetBookings)

	// Start Server
	listenAddr := os.Getenv("HTTP_LISTEN_ADDRESS")
	log.Fatal(app.Listen(listenAddr))
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on system environment variables")
	}
}
