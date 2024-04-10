package controllers

import (
	"context"
	"fiber-mongo-api/models"
	"fiber-mongo-api/responses"
	"log"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	validationErr := validate.Struct(&user)
	if validationErr != nil {
		log.Printf("Validation error: %v", validationErr)
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": validationErr.Error()}})
	}
	log.Printf("User Data Successfully Parsed")
	hashedPassword, err := hashPassword(user.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	user.Password = hashedPassword
	result, err := userCollection.InsertOne(ctx, user)
	if err != nil || result == nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	token, err := generateToken(user.Id.Hex())
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	log.Printf("User %s Successfully Registered", user.Email)
	response := CreateUserResponse{
		OK: true,
		Data: CreateUserPayload{
			Token: token,
			User: models.User{
				Id:        user.Id,
				FirstName: user.FirstName,
				LastName:  user.LastName,
			},
		},
	}
	return c.Status(http.StatusCreated).JSON(response)
}
