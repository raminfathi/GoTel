package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/raminfathi/GoTel/db"
	"github.com/raminfathi/GoTel/db/fixtures"
	"github.com/raminfathi/GoTel/types"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using default values")
	}

	mongoURI := os.Getenv("MONGO_DB_URL")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}

	dbName := os.Getenv("MONGO_DB_NAME")
	if dbName == "" {
		dbName = "GoTel"
	}
	// 1. اتصال به دیتابیس
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatal(err)
	}
	ctx := context.Background()

	// 2. پاک کردن دیتابیس قدیمی
	fmt.Println("🧹 Dropping database...")
	if err := client.Database(dbName).Drop(ctx); err != nil {
		log.Fatal(err)
	}

	// 3. راه‌اندازی Store
	hotelStore := db.NewMongoHotelStore(client)
	store := &db.Store{
		User:    db.NewMongoUserStore(client),
		Hotel:   hotelStore,
		Room:    db.NewMongoRoomStore(client, hotelStore),
		Booking: db.NewMongoBookingStore(client),
	}

	// 4. ساخت هتل با استفاده از Fixture
	fmt.Println("🏨 Seeding Hotel...")
	hotel := fixtures.AddHotel(store, "Espinas Palace", "Tehran", 5, nil)
	fmt.Printf("   -> Created Hotel: %s\n", hotel.Name)

	// 5. ساخت اتاق‌ها
	fmt.Println("🛏️  Seeding Rooms...")
	firstRoom := fixtures.AddRoom(store, types.Single, 99.9, hotel.ID)
	fixtures.AddRoom(store, types.Double, 149.9, hotel.ID)
	fixtures.AddRoom(store, types.SeaView, 199.9, hotel.ID)
	fmt.Println("   -> Created 3 rooms")

	// 6. ساخت کاربر ادمین
	fmt.Println("👤 Seeding Users...")
	admin := fixtures.AddUser(store, "admin", "admin", true)
	printUserCredentials(admin)

	// 7. ساخت کاربر معمولی
	user := fixtures.AddUser(store, "user", "user", false)
	printUserCredentials(user)

	// 8. ساخت رزرو
	fmt.Println("📅 Seeding Booking...")
	fixtures.AddBooking(store, user.ID, firstRoom.ID, time.Now(), time.Now().AddDate(0, 0, 3))
	fmt.Printf("   -> Booking created for user %s in hotel %s\n", user.Email, hotel.Name)

	fmt.Println("---------------------------------------------------------")
	fmt.Println("✅ Seeding completed successfully!")
	fmt.Println("---------------------------------------------------------")
}

func printUserCredentials(u *types.User) {
	// تولید توکن برای نمایش
	token := generateToken(u)

	// تعیین نقش بر اساس IsAdmin (اصلاح شد)
	role := "User"
	if u.IsAdmin {
		role = "Admin"
	}

	fmt.Printf("\n   User: %s %s (%s)\n", u.FirstName, u.LastName, role)
	fmt.Printf("   Email: %s\n", u.Email)
	fmt.Printf("   Password: %s_%s\n", u.FirstName, u.LastName)
	fmt.Printf("   🔑 X-Api-Token: %s\n", token)
}

func generateToken(user *types.User) string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default_secret"
	}

	claims := jwt.MapClaims{
		"id":      user.ID.Hex(),
		"email":   user.Email,
		"expires": time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "ERROR_GENERATING_TOKEN"
	}
	return tokenStr
}
