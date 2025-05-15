package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Biến toàn cục để truy cập MongoDB
var DB *mongo.Database

// LoadDotEnv tìm và load file .env từ thư mục hiện tại lên các thư mục cha
func LoadDotEnv() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			// Tìm thấy file .env
			return godotenv.Load(envPath)
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break // Đã lên tới root folder rồi
		}
		dir = parent
	}

	return fmt.Errorf(".env not found in any parent directory")
}

func ConnectDB() {
	// Cố gắng load .env nếu có
	if err := LoadDotEnv(); err != nil {
		log.Println("⚠️ Không tìm thấy .env, sẽ dùng biến môi trường hệ thống nếu có")
	} else {
		log.Println("✅ Đã load file .env thành công")
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ MONGO_URI không được để trống")
	}

	mongoName := os.Getenv("MONGO_NAME")
	if mongoName == "" {
		log.Fatal("❌ MONGO_NAME không được để trống")
	}

	clientOpts := options.Client().ApplyURI(mongoURI)
	client, err := mongo.NewClient(clientOpts)
	if err != nil {
		log.Fatalf("Lỗi tạo MongoDB client: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("Lỗi kết nối MongoDB: %v", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Không thể ping MongoDB: %v", err)
	}

	DB = client.Database(mongoName)
	log.Println("✅ Kết nối MongoDB thành công")
}
