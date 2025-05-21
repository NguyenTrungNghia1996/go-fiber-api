package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Classroom struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`                                           // Tên lớp, ví dụ: "12A1"
	Grade          int                `bson:"grade,omitempty" json:"grade,omitempty"`                     // Khối lớp (ví dụ 10, 11, 12)
	Description    string             `bson:"description,omitempty" json:"description,omitempty"`         // Mô tả lớp học
	SchoolYear     string             `bson:"school_year,omitempty" json:"school_year,omitempty"`         // Niên khóa, ví dụ: "2024-2025"
	NameNormalized string             `bson:"name_normalized,omitempty" json:"name_normalized,omitempty"` // Tên chuẩn hóa không dấu, viết thường
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsActive       bool               `bson:"is_active" json:"is_active"` // Trạng thái lớp còn hoạt động hay không
}
