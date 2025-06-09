package controllers

import (
	"context"
	"time"

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

func (c *InvoiceController) ListInvoices(ctx *fiber.Ctx) error {
	invoices, err := c.Repo.ListInvoices(context.Background())
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Cannot fetch invoice list",
		})
	}
	return ctx.JSON(fiber.Map{
		"status":  "success",
		"message": "Invoice list",
		"data":    invoices,
	})
}

// 📊 Báo cáo: tổng quan
func (c *InvoiceController) GetInvoiceSummary(ctx *fiber.Ctx) error {
	from, to, err := parseDateRange(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	report, err := c.Repo.GetInvoiceReportByDateRange(context.Background(), from, to)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(report)
}

// 📊 Báo cáo: theo sản phẩm
func (c *InvoiceController) GetProductSales(ctx *fiber.Ctx) error {
	from, to, err := parseDateRange(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	results, err := c.Repo.GetProductSalesByDateRange(context.Background(), from, to)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(results)
}

// 📊 Báo cáo: gộp theo ngày/tháng
func (c *InvoiceController) GetGroupedSales(ctx *fiber.Ctx) error {
	from, to, err := parseDateRange(ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	groupBy := ctx.Query("group", "day")
	results, err := c.Repo.GetSalesByPeriod(context.Background(), from, to, groupBy)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(results)
}

// 📅 Helper: parse ngày
func parseDateRange(ctx *fiber.Ctx) (time.Time, time.Time, error) {
	fromStr := ctx.Query("from")
	toStr := ctx.Query("to")
	if fromStr == "" || toStr == "" {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Missing from or to query param")
	}
	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Invalid from date")
	}
	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return time.Time{}, time.Time{}, fiber.NewError(fiber.StatusBadRequest, "Invalid to date")
	}
	return from, to, nil
}
