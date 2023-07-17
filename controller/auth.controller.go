package controller

import (
	"github.com/gocroot/gocroot/models"
	"github.com/gocroot/gocroot/repository"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type AuthController struct {
	userRepo *repository.UserRepository
}

func NewAuthController(userRepo *repository.UserRepository) *AuthController {
	return &AuthController{
		userRepo: userRepo,
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "Invalid request payload"})
	}

	existingUser, err := ac.userRepo.GetUserByEmail(user.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed to check user existence"})
	}
	if existingUser != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"status": "CONFLICT", "message": "User with this email already exists"})
	}

	// Hash the password before saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed to hash password"})
	}
	user.Password = string(hashedPassword)

	// Save the user to the database
	if err := ac.userRepo.Create(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed to create user"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"status": "CREATED", "message": "User created successfully"})
}
