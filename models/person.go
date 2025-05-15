package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Person struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`

	// Tên đầy đủ (tên thật của người)
	// Ví dụ: "Nguyễn Văn A"
	Name string `bson:"name" json:"name"`

	// Bí danh hoặc tên gọi khác (nếu có)
	// Ví dụ: "Ba Lúa"
	Alias string `bson:"alias,omitempty" json:"alias,omitempty"`

	// Tên đã chuẩn hóa không dấu (để tìm kiếm)
	// Ví dụ: "nguyen van a"
	NameNormalized string `bson:"name_normalized,omitempty" json:"name_normalized,omitempty"`

	// Bí danh đã chuẩn hóa không dấu (để tìm kiếm)
	// Ví dụ: "ba lua"
	AliasNormalized string `bson:"alias_normalized,omitempty" json:"alias_normalized,omitempty"`

	// Giới tính: "male" (nam), "female" (nữ)
	// Ví dụ: "male"
	Gender string `bson:"gender" json:"gender"`

	// Ngày sinh theo Dương lịch, định dạng ISO 8601
	// Ví dụ: "1950-03-25T00:00:00Z"
	BirthDate *time.Time `bson:"birth_date,omitempty" json:"birth_date,omitempty"`

	// Năm sinh theo Can Chi (nếu biết)
	// Ví dụ: "Mậu Dần"
	BirthYearCanChi string `bson:"birth_year_can_chi,omitempty" json:"birth_year_can_chi,omitempty"`

	// Ngày mất theo Dương lịch, định dạng ISO 8601
	// Ví dụ: "2021-07-10T00:00:00Z"
	DeathDate *time.Time `bson:"death_date,omitempty" json:"death_date,omitempty"`

	// Năm mất theo Can Chi (nếu biết)
	// Ví dụ: "Tân Sửu"
	DeathYearCanChi string `bson:"death_year_can_chi,omitempty" json:"death_year_can_chi,omitempty"`

	// URL hình ảnh đại diện
	// Ví dụ: "https://example.com/images/nguyen-van-a.jpg"
	ImageURL string `bson:"image_url,omitempty" json:"image_url,omitempty"`

	// ID người cha (ObjectID MongoDB)
	FatherID *primitive.ObjectID `bson:"father_id,omitempty" json:"father_id,omitempty"`

	// ID người mẹ (ObjectID MongoDB)
	MotherID *primitive.ObjectID `bson:"mother_id,omitempty" json:"mother_id,omitempty"`

	// Danh sách ID vợ/chồng (có thể nhiều người)
	SpouseIDs []primitive.ObjectID `bson:"spouse_ids,omitempty" json:"spouse_ids,omitempty"`

	// Danh sách ID con cái
	ChildrenIDs []primitive.ObjectID `bson:"children_ids,omitempty" json:"children_ids,omitempty"`

	// Thời điểm tạo (do hệ thống tự thêm)
	CreatedAt primitive.DateTime `bson:"created_at" json:"created_at"`

	// Thời điểm cập nhật gần nhất
	UpdatedAt primitive.DateTime `bson:"updated_at" json:"updated_at"`
}
