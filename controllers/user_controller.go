package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"
	"go-fiber-api/utils"

	"github.com/gofiber/fiber/v2"
)

// CreateUser handles the creation of a new user
// POST /api/users
// Body:
//
//	{
//	  "username": "teacher1",
//	  "password": "123456",
//	  "email": "teacher@example.com",
//	  "role": "member",
//	  "person_id": "665abc..."
//	}
func CreateUser(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Dữ liệu không hợp lệ",
			Data:    nil,
		})
	}

	exists, err := repositories.IsUsernameExists(user.Username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Lỗi kiểm tra username",
			Data:    nil,
		})
	}
	if exists {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Username đã tồn tại",
			Data:    nil,
		})
	}

	if err := repositories.CreateUser(&user); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Không thể tạo user",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Tạo user thành công",
		Data:    user,
	})
}

// GetUsersByRole retrieves users by their role
// GET /api/users?role=member
func GetUsersByRole(c *fiber.Ctx) error {
	role := c.Query("role")

	users, err := repositories.GetUsersByRole(role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Không thể lấy danh sách user",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Lấy danh sách user thành công",
		Data:    users,
	})
}

// UpdateUserPersonID updates a user's associated person_id
// PUT /api/users/person
// Body:
//
//	{
//	  "id": "665e1b3fa6ef0c2d7e3e594f",
//	  "person_id": "665e1cdbabc123..."
//	}
func UpdateUserPersonID(c *fiber.Ctx) error {
	var body struct {
		ID       string `json:"id"`
		PersonID string `json:"person_id"`
	}
	if err := c.BodyParser(&body); err != nil || body.ID == "" || body.PersonID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Dữ liệu không hợp lệ",
			Data:    nil,
		})
	}

	err := repositories.UpdateUserPersonID(body.ID, body.PersonID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Không thể cập nhật PersonID",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Đã cập nhật PersonID thành công",
		Data:    nil,
	})
}

// ChangeUserPassword updates a user's password after verifying the old one
// PUT /api/users/password
// Body:
//
//	{
//	  "id": "665e1b3fa6ef0c2d7e3e594f",
//	  "old_password": "admin123",
//	  "new_password": "123456"
//	}
func ChangeUserPassword(c *fiber.Ctx) error {
	var body struct {
		ID          string `json:"id"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
	}

	if err := c.BodyParser(&body); err != nil || body.ID == "" || body.OldPassword == "" || body.NewPassword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Dữ liệu không hợp lệ",
			Data:    nil,
		})
	}

	user, err := repositories.FindUserByID(body.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Không tìm thấy người dùng",
			Data:    nil,
		})
	}

	if !utils.CheckPasswordHash(body.OldPassword, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Status:  "error",
			Message: "Mật khẩu cũ không chính xác",
			Data:    nil,
		})
	}

	hashed, err := utils.HashPassword(body.NewPassword)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Không thể mã hóa mật khẩu",
			Data:    nil,
		})
	}

	err = repositories.UpdateUserPassword(body.ID, hashed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Không thể cập nhật mật khẩu",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Đổi mật khẩu thành công",
		Data:    nil,
	})
}
