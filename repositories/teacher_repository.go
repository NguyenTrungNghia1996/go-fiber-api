package repositories

import (
	"context"
	"time"

	"go-fiber-api/models"
	"go-fiber-api/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TeacherRepository struct {
	Collection *mongo.Collection
}

func NewTeacherRepository(db *mongo.Database) *TeacherRepository {
	return &TeacherRepository{
		Collection: db.Collection("teachers"),
	}
}

// Create a new teacher
func (r *TeacherRepository) Create(ctx context.Context, teacher *models.Teacher) error {
	now := time.Now()
	teacher.ID = primitive.NewObjectID()
	teacher.CreatedAt = now
	teacher.UpdatedAt = now
	teacher.NameNormalized = utils.NormalizeText(teacher.Name)

	_, err := r.Collection.InsertOne(ctx, teacher)
	return err
}

// Update only fields that are sent
func (r *TeacherRepository) Update(ctx context.Context, id primitive.ObjectID, updateData *models.Teacher) error {
	update := bson.M{}

	if updateData.Name != "" {
		update["name"] = updateData.Name
		update["name_normalized"] = utils.NormalizeText(updateData.Name)
	}
	if updateData.Email != "" {
		update["email"] = updateData.Email
	}
	if updateData.Phone != "" {
		update["phone"] = updateData.Phone
	}
	if updateData.Address != "" {
		update["address"] = updateData.Address
	}
	if updateData.AvatarURL != "" {
		update["avatar_url"] = updateData.AvatarURL
	}
	if updateData.DateOfBirth != nil {
		update["date_of_birth"] = updateData.DateOfBirth
	}
	if updateData.SubjectIDs != nil {
		update["subject_ids"] = updateData.SubjectIDs
	}
	// is_active luôn được cập nhật nếu được gửi lên
	update["is_active"] = updateData.IsActive
	update["updated_at"] = time.Now()

	_, err := r.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

// Get one by ID
func (r *TeacherRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Teacher, error) {
	var teacher models.Teacher
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&teacher)
	if err != nil {
		return nil, err
	}
	return &teacher, nil
}

// Get all (không phân trang)
func (r *TeacherRepository) GetAll(ctx context.Context) ([]*models.Teacher, error) {
	var teachers []*models.Teacher

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var t models.Teacher
		if err := cursor.Decode(&t); err != nil {
			return nil, err
		}
		teachers = append(teachers, &t)
	}
	return teachers, nil
}

// List with pagination, keyword search, sorting
func (r *TeacherRepository) List(ctx context.Context, page, limit int64, sortField, sortOrder, keyword string, isActive *bool, subjectIDs []primitive.ObjectID) ([]*models.Teacher, int64, error) {
	var teachers []*models.Teacher

	filter := bson.M{}
	if keyword != "" {
		filter["name_normalized"] = bson.M{
			"$regex":   utils.NormalizeText(keyword),
			"$options": "i",
		}
	}
	if isActive != nil {
		filter["is_active"] = *isActive
	}
	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
		findOptions.SetSkip((page - 1) * limit)
	}
	if len(subjectIDs) > 0 {
		filter["subject_ids"] = bson.M{
			"$in": subjectIDs,
		}
	}
	// Sort
	sort := bson.D{{Key: "created_at", Value: -1}}
	if sortField != "" {
		order := -1
		if sortOrder == "asc" {
			order = 1
		}
		sort = bson.D{{Key: sortField, Value: order}}
	}
	findOptions.SetSort(sort)

	cursor, err := r.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var t models.Teacher
		if err := cursor.Decode(&t); err != nil {
			return nil, 0, err
		}
		teachers = append(teachers, &t)
	}

	total, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return teachers, total, nil
}
