package controllers

import (
	"context"

	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
)

type InvoiceController struct {
	Repo repositories.InvoiceRepository
}

func NewInvoiceController(repo repositories.InvoiceRepository) *InvoiceController {
	return &InvoiceController{Repo: repo}
}

// CreateInvoice tạo hóa đơn mới
//
// POST /api/invoices
//
// Body:
//
//	{
//	  "id": "inv001",
//	  "items": [
//	    {"product_id": "p1", "quantity": 2, "unit_price": 100},
//	    {"product_id": "p2", "quantity": 3, "unit_price": 50}
//	  ]
//	}
func (c *InvoiceController) CreateInvoice(ctx *fiber.Ctx) error {
	var invoice models.Invoice
	if err := ctx.BodyParser(&invoice); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := c.Repo.CreateInvoice(context.Background(), &invoice); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(invoice)
}

// GetInvoiceByID lấy hóa đơn theo ID qua query (?id=...)
//
// GET /api/invoices/detail?id=inv001
func (c *InvoiceController) GetInvoiceByID(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing id in query"})
	}

	invoice, err := c.Repo.GetInvoiceByID(context.Background(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Invoice not found"})
	}

	return ctx.JSON(invoice)
}

// DeleteInvoice xóa hóa đơn theo ID qua query (?id=...)
//
// DELETE /api/invoices?id=inv001
func (c *InvoiceController) DeleteInvoice(ctx *fiber.Ctx) error {
	id := ctx.Query("id")
	if id == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Missing id in query"})
	}

	if err := c.Repo.DeleteInvoice(context.Background(), id); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Delete failed"})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

// ListInvoices trả về danh sách tất cả hóa đơn
//
// GET /api/invoices/list
func (c *InvoiceController) ListInvoices(ctx *fiber.Ctx) error {
	invoices, err := c.Repo.ListInvoices(context.Background())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Không thể lấy danh sách hóa đơn",
		})
	}

	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Danh sách hóa đơn",
		"data":    invoices,
	})
}
