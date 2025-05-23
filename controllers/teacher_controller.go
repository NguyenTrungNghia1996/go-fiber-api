package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TeacherController struct {
	Repo *repositories.TeacherRepository
}

func NewTeacherController(repo *repositories.TeacherRepository) *TeacherController {
	return &TeacherController{Repo: repo}
}

// ListTeachers trả về danh sách giáo viên có phân trang, tìm kiếm và lọc
// GET /api/teachers?page=1&limit=10&sort_field=name&sort_order=asc&keyword=toan&is_active=true&subject_ids=id1,id2
func (tc *TeacherController) ListTeachers(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	sortField := c.Query("sort_field", "")
	sortOrder := c.Query("sort_order", "desc")
	keyword := c.Query("keyword", "")
	isActiveStr := c.Query("is_active", "")
	subjectIDsStr := c.Query("subject_ids", "") // e.g. "id1,id2"

	var isActive *bool
	if isActiveStr != "" {
		val, err := strconv.ParseBool(isActiveStr)
		if err == nil {
			isActive = &val
		}
	}

	var subjectIDs []primitive.ObjectID
	if subjectIDsStr != "" {
		ids := strings.Split(subjectIDsStr, ",")
		for _, idStr := range ids {
			id, err := primitive.ObjectIDFromHex(idStr)
			if err == nil {
				subjectIDs = append(subjectIDs, id)
			}
		}
	}

	teachers, total, err := tc.Repo.List(c.Context(), int64(page), int64(limit), sortField, sortOrder, keyword, isActive, subjectIDs)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to list teachers",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "List of teachers",
		Data: fiber.Map{
			"items": teachers,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetTeacher trả về thông tin giáo viên theo ID
// GET /api/teachers?id=...
func (tc *TeacherController) GetTeacher(c *fiber.Ctx) error {
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
			Message: "Invalid teacher ID",
			Data:    nil,
		})
	}

	teacher, err := tc.Repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Teacher not found",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Teacher found",
		Data:    teacher,
	})
}

// CreateTeacher tạo mới giáo viên
// POST /api/teachers
// Body: { "name": "...", "email": "...", ... }
func (tc *TeacherController) CreateTeacher(c *fiber.Ctx) error {
	var teacher models.Teacher
	if err := c.BodyParser(&teacher); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if err := tc.Repo.Create(c.Context(), &teacher); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create teacher",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Teacher created successfully",
		Data:    teacher,
	})
}

// UpdateTeacher cập nhật thông tin giáo viên
// PUT /api/teachers
// Body: { "id": "...", "name": "...", ... }
func (tc *TeacherController) UpdateTeacher(c *fiber.Ctx) error {
	var updateData models.Teacher
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if updateData.ID.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID is required in body",
			Data:    nil,
		})
	}

	if err := tc.Repo.Update(c.Context(), updateData.ID, &updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update teacher",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Teacher updated successfully",
	})
}
