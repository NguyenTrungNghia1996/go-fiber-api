package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
)

// POST /api/invoices
// Tạo một hóa đơn mới từ dữ liệu trong body request
func CreateInvoice(c *fiber.Ctx) error {
	var invoice models.Invoice

	// Phân tích dữ liệu JSON từ body request vào struct Invoice
	if err := c.BodyParser(&invoice); err != nil {
		// Trả về lỗi 400 nếu dữ liệu không hợp lệ
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid invoice data",
			Data:    nil,
		})
	}

	// Gọi repository để tạo hóa đơn trong database
	if err := repositories.CreateInvoice(&invoice); err != nil {
		// Trả về lỗi 500 nếu quá trình lưu thất bại
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create invoice",
			Data:    nil,
		})
	}

	// Trả về kết quả thành công và hóa đơn vừa tạo
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Invoice created",
		Data:    invoice,
	})
}

// GET /api/invoices
// Truy vấn và trả về toàn bộ danh sách hóa đơn
func GetAllInvoices(c *fiber.Ctx) error {
	// Gọi repository để lấy danh sách hóa đơn
	invoices, err := repositories.GetAllInvoices()
	if err != nil {
		// Trả về lỗi nếu không lấy được dữ liệu
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to fetch invoices",
			Data:    nil,
		})
	}

	// Trả về danh sách hóa đơn
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Invoice list retrieved",
		Data:    invoices,
	})
}

// GET /api/invoices?id=xxx
// Lấy thông tin hóa đơn theo ID truyền qua query string
func GetInvoiceByID(c *fiber.Ctx) error {
	// Lấy ID từ query string
	id := c.Query("id")

	// Kiểm tra nếu không có ID được cung cấp
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing invoice ID in query",
			Data:    nil,
		})
	}

	// Gọi repository để truy vấn hóa đơn theo ID
	invoice, err := repositories.GetInvoiceByID(id)
	if err != nil {
		// Trả về lỗi nếu không tìm thấy hóa đơn
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invoice not found",
			Data:    nil,
		})
	}

	// Trả về dữ liệu hóa đơn nếu tìm thấy
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Invoice retrieved",
		Data:    invoice,
	})
}
