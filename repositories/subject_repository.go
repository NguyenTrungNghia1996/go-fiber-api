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

type SubjectRepository struct {
	Collection *mongo.Collection
}

func NewSubjectRepository(db *mongo.Database) *SubjectRepository {
	return &SubjectRepository{
		Collection: db.Collection("subjects"),
	}
}

// Create a new subject
func (r *SubjectRepository) Create(ctx context.Context, subject *models.Subject) error {
	now := time.Now()
	subject.ID = primitive.NewObjectID()
	subject.CreatedAt = now
	subject.UpdatedAt = now
	subject.NameNormalized = utils.NormalizeText(subject.Name)

	_, err := r.Collection.InsertOne(ctx, subject)
	return err
}

// Update only sent fields
func (r *SubjectRepository) Update(ctx context.Context, id primitive.ObjectID, updateData *models.Subject) error {
	update := bson.M{}

	if updateData.Name != "" {
		update["name"] = updateData.Name
		update["name_normalized"] = utils.NormalizeText(updateData.Name)
	}
	if updateData.Code != "" {
		update["code"] = updateData.Code
	}
	if updateData.Description != "" {
		update["description"] = updateData.Description
	}
	update["updated_at"] = time.Now()

	// is_active must be explicitly set
	if updateData.IsActive || !updateData.IsActive {
		update["is_active"] = updateData.IsActive
	}

	if len(update) == 1 && update["updated_at"] != nil {
		// No real data to update
		return nil
	}

	_, err := r.Collection.UpdateByID(ctx, id, bson.M{"$set": update})
	return err
}

// Get by ID
func (r *SubjectRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Subject, error) {
	var subject models.Subject
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&subject)
	if err != nil {
		return nil, err
	}
	return &subject, nil
}

// List with pagination, sorting, filtering
func (r *SubjectRepository) List(ctx context.Context, page, limit int64, sortField, sortOrder, keyword string) ([]*models.Subject, int64, error) {
	var subjects []*models.Subject

	filter := bson.M{}
	if keyword != "" {
		filter["name_normalized"] = bson.M{
			"$regex":   utils.NormalizeText(keyword),
			"$options": "i",
		}
	}

	findOptions := options.Find()
	if limit > 0 {
		findOptions.SetLimit(limit)
		findOptions.SetSkip((page - 1) * limit)
	}

	// Default sort
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
		var subject models.Subject
		if err := cursor.Decode(&subject); err != nil {
			return nil, 0, err
		}
		subjects = append(subjects, &subject)
	}

	total, err := r.Collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	return subjects, total, nil
}
