package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type SchedulePeriod struct {
	Period      int                `bson:"period" json:"period"`                                 // Tiết số (1-5)
	SubjectID   primitive.ObjectID `bson:"subject_id" json:"subject_id"`                         // Môn học
	TeacherID   primitive.ObjectID `bson:"teacher_id" json:"teacher_id"`                         // Giáo viên dạy tiết đó
	ClassroomID primitive.ObjectID `bson:"classroom_id,omitempty" json:"classroom_id,omitempty"` // Lớp học (nếu tiết dành riêng cho lớp)
	StartTime   *time.Time         `bson:"start_time,omitempty" json:"start_time,omitempty"`     // Thời gian bắt đầu tiết học (nếu cần)
	EndTime     *time.Time         `bson:"end_time,omitempty" json:"end_time,omitempty"`         // Thời gian kết thúc tiết học (nếu cần)
	Note        string             `bson:"note,omitempty" json:"note,omitempty"`                 // Ghi chú (ví dụ: tiết thay, nghỉ...)
	IsActive    bool               `bson:"is_active" json:"is_active"`                           // Trạng thái tiết học (đang áp dụng hay không)
}
