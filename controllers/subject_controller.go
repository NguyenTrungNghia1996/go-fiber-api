package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SubjectController xử lý các logic liên quan đến môn học
type SubjectController struct {
	Repo *repositories.SubjectRepository
}

// NewSubjectController tạo một controller mới cho Subject
func NewSubjectController(repo *repositories.SubjectRepository) *SubjectController {
	return &SubjectController{Repo: repo}
}

// CreateSubject tạo mới một môn học
// POST /api/subjects
// Body:
//
//	{
//	  "name": "Toán học",
//	  "code": "MATH101",
//	  "description": "Môn Toán cơ bản",
//	  "is_active": true
//	}
func (s *SubjectController) CreateSubject(c *fiber.Ctx) error {
	var input models.Subject
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if input.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Name is required",
			Data:    nil,
		})
	}

	err := s.Repo.Create(c.Context(), &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create subject",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Subject created successfully",
		Data:    input,
	})
}

// UpdateSubject cập nhật thông tin môn học theo ID trong body
// PUT /api/subjects
// Body:
//
//	{
//	  "id": "665e1b3fa6ef0c2d7e3e594f",
//	  "name": "Toán nâng cao",
//	  "code": "MATH201",
//	  "description": "Môn học về toán nâng cao",
//	  "is_active": true
//	}
func (s *SubjectController) UpdateSubject(c *fiber.Ctx) error {
	var input models.Subject
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if input.ID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID is required",
			Data:    nil,
		})
	}

	err := s.Repo.Update(c.Context(), input.ID, &input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update subject",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Subject updated successfully",
		Data:    nil,
	})
}

// GetSubjectByID lấy thông tin môn học theo ID truyền qua query
// GET /api/subjects?id=665e1b3fa6ef0c2d7e3e594f
func (s *SubjectController) GetSubjectByID(c *fiber.Ctx) error {
	idStr := c.Query("id")
	if idStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID is required",
			Data:    nil,
		})
	}

	id, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
	}

	subject, err := s.Repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Subject not found",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Subject fetched successfully",
		Data:    subject,
	})
}

// ListSubjects trả về danh sách môn học với phân trang và tìm kiếm
// GET /api/subjects/list?page=1&limit=10&sort_field=name&sort_order=asc&keyword=toan
func (s *SubjectController) ListSubjects(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	sortField := c.Query("sort_field", "")
	sortOrder := c.Query("sort_order", "desc")
	keyword := c.Query("keyword", "")
	isActiveStr := c.Query("is_active", "")

	var isActive *bool
	if isActiveStr == "true" {
		val := true
		isActive = &val
	} else if isActiveStr == "false" {
		val := false
		isActive = &val
	}

	subjects, total, err := s.Repo.List(c.Context(), int64(page), int64(limit), sortField, sortOrder, keyword, isActive)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to list subjects",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Subjects listed successfully",
		Data: fiber.Map{
			"items": subjects,
			"total": total,
		},
	})
}

// DeleteSubject handles
// DELETE /subjects?id=..
func (sc *SubjectController) DeleteSubject(c *fiber.Ctx) error {
	idParam := c.Query("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID is required",
			Data:    nil,
		})
	}
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid subject ID",
			Data:    nil,
		})
	}

	if err := sc.Repo.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to delete subject",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Subject deleted successfully",
		Data:    nil,
	})
}
