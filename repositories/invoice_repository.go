package repositories

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go-fiber-api/models"
)

type InvoiceRepository interface {
	CreateInvoice(ctx context.Context, invoice *models.Invoice) error
	GetInvoiceByID(ctx context.Context, id string) (*models.Invoice, error)
	ListInvoices(ctx context.Context) ([]models.Invoice, error)
	DeleteInvoice(ctx context.Context, id string) error
}

type invoiceRepository struct {
	collection *mongo.Collection
}

func NewInvoiceRepository(db *mongo.Database) InvoiceRepository {
	return &invoiceRepository{
		collection: db.Collection("invoices"),
	}
}

// CreateInvoice inserts a new invoice and calculates totals
func (r *invoiceRepository) CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	if invoice.ID == "" {
		return errors.New("invoice ID is required")
	}
	invoice.CreatedAt = time.Now()
	calculateInvoiceTotals(invoice)

	_, err := r.collection.InsertOne(ctx, invoice)
	return err
}

func (r *invoiceRepository) GetInvoiceByID(ctx context.Context, id string) (*models.Invoice, error) {
	var invoice models.Invoice
	err := r.collection.FindOne(ctx, bson.M{"_id": id}).Decode(&invoice)
	if err != nil {
		return nil, err
	}
	return &invoice, nil
}

func (r *invoiceRepository) ListInvoices(ctx context.Context) ([]models.Invoice, error) {
	cursor, err := r.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var invoices []models.Invoice
	for cursor.Next(ctx) {
		var invoice models.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			return nil, err
		}
		invoices = append(invoices, invoice)
	}
	return invoices, nil
}

func (r *invoiceRepository) DeleteInvoice(ctx context.Context, id string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"_id": id})
	return err
}

// Helper: Calculate totals before insert
func calculateInvoiceTotals(invoice *models.Invoice) {
	total := 0.0
	for i := range invoice.Items {
		item := &invoice.Items[i]
		item.TotalPrice = float64(item.Quantity) * item.UnitPrice
		total += item.TotalPrice
	}
	invoice.TotalAmount = total
}
