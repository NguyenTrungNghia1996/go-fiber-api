package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ScheduleController xử lý các request liên quan đến lịch học
type ScheduleController struct {
	Repo *repositories.ScheduleRepository
}

// NewScheduleController khởi tạo ScheduleController mới
func NewScheduleController(repo *repositories.ScheduleRepository) *ScheduleController {
	return &ScheduleController{
		Repo: repo,
	}
}

// ListSchedules trả về danh sách lịch học với phân trang, sắp xếp, và filter tùy chọn
// GET /api/schedules?page=1&limit=10&sort_field=week&sort_order=asc&classroom_id=...&academic_year=...&semester=...&week=...
//
// Ví dụ:
// GET /api/schedules?page=1&limit=5&sort_field=week&sort_order=asc&classroom_id=6422d3bffcbf843a59fcd108&academic_year=2024-2025&semester=1&week=10
//
// Response:
//
//	{
//	  "status": "success",
//	  "message": "List of schedules",
//	  "data": {
//	    "items": [ ...list of schedules... ],
//	    "total": 100,
//	    "page": 1,
//	    "limit": 5
//	  }
//	}
func (sc *ScheduleController) ListSchedules(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	sortField := c.Query("sort_field", "created_at")
	sortOrder := c.Query("sort_order", "desc")

	filters := make(map[string]interface{})

	if classroomIDStr := c.Query("classroom_id", ""); classroomIDStr != "" {
		classroomID, err := primitive.ObjectIDFromHex(classroomIDStr)
		if err == nil {
			filters["classroom_id"] = classroomID
		}
	}

	if academicYear := c.Query("academic_year", ""); academicYear != "" {
		filters["academic_year"] = academicYear
	}

	if semester := c.QueryInt("semester", 0); semester != 0 {
		filters["semester"] = semester
	}

	if week := c.QueryInt("week", 0); week != 0 {
		filters["week"] = week
	}
	if isActiveStr := c.Query("is_active", ""); isActiveStr != "" {
		var isActive *bool
		if isActiveStr != "" {
			val := isActiveStr == "true"
			isActive = &val
		}
		filters["is_active"] = isActive
	}
	// is_active
	schedules, total, err := sc.Repo.List(c.Context(), int64(page), int64(limit), sortField, sortOrder, filters)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to list schedules",
			Data:    nil,
		})

	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "List of schedules",
		Data: fiber.Map{
			"items": schedules,
			"total": total,
			"page":  page,
			"limit": limit,
		},
	})
}

// GetSchedule trả về lịch học theo ID
// GET /api/schedules/detail?id=...
//
// Ví dụ:
// GET /api/schedules/detail?id=6422d3bffcbf843a59fcd109
//
// Response:
//
//	{
//	  "status": "success",
//	  "message": "Schedule found",
//	  "data": {
//	    "id": "6422d3bffcbf843a59fcd109",
//	    "classroom_id": "6422d3bffcbf843a59fcd108",
//	    "academic_year": "2024-2025",
//	    "semester": 1,
//	    "week": 10,
//	    "days": [
//	      {
//	        "day_of_week": 1,
//	        "morning": [
//	          {
//	            "period": 1,
//	            "subject_id": "...",
//	            "teacher_id": "...",
//	            "classroom_id": "...",
//	            "start_time": "2024-09-01T07:00:00Z",
//	            "end_time": "2024-09-01T07:45:00Z",
//	            "note": "",
//	            "is_active": true
//	          }
//	        ],
//	        "afternoon": []
//	      },
//	      ...
//	    ],
//	    "is_active": true,
//	    "created_at": "...",
//	    "updated_at": "..."
//	  }
//	}
func (sc *ScheduleController) GetSchedule(c *fiber.Ctx) error {
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
			Message: "Invalid schedule ID",
			Data:    nil,
		})
	}

	schedule, err := sc.Repo.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Schedule not found",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Schedule found",
		Data:    schedule,
	})
}

// CreateSchedule tạo mới lịch học
// POST /api/schedules
//
// Ví dụ request body:
//
//	{
//	  "classroom_id": "6422d3bffcbf843a59fcd108",
//	  "academic_year": "2024-2025",
//	  "semester": 1,
//	  "week": 10,
//	  "days": [
//	    {
//	      "day_of_week": 1,
//	      "morning": [
//	        {
//	          "period": 1,
//	          "subject_id": "6422d3bffcbf843a59fcd200",
//	          "teacher_id": "6422d3bffcbf843a59fcd300",
//	          "is_active": true
//	        }
//	      ],
//	      "afternoon": []
//	    }
//	  ],
//	  "is_active": true
//	}
//
// Response:
//
//	{
//	  "status": "success",
//	  "message": "Schedule created successfully",
//	  "data": { ...created schedule object... }
//	}
func (sc *ScheduleController) CreateSchedule(c *fiber.Ctx) error {
	var schedule models.Schedule
	if err := c.BodyParser(&schedule); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if err := sc.Repo.Create(c.Context(), &schedule); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to create schedule",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Schedule created successfully",
		Data:    schedule,
	})
}

// UpdateSchedule cập nhật lịch học theo ID gửi trong body JSON
// PUT /api/schedules
//
// Ví dụ request body:
//
//	{
//	  "id": "6422d3bffcbf843a59fcd109",
//	  "week": 11,
//	  "is_active": false,
//	  "days": [ ...cập nhật các ngày hoặc để nguyên nếu không muốn thay đổi... ]
//	}
//
// Response:
//
//	{
//	  "status": "success",
//	  "message": "Schedule updated successfully"
//	}
func (sc *ScheduleController) UpdateSchedule(c *fiber.Ctx) error {
	var updateData models.Schedule
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid request body",
			Data:    nil,
		})
	}

	if updateData.ID == primitive.NilObjectID {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID is required in request body",
			Data:    nil,
		})
	}

	if err := sc.Repo.Update(c.Context(), updateData.ID, &updateData); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to update schedule",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Schedule updated successfully",
		Data:    nil,
	})
}

// DeleteSchedule xóa lịch học theo ID
// DELETE /api/schedules?id=...
//
// Ví dụ:
// DELETE /api/schedules?id=6422d3bffcbf843a59fcd109
//
// Response:
//
//	{
//	  "status": "success",
//	  "message": "Schedule deleted successfully"
//	}
func (sc *ScheduleController) DeleteSchedule(c *fiber.Ctx) error {
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
			Message: "Invalid schedule ID",
			Data:    nil,
		})
	}

	if err := sc.Repo.Delete(c.Context(), id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to delete schedule",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Schedule deleted successfully",
		Data:    nil,
	})
}
