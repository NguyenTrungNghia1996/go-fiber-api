package repositories

import (
	"context"

	"go-fiber-api/config"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tìm user theo username
func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Tạo user mới và gắn với giáo viên (PersonID)
func CreateUser(user *models.User) error {
	user.ID = primitive.NewObjectID().Hex()
	_, err := config.DB.Collection("users").InsertOne(context.TODO(), user)
	return err
}

// Kiểm tra username đã tồn tại
func IsUsernameExists(username string) (bool, error) {
	count, err := config.DB.Collection("users").CountDocuments(context.TODO(), bson.M{"username": username})
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
