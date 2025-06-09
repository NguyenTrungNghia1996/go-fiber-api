package models

type InvoiceItem struct {
	ProductID   string  `json:"product_id" bson:"product_id"`
	Quantity    int     `json:"quantity" bson:"quantity"`
	UnitPrice   float64 `json:"unit_price" bson:"unit_price"`
	ProductName string  `json:"product_name" bson:"product_name"`
	TotalPrice  float64 `json:"total_price" bson:"total_price"`
}
