package controller

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/dto"
	"github.com/gocroot/gocroot/models"
	"github.com/gocroot/gocroot/repository"
	"github.com/gocroot/gocroot/utils"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type AuthController struct {
	userRepo *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewAuthController(userRepo *repository.UserRepository, roleRepo *repository.RoleRepository) *AuthController {
	return &AuthController{
		userRepo: userRepo,
		roleRepo: roleRepo,
	}
}

func (ac *AuthController) Register(c *fiber.Ctx) error {
	var registerDto dto.RegisterDTO
	var user models.User

	if err := c.BodyParser(&registerDto); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "Invalid request payload"})
	}

	existingUser, err := ac.userRepo.GetUserByEmail(registerDto.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed to check user existence"})
	}
	if existingUser != nil {
		return c.Status(http.StatusConflict).JSON(fiber.Map{"status": "CONFLICT", "message": "User with this email already exists"})
	}

	user.Username = registerDto.Username
	user.Email = registerDto.Email

	// Hash the password before saving it to the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed to hash password"})
	}
	user.Password = string(hashedPassword)

	roleUser, err := ac.roleRepo.GetRoleByName("USER")
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Role user not found"})
	}
	user.RoleID = roleUser.ID

	// Save the user to the database
	if err := ac.userRepo.Create(user); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed to create user"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{"status": "CREATED", "message": "User created successfully"})
}

func (ac *AuthController) Login(c *fiber.Ctx) error {
	var loginDto dto.LoginDTO
	err := c.BodyParser(&loginDto)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "Invalid request payload"})
	}

	user, err := ac.userRepo.GetUserByEmail(loginDto.Email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed to check user existence"})
	}
	if user == nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "Email not registered yet"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginDto.Password))
	if err != nil {
		log.Printf("%s signin failed: %v\n", user.Email, err.Error())
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"status": "BAD_REQUEST", "message": "Wrong password !!"})
	}

	tokenMaker, err := utils.NewPasetoMaker([]byte(config.PrivateKey))
	if err != nil {
		log.Println(" cannot create token maker: \n", user.Email, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed when generate token"})
	}

	token, err := tokenMaker.CreateToken(
		user.ID.Hex(),
		user.RoleID.Hex(),
	)
	log.Println(token)
	if err != nil {
		log.Println("cannot create token: %w", user.Email, err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed when generate token"})
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  "OK",
			"message": "Success logging in",
			"data":    token,
		})
}

func (ac *AuthController) GetUsers(c *fiber.Ctx) error {
	users, err := ac.userRepo.GetAllUser()

	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": err.Error()})
	}
	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  "OK",
			"message": "Success get all user",
			"data": fiber.Map{
				"users": users,
			},
		})
}

func (ac *AuthController) Me(c *fiber.Ctx) error {
	userId, _ := c.Locals("user_id").(string)

	log.Println(userId)

	if userId == "nil" {
		return c.
			Status(http.StatusNotFound).
			JSON(fiber.Map{"status": "NOT_FOUND", "message": "UserID not found"})
	}

	users, err := ac.userRepo.GetUserByID(userId)
	if err != nil {
		return c.
			Status(http.StatusInternalServerError).
			JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": err.Error()})
	}

	return c.
		Status(http.StatusOK).
		JSON(fiber.Map{
			"status":  "OK",
			"message": "Success get all user",
			"data": fiber.Map{
				"users": users,
			},
		})
}
