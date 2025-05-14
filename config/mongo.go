package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Biến toàn cục để truy cập MongoDB
var DB *mongo.Database

// Hàm khởi tạo kết nối MongoDB
func ConnectDB() {
	// Lấy URI từ biến môi trường
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("MONGO_URI không được để trống trong .env")
	}

	// Tạo client MongoDB
	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatalf("Lỗi tạo MongoDB client: %v", err)
	}

	// Tạo context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Kết nối
	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Lỗi kết nối MongoDB: %v", err)
	}

	// Ping để đảm bảo kết nối thành công
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	// Gán Database
	mongoName := os.Getenv("MONGO_NAME")
	if mongoName == "" {
		log.Fatal("MONGO_NAME không được để trống trong .env")
	}
	DB = client.Database(mongoName)
	log.Println("✅ Kết nối MongoDB thành công")
}
