package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/middleware"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(app *fiber.App, db *mongo.Database) {
	// Auth
	app.Post("/login", controllers.Login)

	// Protected API group
	api := app.Group("/api", middleware.Protected())

	// Upload URL
	api.Put("/presigned_url", controllers.GetUploadUrl)

	// Person routes (đã comment)
	// persons := api.Group("/persons")
	// persons.Post("/", controllers.CreatePerson)
	// persons.Get("/search", controllers.SearchPersons)
	// persons.Get("/family", controllers.GetFamilyInfo)
	// persons.Get("/", controllers.GetPersonByID)
	// persons.Put("/", controllers.UpdatePerson)
	// persons.Delete("/", controllers.DeletePerson)

	// User routes
	usersGroup := api.Group("/users")

	usersGroup.Post("/", controllers.CreateUser)                // Tạo user mới
	usersGroup.Get("/", controllers.GetUsersByRole)             // Lấy danh sách user theo role (?role=)
	usersGroup.Put("/person", controllers.UpdateUserPersonID)   // Cập nhật person_id cho user
	usersGroup.Put("/password", controllers.ChangeUserPassword) // Đổi mật khẩu (kiểm tra mật khẩu cũ)

	// Subject routes
	subjectController := controllers.NewSubjectController(repositories.NewSubjectRepository(db))

	subjectGroup := api.Group("/subjects")
	subjectGroup.Post("/", subjectController.CreateSubject)       // POST /api/subjects
	subjectGroup.Put("/", subjectController.UpdateSubject)        // PUT /api/subjects (cập nhật theo body có id)
	subjectGroup.Get("/detail", subjectController.GetSubjectByID) // GET /api/subjects/detail?id=...
	subjectGroup.Get("/", subjectController.ListSubjects)         // GET /api/subjects?page=...&limit=...
	subjectGroup.Delete("/", subjectController.DeleteSubject)     // DELETE /api/subjects?id=...

	// Teacher routes
	teacherController := controllers.NewTeacherController(repositories.NewTeacherRepository(db)) // Khởi tạo TeacherController với TeacherRepository (truyền database đã kết nối)
	teachers := api.Group("/teachers")
	teachers.Get("/", teacherController.ListTeachers)     // GET /api/teachers?page=1&limit=10&sort_field=name&sort_order=asc&keyword=toan&is_active=true&subject_ids=...
	teachers.Get("/detail", teacherController.GetTeacher) // GET /api/teachers/detail?id=...
	teachers.Post("/", teacherController.CreateTeacher)   // POST /api/teachers
	teachers.Put("/", teacherController.UpdateTeacher)    // PUT /api/teachers
}
