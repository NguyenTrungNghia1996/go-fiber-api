package models

import "time"

type Invoice struct {
	ID          string        `json:"id" bson:"_id,omitempty"`
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	Items       []InvoiceItem `json:"items" bson:"items"`
	TotalAmount float64       `json:"total_amount" bson:"total_amount"` // tổng giá trị hóa đơn
}
