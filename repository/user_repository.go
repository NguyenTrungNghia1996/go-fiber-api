package repository

import (
	"context"
	"go-fiber-api/config"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
)

func FindUserByUsername(username string) (*models.User, error) {
	var user models.User
	err := config.DB.Collection("users").FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
