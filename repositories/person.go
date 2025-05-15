package repositories

import (
	"context"
	"strings"
	"sync"
	"time"
	"unicode"

	"go-fiber-api/models" // Thay bằng tên module của bạn

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type PersonRepository struct {
	Collection *mongo.Collection
}

func NewPersonRepository(db *mongo.Database) *PersonRepository {
	return &PersonRepository{
		Collection: db.Collection("persons"),
	}
}

// ====== TẠO ======
func (r *PersonRepository) Create(ctx context.Context, person *models.Person) error {
	now := primitive.NewDateTimeFromTime(time.Now())
	person.CreatedAt = now
	person.UpdatedAt = now

	res, err := r.Collection.InsertOne(ctx, person)
	if err != nil {
		return err
	}

	person.ID = res.InsertedID.(primitive.ObjectID)
	return r.syncRelationships(ctx, person)
}

// ====== LẤY THÔNG TIN THEO ID ======
func (r *PersonRepository) GetByID(ctx context.Context, id primitive.ObjectID) (*models.Person, error) {
	var person models.Person

	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil // Không tìm thấy nhưng không phải lỗi
		}
		return nil, err
	}

	return &person, nil
}
func (r *PersonRepository) GetChildren(ctx context.Context, parentID primitive.ObjectID) ([]models.Person, error) {
	filter := bson.M{"$or": []bson.M{
		{"father_id": parentID},
		{"mother_id": parentID},
	}}

	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var children []models.Person
	if err := cursor.All(ctx, &children); err != nil {
		return nil, err
	}

	return children, nil
}
func (r *PersonRepository) GetSpouses(ctx context.Context, personID primitive.ObjectID) ([]models.Person, error) {
	var person models.Person
	err := r.Collection.FindOne(ctx, bson.M{"_id": personID}).Decode(&person)
	if err != nil {
		return nil, err
	}

	if len(person.SpouseIDs) == 0 {
		return []models.Person{}, nil
	}

	filter := bson.M{"_id": bson.M{"$in": person.SpouseIDs}}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var spouses []models.Person
	if err := cursor.All(ctx, &spouses); err != nil {
		return nil, err
	}

	return spouses, nil
}

type FamilyInfo struct {
	Person   *models.Person
	Father   *models.Person
	Mother   *models.Person
	Spouses  []models.Person
	Children []models.Person
}

// ====== LẤY THÔNG GIA ĐÌNH ======
func (r *PersonRepository) GetFamilyInfo(ctx context.Context, personID primitive.ObjectID) (*FamilyInfo, error) {
	// Lấy thông tin chính
	person, err := r.GetByID(ctx, personID)
	if err != nil || person == nil {
		return nil, err
	}

	var father, mother *models.Person
	var spouses, children []models.Person

	// Lấy thông tin cha mẹ song song
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		if person.FatherID != nil {
			father, _ = r.GetByID(ctx, *person.FatherID)
		}
	}()

	go func() {
		defer wg.Done()
		if person.MotherID != nil {
			mother, _ = r.GetByID(ctx, *person.MotherID)
		}
	}()

	// Lấy thông tin vợ/chồng và con cái
	spouses, _ = r.GetSpouses(ctx, personID)
	children, _ = r.GetChildren(ctx, personID)

	wg.Wait()

	return &FamilyInfo{
		Person:   person,
		Father:   father,
		Mother:   mother,
		Spouses:  spouses,
		Children: children,
	}, nil
}

// ====== CẬP NHẬT ======
func (r *PersonRepository) Update(ctx context.Context, person *models.Person) error {
	person.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	filter := bson.M{"_id": person.ID}
	update := bson.M{"$set": person}

	_, err := r.Collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	return r.syncRelationships(ctx, person)
}

