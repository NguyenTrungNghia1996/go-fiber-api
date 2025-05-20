package controllers

import (
	"context"
	"time"

	"go-fiber-api/config"
	"go-fiber-api/models"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 📌 TẠO MỘT NGƯỜI MỚI
/*
@route   POST /api/persons
@body    {
  "name": "Nguyễn Văn A",
  "alias": "Ba A",
  "gender": "male",
	"image_url":"https://example.com/images/nguyen-van-a.jpg",
  "dob": "1980-05-10T00:00:00Z",
  "dod": null,
  "can_chi_dob": "Canh Thân",
  "can_chi_dod": "",
  "father_id": "664481c48fa7b11be59f53ad",
  "mother_id": null,
  "children_ids": ["664481c48fa7b11be59f53ae"],
  "spouse_ids": ["664481c48fa7b11be59f53af"]
}
*/
func CreatePerson(c *fiber.Ctx) error {
	var person models.Person
	if err := c.BodyParser(&person); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Data:    nil,
		})
	}

	repo := repositories.NewPersonRepository(config.DB)
	if err := repo.Create(context.TODO(), &person); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.Status(fiber.StatusCreated).JSON(models.APIResponse{
		Status:  "success",
		Message: "Person created successfully",
		Data:    person,
	})
}

// ✏️ CẬP NHẬT NGƯỜI THEO ID
/*
@route   PUT /api/persons?id=
@body    Giống như CreatePerson (JSON)
*/
func UpdatePerson(c *fiber.Ctx) error {
	var updateData map[string]interface{}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Data:    nil,
		})
	}

	idValue, ok := updateData["id"].(string)
	if !ok || idValue == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Missing or invalid 'id' field",
			Data:    nil,
		})
	}

	personID, err := primitive.ObjectIDFromHex(idValue)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    nil,
		})
	}

	repo := repositories.NewPersonRepository(config.DB)
	existingPerson, err := repo.GetByID(context.TODO(), personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}
	if existingPerson == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Person not found",
			Data:    nil,
		})
	}

	updatedPerson := *existingPerson

	// Cập nhật các trường từ body
	if name, ok := updateData["name"].(string); ok {
		updatedPerson.Name = name
	}
	if alias, ok := updateData["alias"].(string); ok {
		updatedPerson.Alias = alias
	}
	if gender, ok := updateData["gender"].(string); ok {
		updatedPerson.Gender = gender
	}
	if imageURL, ok := updateData["image_url"].(string); ok {
		updatedPerson.ImageURL = imageURL
	}
	if dob, ok := updateData["dob"].(string); ok {
		if t, err := time.Parse(time.RFC3339, dob); err == nil {
			updatedPerson.BirthDate = &t
		}
	}
	if dod, ok := updateData["dod"].(string); ok {
		if t, err := time.Parse(time.RFC3339, dod); err == nil {
			updatedPerson.DeathDate = &t
		}
	}
	if canChiDob, ok := updateData["birth_year_can_chi"].(string); ok {
		updatedPerson.BirthYearCanChi = canChiDob
	}
	if canChiDod, ok := updateData["death_year_can_chi"].(string); ok {
		updatedPerson.DeathYearCanChi = canChiDod
	}
	if fatherID, ok := updateData["father_id"].(string); ok {
		if objID, err := primitive.ObjectIDFromHex(fatherID); err == nil {
			updatedPerson.FatherID = &objID
		}
	}
	if motherID, ok := updateData["mother_id"].(string); ok {
		if objID, err := primitive.ObjectIDFromHex(motherID); err == nil {
			updatedPerson.MotherID = &objID
		}
	}
	if childrenIDs, ok := updateData["children_ids"].([]interface{}); ok {
		var ids []primitive.ObjectID
		for _, id := range childrenIDs {
			if strID, ok := id.(string); ok {
				if objID, err := primitive.ObjectIDFromHex(strID); err == nil {
					ids = append(ids, objID)
				}
			}
		}
		updatedPerson.ChildrenIDs = ids
	}
	if spouseIDs, ok := updateData["spouse_ids"].([]interface{}); ok {
		var ids []primitive.ObjectID
		for _, id := range spouseIDs {
			if strID, ok := id.(string); ok {
				if objID, err := primitive.ObjectIDFromHex(strID); err == nil {
					ids = append(ids, objID)
				}
			}
		}
		updatedPerson.SpouseIDs = ids
	}

	// Cập nhật thời gian chỉnh sửa
	updatedPerson.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	if err := repo.Update(context.TODO(), &updatedPerson); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Person updated successfully",
		Data:    updatedPerson,
	})
}

