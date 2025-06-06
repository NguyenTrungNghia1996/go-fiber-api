package repositories

import (
	"context"

	"go-fiber-api/config"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tạo sản phẩm mới
func CreateProduct(product *models.Product) error {
	product.ID = primitive.NewObjectID().Hex()
	_, err := config.DB.Collection("products").InsertOne(context.TODO(), product)
	return err
}

// Lấy danh sách sản phẩm
func GetAllProducts() ([]models.Product, error) {
	cursor, err := config.DB.Collection("products").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var products []models.Product
	for cursor.Next(context.TODO()) {
		var product models.Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	return products, nil
}

// Lấy chi tiết sản phẩm theo ID
func GetProductByID(id string) (*models.Product, error) {
	var product models.Product
	err := config.DB.Collection("products").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&product)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

// Cập nhật sản phẩm
func UpdateProduct(id string, updatedData bson.M) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": updatedData}
	_, err := config.DB.Collection("products").UpdateOne(context.TODO(), filter, update)
	return err
}

// Xóa sản phẩm
func DeleteProduct(id string) error {
	_, err := config.DB.Collection("products").DeleteOne(context.TODO(), bson.M{"_id": id})
	return err
}
