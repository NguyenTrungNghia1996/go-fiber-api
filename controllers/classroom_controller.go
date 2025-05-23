package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ClassroomController struct {
	Repo *repositories.ClassroomRepository
}

func NewClassroomController(repo *repositories.ClassroomRepository) *ClassroomController {
	return &ClassroomController{
		Repo: repo,
	}
}

// CreateClassroom godoc
// @Summary Tạo lớp học mới
// @Description Thêm một lớp học mới vào cơ sở dữ liệu
// @Tags Classrooms
// @Accept json
// @Produce json
// @Param classroom body models.Classroom true "Thông tin lớp học"
// @Success 201 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Router /api/classrooms/create [post]
func (c *ClassroomController) CreateClassroom(ctx *fiber.Ctx) error {
	var classroom models.Classroom
	if err := ctx.BodyParser(&classroom); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if err := c.Repo.Create(ctx.Context(), &classroom); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create classroom",
			Data:    nil,
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Status:  "success",
		Message: "Classroom created successfully",
		Data:    classroom,
	})
}

// GetClassroomByID godoc
// @Summary Lấy thông tin lớp học theo ID
// @Description Truy xuất thông tin lớp học dựa trên ID được truyền qua query string
// @Tags Classrooms
// @Accept json
// @Produce json
// @Param id query string true "ID của lớp học (ObjectID)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Router /api/classrooms/get [get]
//
// 📌 Ví dụ gọi API:
// GET /api/classrooms/get?id=66547a1bfb3e401aadb45201
func (c *ClassroomController) GetClassroomByID(ctx *fiber.Ctx) error {
	idParam := ctx.Query("id")
	id, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid classroom ID",
			Data:    nil,
		})
	}

	classroom, err := c.Repo.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Classroom not found",
			Data:    nil,
		})
	}

	return ctx.JSON(models.APIResponse{
		Status:  "success",
		Message: "Classroom retrieved successfully",
		Data:    classroom,
	})
}

// UpdateClassroom godoc
// @Summary Cập nhật thông tin lớp học
// @Description Cập nhật các trường thông tin của lớp học, ID được truyền trong body JSON
// @Tags Classrooms
// @Accept json
// @Produce json
// @Param classroom body models.Classroom true "Dữ liệu lớp học cần cập nhật (phải bao gồm _id)"
// @Success 200 {object} models.APIResponse
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /api/classrooms/update [put]
//
// 📌 Ví dụ gọi API:
// PUT /api/classrooms/update
// Body JSON:
//
//	{
//	  "id": "66547a1bfb3e401aadb45201",
//	  "name": "12A2 Updated",
//	  "grade": 12,
//	  "description": "Lớp 12A2 nâng cao",
//	  "is_active": true
//	}
func (c *ClassroomController) UpdateClassroom(ctx *fiber.Ctx) error {
	var classroom models.Classroom
	if err := ctx.BodyParser(&classroom); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if classroom.ID.IsZero() {
		return ctx.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing classroom ID",
			Data:    nil,
		})
	}

	if err := c.Repo.Update(ctx.Context(), classroom.ID, &classroom); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update classroom",
			Data:    nil,
		})
	}

	return ctx.JSON(models.APIResponse{
		Status:  "success",
		Message: "Classroom updated successfully",
	})
}

// ListClassrooms godoc
// @Summary Danh sách lớp học có phân trang, tìm kiếm, lọc
// @Description Danh sách lớp học theo trang, từ khóa, trạng thái
// @Tags Classrooms
// @Accept json
// @Produce json
// @Param page query int false "Trang hiện tại" default(1)
// @Param limit query int false "Số bản ghi mỗi trang" default(10)
// @Param keyword query string false "Từ khóa tìm kiếm theo tên lớp"
// @Param sort_field query string false "Trường sắp xếp"
// @Param sort_order query string false "asc | desc" Enums(asc, desc)
// @Param is_active query bool false "Lọc theo trạng thái hoạt động"
// @Success 200 {object} models.APIResponse
// @Router /api/classrooms/list [get]
//
// 📌 Ví dụ gọi API:
// GET /api/classrooms/list?page=1&limit=5&keyword=12a&is_active=true&school_year=...
func (c *ClassroomController) ListClassrooms(ctx *fiber.Ctx) error {
	page, _ := strconv.ParseInt(ctx.Query("page", "1"), 10, 64)
	limit, _ := strconv.ParseInt(ctx.Query("limit", "10"), 10, 64)
	keyword := ctx.Query("keyword")
	sortField := ctx.Query("sort_field")
	sortOrder := ctx.Query("sort_order")
	schoolYear := ctx.Query("school_year")

	var isActive *bool
	if ctx.Query("is_active") != "" {
		val := ctx.Query("is_active") == "true"
		isActive = &val
	}

	classrooms, total, err := c.Repo.List(ctx.Context(), page, limit, sortField, sortOrder, keyword, isActive, schoolYear)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to list classrooms",
			Data:    nil,
		})
	}

	return ctx.JSON(models.APIResponse{
		Status:  "success",
		Message: "Classrooms retrieved successfully",
		Data: fiber.Map{
			"items": classrooms,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}
