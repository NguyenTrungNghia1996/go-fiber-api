package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// POST /api/products
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid product data",
			Data:    nil,
		})
	}

	if err := repositories.CreateProduct(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create product",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product created",
		Data:    product,
	})
}

// GET /api/products
func GetAllProducts(c *fiber.Ctx) error {
	products, err := repositories.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to fetch products",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product list retrieved",
		Data:    products,
	})
}

// GET /api/products/:id
func GetProductByID(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := repositories.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Product not found",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product retrieved",
		Data:    product,
	})
}

// PUT /api/products/:id
func UpdateProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid update data",
			Data:    nil,
		})
	}

	if err := repositories.UpdateProduct(id, bson.M(updateData)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update product",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product updated",
		Data:    nil,
	})
}

// DELETE /api/products/:id
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := repositories.DeleteProduct(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to delete product",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product deleted",
		Data:    nil,
	})
}
