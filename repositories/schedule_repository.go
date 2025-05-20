package repositories

import (
	"context"
	"errors"
	"time"

	"go-fiber-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ScheduleRepository struct {
	Collection *mongo.Collection
}

func NewScheduleRepository(db *mongo.Database) *ScheduleRepository {
	return &ScheduleRepository{
		Collection: db.Collection("schedules"),
	}
}

// Create a new Schedule
func (r *ScheduleRepository) Create(ctx context.Context, schedule *models.Schedule) error {
	now := time.Now()
	schedule.ID = primitive.NewObjectID()
	schedule.CreatedAt = now
	schedule.UpdatedAt = now

	_, err := r.Collection.InsertOne(ctx, schedule)
	return err
}

// Update Schedule by ID, chỉ update các trường được gửi, giữ nguyên các trường còn lại
func (r *ScheduleRepository) Update(ctx context.Context, id primitive.ObjectID, updateData *models.Schedule) error {
	update := bson.M{}
	if updateData.ClassroomID != primitive.NilObjectID {
		update["classroom_id"] = updateData.ClassroomID
	}
	if updateData.AcademicYear != "" {
		update["academic_year"] = updateData.AcademicYear
	}
	if updateData.Semester != 0 {
		update["semester"] = updateData.Semester
	}
	if updateData.Week != 0 {
		update["week"] = updateData.Week
	}
	if updateData.Days != nil {
		update["days"] = updateData.Days
	}

	if len(update) == 0 {
		return errors.New("no data to update")
	}

	update["updated_at"] = time.Now()

	_, err := r.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

// Delete Schedule by ID
func (r *ScheduleRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	if res.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}
	return nil
}

// Get Schedule by ID
func (r *ScheduleRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Schedule, error) {
	var schedule models.Schedule
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&schedule)
	if err != nil {
		return nil, err
	}
	return &schedule, nil
}

// Find Schedule by classroomID, academicYear, semester, week (exact match)
func (r *ScheduleRepository) FindByClassroomWeek(ctx context.Context, classroomID primitive.ObjectID, academicYear string, semester, week int) ([]*models.Schedule, error) {
	filter := bson.M{
		"classroom_id":  classroomID,
		"academic_year": academicYear,
		"semester":      semester,
		"week":          week,
	}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var schedules []*models.Schedule
	for cursor.Next(ctx) {
		var s models.Schedule
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		schedules = append(schedules, &s)
	}
	return schedules, nil
}

// List schedules with pagination, sorting, and optional filters (classroom, academicYear, semester, week)
func (r *ScheduleRepository) List(ctx context.Context, page, limit int64, sortField, sortOrder string, filters map[string]interface{}) ([]*models.Schedule, int64, error) {
	filter := bson.M{}

	// Build filter from filters map
	for key, val := range filters {
		if val != nil {
			filter[key] = val
		}
	}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
		findOptions.SetSkip((page - 1) * limit)
	}

	// Sorting
	order := -1
	if sortOrder == "asc" {
		order = 1
	}
	if sortField == "" {
		sortField = "created_at"
	}
	findOptions.SetSort(bson.D{{Key: sortField, Value: order}})

	cursor, err := r.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var schedules []*models.Schedule
	for cursor.Next(ctx) {
		var s models.Schedule
		if err := cursor.Decode(&s); err != nil {
			return nil, 0, err
		}
		schedules = append(schedules, &s)
	}

	total, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return schedules, total, nil
}
