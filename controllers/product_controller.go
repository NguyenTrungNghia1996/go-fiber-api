package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// POST /api/products
// Tạo một sản phẩm mới từ dữ liệu gửi lên trong body
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	// Phân tích dữ liệu JSON từ body vào struct Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid product data",
			Data:    nil,
		})
	}

	// Gọi repository để lưu sản phẩm vào database
	if err := repositories.CreateProduct(&product); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create product",
			Data:    nil,
		})
	}

	// Trả về thông tin sản phẩm đã được tạo thành công
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product created",
		Data:    product,
	})
}

// GET /api/products
// Lấy toàn bộ danh sách sản phẩm
func GetAllProducts(c *fiber.Ctx) error {
	// Gọi repository để lấy danh sách sản phẩm
	products, err := repositories.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to fetch products",
			Data:    nil,
		})
	}

	// Trả về danh sách sản phẩm
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product list retrieved",
		Data:    products,
	})
}

// GET /api/products?id=xxx
// Lấy thông tin sản phẩm theo ID truyền qua query string
func GetProductByID(c *fiber.Ctx) error {
	id := c.Query("id") // Lấy ID từ query string

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing product ID in query",
			Data:    nil,
		})
	}

	// Gọi repository để tìm sản phẩm theo ID
	product, err := repositories.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Product not found",
			Data:    nil,
		})
	}

	// Trả về thông tin sản phẩm nếu tìm thấy
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product retrieved",
		Data:    product,
	})
}

// PUT /api/products
// Cập nhật thông tin sản phẩm dựa trên dữ liệu gửi trong body (bao gồm cả ID)
func UpdateProduct(c *fiber.Ctx) error {
	var updateData map[string]interface{}

	// Phân tích body request thành map để xử lý động
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid update data",
			Data:    nil,
		})
	}

	// Lấy ID từ dữ liệu gửi lên
	idValue, ok := updateData["id"]
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing product ID in request body",
			Data:    nil,
		})
	}

	// Kiểm tra ID có phải dạng string hợp lệ hay không
	id, ok := idValue.(string)
	if !ok || id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid product ID",
			Data:    nil,
		})
	}

	// Xóa ID ra khỏi dữ liệu cập nhật để tránh ghi đè ID
	delete(updateData, "id")

	// Gọi repository để cập nhật thông tin sản phẩm
	if err := repositories.UpdateProduct(id, bson.M(updateData)); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update product",
			Data:    nil,
		})
	}

	// Trả về kết quả thành công
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product updated",
		Data:    nil,
	})
}

// DELETE /api/products?id=xxx
// Xóa sản phẩm dựa trên ID truyền qua query
func DeleteProduct(c *fiber.Ctx) error {
	id := c.Query("id") // Lấy ID từ query string

	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing product ID in query",
			Data:    nil,
		})
	}

	// Gọi repository để xóa sản phẩm
	if err := repositories.DeleteProduct(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to delete product",
			Data:    nil,
		})
	}

	// Trả về phản hồi thành công
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Product deleted",
		Data:    nil,
	})
}
