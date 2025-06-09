package routes

import (
	"go-fiber-api/controllers"
	"go-fiber-api/repositories"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func Setup(app *fiber.App, db *mongo.Database) {
	// Auth
	app.Post("/login", controllers.Login)
	app.Get("/test", controllers.Hello)
	// Protected API group
	// api := app.Group("/api", middleware.Protected())
	api := app.Group("/api")
	api.Get("/test2", controllers.Hello)
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

	// Schedule routes
	scheduleController := controllers.NewScheduleController(repositories.NewScheduleRepository(db))
	schedules := api.Group("/schedules")
	schedules.Get("/", scheduleController.ListSchedules)     // GET /api/schedules?page=1&limit=10&sort=created_at&order=desc&classroom_id=xxx&academic_year=2024-2025&semester=1&week=3
	schedules.Get("/detail", scheduleController.GetSchedule) // GET /api/schedules/detail?id=xxx
	schedules.Post("/", scheduleController.CreateSchedule)   // POST /api/schedules
	schedules.Put("/", scheduleController.UpdateSchedule)    // PUT /api/schedules
	schedules.Delete("/", scheduleController.DeleteSchedule) // DELETE /api/schedules?id=xxx

	//Classroom routes
	classroomController := controllers.NewClassroomController(repositories.NewClassroomRepository(db))
	classroom := api.Group("/classroom")
	classroom.Get("/", classroomController.ListClassrooms)         //GET /api/classroom?page=1&limit=10&sort_field=name&sort_order=asc&keyword=12A1&is_active=true
	classroom.Get("/detail", classroomController.GetClassroomByID) // GET /api/classroom/detail?id=664c3179f5a36b935b674f9d
	classroom.Post("/", classroomController.CreateClassroom)       //POST /api/classroom
	classroom.Put("/", classroomController.UpdateClassroom)        // PUT /api/classroom

	// Product routes
	productsGroup := api.Group("/products")
	productsGroup.Post("/", controllers.CreateProduct)       // POST /api/products
	productsGroup.Get("/", controllers.GetAllProducts)       // GET /api/products
	productsGroup.Get("/detail", controllers.GetProductByID) // GET /api/products/detail?id=...
	productsGroup.Put("/", controllers.UpdateProduct)        // PUT /api/products (id nằm trong body)
	productsGroup.Delete("/", controllers.DeleteProduct)     // DELETE /api/products?id=...

	// Invoice routes
	invoiceController := controllers.NewInvoiceController(repositories.NewInvoiceRepository(db))
	invoices := api.Group("/invoices")

	invoices.Post("/", invoiceController.CreateInvoice)   // POST /api/invoices
	invoices.Get("/", invoiceController.GetInvoiceByID)   // GET /api/invoices?id=...
	invoices.Delete("/", invoiceController.DeleteInvoice) // DELETE /api/invoices?id=...
	invoices.Get("/list", invoiceController.ListInvoices) // GET /api/invoices/list

	// 📊 Báo cáo hóa đơn
	invoices.Get("/report/summary", invoiceController.GetInvoiceSummary) // GET /api/invoices/report/summary?from=YYYY-MM-DD&to=YYYY-MM-DD
	invoices.Get("/report/products", invoiceController.GetProductSales)  // GET /api/invoices/report/products?from=...&to=...
	invoices.Get("/report/grouped", invoiceController.GetGroupedSales)   // GET /api/invoices/report/grouped?from=...&to=...&group=day|month

}