// ====== XÓA NGƯỜI VÀ CẬP NHẬT CÁC MỐI QUAN HỆ LIÊN QUAN ======
func (r *PersonRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	// Lấy thông tin người cần xóa trước khi xóa
	var person models.Person
	err := r.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(&person)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil // Không tìm thấy -> coi như đã xóa
		}
		return err
	}

	// Xử lý trong transaction để đảm bảo tính toàn vẹn
	session, err := r.Collection.Database().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (interface{}, error) {
		// 1. Xóa người này khỏi danh sách con của cha/mẹ
		if person.FatherID != nil {
			_, err = r.Collection.UpdateOne(
				sessCtx,
				bson.M{"_id": *person.FatherID},
				bson.M{"$pull": bson.M{"children_ids": id}},
			)
			if err != nil {
				return nil, err
			}
		}

		if person.MotherID != nil {
			_, err = r.Collection.UpdateOne(
				sessCtx,
				bson.M{"_id": *person.MotherID},
				bson.M{"$pull": bson.M{"children_ids": id}},
			)
			if err != nil {
				return nil, err
			}
		}

		// 2. Xóa người này khỏi danh sách vợ/chồng của các spouse
		if len(person.SpouseIDs) > 0 {
			_, err = r.Collection.UpdateMany(
				sessCtx,
				bson.M{"_id": bson.M{"$in": person.SpouseIDs}},
				bson.M{"$pull": bson.M{"spouse_ids": id}},
			)
			if err != nil {
				return nil, err
			}
		}

		// 3. Xóa thông tin cha/mẹ của các con
		if len(person.ChildrenIDs) > 0 {
			update := bson.M{}
			if person.Gender == "male" {
				update = bson.M{"$set": bson.M{"father_id": nil}}
			} else if person.Gender == "female" {
				update = bson.M{"$set": bson.M{"mother_id": nil}}
			}

			_, err = r.Collection.UpdateMany(
				sessCtx,
				bson.M{"_id": bson.M{"$in": person.ChildrenIDs}},
				update,
			)
			if err != nil {
				return nil, err
			}
		}

		// 4. Cuối cùng mới xóa bản thân người đó
		_, err = r.Collection.DeleteOne(sessCtx, bson.M{"_id": id})
		return nil, err
	})

	return err
}

// ====== TÌM KIẾM THEO TÊN / BÍ DANH KHÔNG DẤU ======
func (r *PersonRepository) SearchByNameOrAlias(ctx context.Context, keyword string, limit int64) ([]models.Person, error) {
	normalizedKeyword := normalizeText(keyword)

	filter := bson.M{
		"$or": []bson.M{
			{
				"name": bson.M{
					"$regex":   normalizedKeyword,
					"$options": "i",
				},
			},
			{
				"alias": bson.M{
					"$regex":   normalizedKeyword,
					"$options": "i",
				},
			},
		},
	}

	findOptions := options.Find()
	findOptions.SetLimit(limit)
	findOptions.SetSort(bson.M{"created_at": -1})

	cursor, err := r.Collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var persons []models.Person
	if err := cursor.All(ctx, &persons); err != nil {
		return nil, err
	}
	return persons, nil
}

// ====== ĐỒNG BỘ QUAN HỆ ======
func (r *PersonRepository) syncRelationships(ctx context.Context, person *models.Person) error {
	// Update Father
	if person.FatherID != nil {
		_, _ = r.Collection.UpdateOne(ctx, bson.M{"_id": *person.FatherID}, bson.M{
			"$addToSet": bson.M{"children_ids": person.ID},
		})
	}

	// Update Mother
	if person.MotherID != nil {
		_, _ = r.Collection.UpdateOne(ctx, bson.M{"_id": *person.MotherID}, bson.M{
			"$addToSet": bson.M{"children_ids": person.ID},
		})
	}

	// Update Children
	for _, childID := range person.ChildrenIDs {
		update := bson.M{}
		if person.Gender == "male" {
			update = bson.M{"$set": bson.M{"father_id": person.ID}}
		} else if person.Gender == "female" {
			update = bson.M{"$set": bson.M{"mother_id": person.ID}}
		}
		_, _ = r.Collection.UpdateOne(ctx, bson.M{"_id": childID}, update)
	}

	// Update Spouses (vợ/chồng)
	for _, spouseID := range person.SpouseIDs {
		_, _ = r.Collection.UpdateOne(ctx, bson.M{"_id": spouseID}, bson.M{
			"$addToSet": bson.M{"spouse_ids": person.ID},
		})
	}

	return nil
}

// ====== HỖ TRỢ TÌM KIẾM KHÔNG DẤU ======
func normalizeText(input string) string {
	input = strings.ToLower(input)
	var output []rune
	for _, r := range input {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		r = removeVietnameseAccent(r)
		output = append(output, r)
	}
	return string(output)
}

func removeVietnameseAccent(r rune) rune {
	accentMap := map[rune]rune{
		'à': 'a', 'á': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a',
		'â': 'a', 'ầ': 'a', 'ấ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
		'ă': 'a', 'ằ': 'a', 'ắ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
		'è': 'e', 'é': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e',
		'ê': 'e', 'ề': 'e', 'ế': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
		'ì': 'i', 'í': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
		'ò': 'o', 'ó': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o',
		'ô': 'o', 'ồ': 'o', 'ố': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
		'ơ': 'o', 'ờ': 'o', 'ớ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
		'ù': 'u', 'ú': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u',
		'ư': 'u', 'ừ': 'u', 'ứ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
		'ỳ': 'y', 'ý': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
		'đ': 'd',
	}
	if val, ok := accentMap[r]; ok {
		return val
	}
	return r
}
