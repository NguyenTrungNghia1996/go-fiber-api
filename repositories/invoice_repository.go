package repositories

import (
	"context"
	"time"

	"go-fiber-api/config"
	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tạo hóa đơn mới
func CreateInvoice(invoice *models.Invoice) error {
	invoice.ID = primitive.NewObjectID().Hex()
	invoice.CreatedAt = time.Now()
	_, err := config.DB.Collection("invoices").InsertOne(context.TODO(), invoice)
	return err
}

// Lấy danh sách hóa đơn
func GetAllInvoices() ([]models.Invoice, error) {
	cursor, err := config.DB.Collection("invoices").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var invoices []models.Invoice
	for cursor.Next(context.TODO()) {
		var invoice models.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}

	return invoices, nil
}

// Lấy hóa đơn theo ID
func GetInvoiceByID(id string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := config.DB.Collection("invoices").FindOne(context.TODO(), bson.M{"_id": id}).Decode(&invoice)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}
