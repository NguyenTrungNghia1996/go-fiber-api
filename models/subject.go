package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Subject struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name           string             `bson:"name" json:"name"`                                           // Tên môn học đầy đủ
	Code           string             `bson:"code,omitempty" json:"code,omitempty"`                       // Mã môn học (ví dụ: "MATH101")
	NameNormalized string             `bson:"name_normalized,omitempty" json:"name_normalized,omitempty"` // Tên chuẩn hóa không dấu, viết thường
	Description    string             `bson:"description,omitempty" json:"description,omitempty"`         // Mô tả môn học
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt      time.Time          `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
	IsActive       bool               `bson:"is_active" json:"is_active"` // Trạng thái môn học có được sử dụng hay không
}
