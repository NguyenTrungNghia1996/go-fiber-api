package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
)

// POST /api/invoices
func CreateInvoice(c *fiber.Ctx) error {
	var invoice models.Invoice
	if err := c.BodyParser(&invoice); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid invoice data",
			Data:    nil,
		})
	}

	if err := repositories.CreateInvoice(&invoice); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create invoice",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Invoice created",
		Data:    invoice,
	})
}

// GET /api/invoices
func GetAllInvoices(c *fiber.Ctx) error {
	invoices, err := repositories.GetAllInvoices()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to fetch invoices",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Invoice list retrieved",
		Data:    invoices,
	})
}

// GET /api/invoices/:id
func GetInvoiceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	invoice, err := repositories.GetInvoiceByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invoice not found",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Invoice retrieved",
		Data:    invoice,
	})
}