// ❌ XOÁ NGƯỜI THEO ID VÀ CÁC MỐI QUAN HỆ LIÊN QUAN
/*
@route   DELETE /api/persons?id=
@desc    Xóa người và cập nhật tất cả các mối quan hệ liên quan
*/
func DeletePerson(c *fiber.Ctx) error {
	idParam := c.Query("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID query parameter is required",
			Data:    nil,
		})
	}

	personID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid ID format",
			Data:    fiber.Map{"details": err.Error()},
		})
	}

	repo := repositories.NewPersonRepository(config.DB)
	err = repo.Delete(context.TODO(), personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to delete person",
			Data:    fiber.Map{"details": err.Error()},
		})
	}

	return c.Status(fiber.StatusOK).JSON(models.APIResponse{
		Status:  "success",
		Message: "Person deleted successfully",
		Data:    fiber.Map{"person_id": idParam},
	})
}

// 🔍 TÌM KIẾM THEO TÊN / BÍ DANH
/*
@route   GET /api/persons/search?keyword=van a&limit=10
@return  []models.Person
*/
func SearchPersons(c *fiber.Ctx) error {
	keyword := c.Query("keyword")
	if keyword == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Keyword is required",
			Data:    nil,
		})
	}

	limit := int64(10)
	repo := repositories.NewPersonRepository(config.DB)
	results, err := repo.SearchByNameOrAlias(context.TODO(), keyword, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Search results",
		Data:    results,
	})
}

// 👨‍👩‍👧‍👦 LẤY THÔNG TIN GIA ĐÌNH
/*
@route   GET /api/persons/family?id=
@return  repositories.FamilyInfo
*/
func GetFamilyInfo(c *fiber.Ctx) error {
	idParam := c.Query("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID query parameter is required",
			Data:    nil,
		})
	}
	personID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid ID",
			Data:    nil,
		})
	}

	repo := repositories.NewPersonRepository(config.DB)
	familyInfo, err := repo.GetFamilyInfo(context.TODO(), personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}

	if familyInfo == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Person not found",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Family information retrieved",
		Data:    familyInfo,
	})
}

// 🔍 LẤY THÔNG TIN CHI TIẾT NGƯỜI THEO ID
/*
@route   GET /api/persons?id=
@return  models.Person
*/
func GetPersonByID(c *fiber.Ctx) error {
	idParam := c.Query("id")
	if idParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "ID query parameter is required",
			Data:    nil,
		})
	}
	personID, err := primitive.ObjectIDFromHex(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid ID",
			Data:    nil,
		})
	}

	repo := repositories.NewPersonRepository(config.DB)
	person, err := repo.GetByID(context.TODO(), personID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: err.Error(),
			Data:    nil,
		})
	}

	if person == nil {
		return c.Status(fiber.StatusNotFound).JSON(models.APIResponse{
			Status:  "error",
			Message: "Person not found",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Person retrieved",
		Data:    person,
	})
}

// 📋 LẤY DANH SÁCH TẤT CẢ NGƯỜI
/*
@route   GET /api/persons/all
@return  []models.Person
*/
func GetAllPersons(c *fiber.Ctx) error {
	repo := repositories.NewPersonRepository(config.DB)
	limit := int64(100) // or any default value you want
	offset := int64(0)  // start from the beginning
	persons, err := repo.GetAll(context.TODO(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
			Status:  "error",
			Message: "Failed to fetch persons",
			Data:    nil,
		})
	}

	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "All persons retrieved",
		Data:    persons,
	})
}
