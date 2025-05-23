package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Teacher struct {
	ID             primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name           string               `bson:"name" json:"name"`                       // Tên đầy đủ
	NameNormalized string               `bson:"name_normalized" json:"name_normalized"` // Tên đã chuẩn hóa (không dấu, viết thường)
	Email          string               `bson:"email,omitempty" json:"email,omitempty"`
	Phone          string               `bson:"phone,omitempty" json:"phone,omitempty"`
	DateOfBirth    *time.Time           `bson:"date_of_birth,omitempty" json:"date_of_birth,omitempty"`
	Address        string               `bson:"address,omitempty" json:"address,omitempty"`
	SubjectIDs     []primitive.ObjectID `bson:"subject_ids" json:"subject_ids"`
	AvatarURL      string               `bson:"avatar_url,omitempty" json:"avatar_url,omitempty"`
	CreatedAt      time.Time            `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      time.Time            `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsActive       bool                 `bson:"is_active" json:"is_active"`
}
