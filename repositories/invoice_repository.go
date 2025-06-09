package repositories

import (
	"context"
	"errors"
	"fmt"
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

	GetInvoiceReportByDateRange(ctx context.Context, from, to time.Time) (*models.InvoiceReport, error)
	GetProductSalesByDateRange(ctx context.Context, from, to time.Time) ([]models.ProductSalesReport, error)
	GetSalesByPeriod(ctx context.Context, from, to time.Time, groupBy string) ([]models.SalesByPeriod, error)
}

type invoiceRepository struct {
	collection *mongo.Collection
}

func NewInvoiceRepository(db *mongo.Database) InvoiceRepository {
	return &invoiceRepository{
		collection: db.Collection("invoices"),
	}
}

// Tổng quan theo thời gian
func (r *invoiceRepository) GetInvoiceReportByDateRange(ctx context.Context, from, to time.Time) (*models.InvoiceReport, error) {
	filter := bson.M{
		"created_at": bson.M{"$gte": from, "$lte": to},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var report models.InvoiceReport
	for cursor.Next(ctx) {
		var invoice models.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			return nil, err
		}
		report.TotalInvoices++
		report.TotalAmount += invoice.TotalAmount
		for _, item := range invoice.Items {
			report.TotalProductUnits += item.Quantity
		}
	}
	return &report, nil
}

// Báo cáo theo từng sản phẩm
func (r *invoiceRepository) GetProductSalesByDateRange(ctx context.Context, from, to time.Time) ([]models.ProductSalesReport, error) {
	filter := bson.M{
		"created_at": bson.M{"$gte": from, "$lte": to},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	reportMap := map[string]*models.ProductSalesReport{}
	for cursor.Next(ctx) {
		var invoice models.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			return nil, err
		}
		for _, item := range invoice.Items {
			key := item.ProductID
			if _, exists := reportMap[key]; !exists {
				reportMap[key] = &models.ProductSalesReport{
					ProductID:   item.ProductID,
					ProductName: item.ProductName,
				}
			}
			reportMap[key].TotalSold += item.Quantity
			reportMap[key].Revenue += item.TotalPrice
		}
	}

	var results []models.ProductSalesReport
	for _, v := range reportMap {
		results = append(results, *v)
	}
	return results, nil
}

// Báo cáo theo ngày hoặc tháng
func (r *invoiceRepository) GetSalesByPeriod(ctx context.Context, from, to time.Time, groupBy string) ([]models.SalesByPeriod, error) {
	if groupBy != "day" && groupBy != "month" {
		return nil, errors.New("groupBy must be 'day' or 'month'")
	}

	filter := bson.M{
		"created_at": bson.M{"$gte": from, "$lte": to},
	}
	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	reportMap := map[string]*models.SalesByPeriod{}
	for cursor.Next(ctx) {
		var invoice models.Invoice
		if err := cursor.Decode(&invoice); err != nil {
			return nil, err
		}
		var key string
		if groupBy == "day" {
			key = invoice.CreatedAt.Format("2006-01-02")
		} else {
			key = invoice.CreatedAt.Format("2006-01")
		}
		if _, exists := reportMap[key]; !exists {
			reportMap[key] = &models.SalesByPeriod{Period: key}
		}
		reportMap[key].Revenue += invoice.TotalAmount
		for _, item := range invoice.Items {
			reportMap[key].Quantity += item.Quantity
		}
	}

	var results []models.SalesByPeriod
	for _, v := range reportMap {
		results = append(results, *v)
	}
	return results, nil
}

func (r *invoiceRepository) CreateInvoice(ctx context.Context, invoice *models.Invoice) error {
	if invoice.ID == "" {
		invoice.ID = fmt.Sprintf("%d", time.Now().UnixNano())
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

func calculateInvoiceTotals(invoice *models.Invoice) {
	total := 0.0
	for i := range invoice.Items {
		item := &invoice.Items[i]
		item.TotalPrice = float64(item.Quantity) * item.UnitPrice
		total += item.TotalPrice
	}
	invoice.TotalAmount = total
}
