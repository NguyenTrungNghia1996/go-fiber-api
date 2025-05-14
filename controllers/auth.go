// package controllers

// import (
// 	"context"
// 	"log"
// 	"os"
// 	"time"

// 	"go-fiber-api/config"
// 	"go-fiber-api/models"
// 	"go-fiber-api/utils"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/golang-jwt/jwt/v5"
// 	"go.mongodb.org/mongo-driver/bson"
// )

// // Login handles user authentication and returns a JWT if successful.
// func Login(c *fiber.Ctx) error {
// 	// Parse input
// 	var input models.User
// 	if err := c.BodyParser(&input); err != nil {
// 		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
// 			Status:  "error",
// 			Message: "Invalid input",
// 			Data:    nil,
// 		})
// 	}

// 	// Find user by username
// 	var user models.User
// 	collection := config.DB.Collection("users")
// 	err := collection.FindOne(context.Background(), bson.M{"username": input.Username}).Decode(&user)
// 	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
// 		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
// 			Status:  "error",
// 			Message: "Invalid credentials",
// 			Data:    nil,
// 		})
// 	}

// 	// Generate JWT
// 	claims := jwt.MapClaims{
// 		"id":   user.ID,
// 		"role": user.Role,
// 		"pid":  user.PersonID,
// 		"exp":  time.Now().Add(12 * time.Hour).Unix(),
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

// 	jwtStr, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
// 	if err != nil {
// 		return c.Status(fiber.StatusInternalServerError).JSON(models.APIResponse{
// 			Status:  "error",
// 			Message: "Token creation failed",
// 			Data:    nil,
// 		})
// 	}

// 	// Log and respond
// 	log.Println("User logged in:", user.Username)

//		return c.JSON(models.APIResponse{
//			Status:  "success",
//			Message: "Login successful",
//			Data: fiber.Map{
//				"id":       user.ID,
//				"role":     user.Role,
//				"personID": user.PersonID,
//				"token":    jwtStr,
//			},
//		})
//	}
package controllers

import (
	"go-fiber-api/models"
	"go-fiber-api/repository"
	"go-fiber-api/utils"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid input",
			Data:    nil,
		})
	}

	user, err := repository.FindUserByUsername(input.Username)
	if err != nil || !utils.CheckPasswordHash(input.Password, user.Password) {
		return c.Status(fiber.StatusUnauthorized).JSON(models.APIResponse{
			Status:  "error",
			Message: "Invalid credentials",
			Data:    nil,
		})
	}

	token, _ := utils.GenerateJWT(user.ID, user.Role, user.PersonID)
	return c.JSON(models.APIResponse{
		Status:  "success",
		Message: "Login successful",
		Data: fiber.Map{
			"id":       user.ID,
			"role":     user.Role,
			"personID": user.PersonID,
			"token":    token,
		},
	})
}
