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

type ClassroomRepository struct {
	Collection *mongo.Collection
}

func NewClassroomRepository(db *mongo.Database) *ClassroomRepository {
	return &ClassroomRepository{
		Collection: db.Collection("classrooms"),
	}
}

// Create a new classroom
func (r *ClassroomRepository) Create(ctx context.Context, classroom *models.Classroom) error {
	now := time.Now()
	classroom.ID = primitive.NewObjectID()
	classroom.CreatedAt = now
	classroom.UpdatedAt = now
	classroom.NameNormalized = utils.NormalizeText(classroom.Name)

	_, err := r.Collection.InsertOne(ctx, classroom)
	return err
}

// Update only fields that are sent
func (r *ClassroomRepository) Update(ctx context.Context, id primitive.ObjectID, updateData *models.Classroom) error {
	update := bson.M{}

	if updateData.Name != "" {
		update["name"] = updateData.Name
		update["name_normalized"] = utils.NormalizeText(updateData.Name)
	}
	if updateData.Description != "" {
		update["description"] = updateData.Description
	}
	if updateData.Grade != 0 {
		update["grade"] = updateData.Grade
	}
	update["is_active"] = updateData.IsActive
	update["updated_at"] = time.Now()

	_, err := r.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

// Get one classroom by ID
func (r *ClassroomRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Classroom, error) {
	var classroom models.Classroom
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&classroom)
	if err != nil {
		return nil, err
	}
	return &classroom, nil
}

// Get all classrooms (no pagination)
func (r *ClassroomRepository) GetAll(ctx context.Context) ([]*models.Classroom, error) {
	var classrooms []*models.Classroom

	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var c models.Classroom
		if err := cursor.Decode(&c); err != nil {
			return nil, err
		}
		classrooms = append(classrooms, &c)
	}
	return classrooms, nil
}

// List with pagination, keyword search, sorting
func (r *ClassroomRepository) List(ctx context.Context, page, limit int64, sortField, sortOrder, keyword string, isActive *bool) ([]*models.Classroom, int64, error) {
	var classrooms []*models.Classroom

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
		var c models.Classroom
		if err := cursor.Decode(&c); err != nil {
			return nil, 0, err
		}
		classrooms = append(classrooms, &c)
	}

	total, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return classrooms, total, nil
}
