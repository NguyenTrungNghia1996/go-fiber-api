package models

// Tổng quan toàn bộ
type InvoiceReport struct {
	TotalInvoices     int     `json:"total_invoices"`
	TotalAmount       float64 `json:"total_amount"`
	TotalProductUnits int     `json:"total_product_units"`
}

// Doanh thu theo từng sản phẩm
type ProductSalesReport struct {
	ProductID   string  `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalSold   int     `json:"total_sold"` // Tổng số lượng bán
	Revenue     float64 `json:"revenue"`    // Tổng doanh thu
}

// Doanh thu theo ngày hoặc tháng
type SalesByPeriod struct {
	Period   string  `json:"period"` // ví dụ: "2024-06-06" hoặc "2024-06"
	Revenue  float64 `json:"revenue"`
	Quantity int     `json:"quantity"` // tổng sản phẩm bán ra
}
