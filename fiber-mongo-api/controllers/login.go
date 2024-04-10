package controllers

import (
	"context"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	log.Printf("User Data Successfully Parsed")

	if user.Email == "" || user.Password == "" {
        return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Email and password are required"}})
    }
	existingUser := models.User{}
	if err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&existingUser); err != nil {
		if err == mongo.ErrNoDocuments {
			return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Invalid email or password"}})
		}
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Invalid email or password"}})
	}
	token, err := generateToken(existingUser.Id.Hex())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	log.Printf("User %s Successfully Logged In", user.Email)
	response := LoginResponse{
		OK: true,
		Data: LoginPayload{
			Token: token,
			User:  existingUser,
		},
	}
	return c.Status(http.StatusOK).JSON(response)
}
