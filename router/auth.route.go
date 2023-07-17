package router

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/controller"
	"github.com/gocroot/gocroot/repository"
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(page *fiber.App) {
	collection := config.Ulbimongoconn
	userRepo := repository.NewUsersRepository(collection, "users")
	roleRepo := repository.NewRoleRepository(collection, "roles")
	authController := controller.NewAuthController(userRepo, roleRepo)

	var path = "/auth"

	page.Post(path+"/register", authController.Register)
	page.Post(path+"/login", authController.Login)
}
