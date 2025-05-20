package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Lịch học của 1 lớp theo tuần
type Schedule struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	ClassroomID  primitive.ObjectID `bson:"classroom_id" json:"classroom_id"`   // Lớp áp dụng
	AcademicYear string             `bson:"academic_year" json:"academic_year"` // Năm học (VD: "2024-2025")
	Semester     int                `bson:"semester" json:"semester"`           // Học kỳ (1 hoặc 2)
	Week         int                `bson:"week" json:"week"`                   // Tuần học thứ mấy (1–40+)

	Days []ScheduleDay `bson:"days" json:"days"` // Lịch theo các ngày trong tuần (Thứ 2 -> CN)

	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

// Lịch học theo từng ngày
type ScheduleDay struct {
	DayOfWeek int `bson:"day_of_week" json:"day_of_week"` // 1 = Monday, 7 = Sunday

	Morning   []SchedulePeriod `bson:"morning,omitempty" json:"morning,omitempty"`     // Tiết sáng
	Afternoon []SchedulePeriod `bson:"afternoon,omitempty" json:"afternoon,omitempty"` // Tiết chiều
}
