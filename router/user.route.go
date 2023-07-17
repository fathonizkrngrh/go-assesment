package router

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/controller"
	"github.com/gocroot/gocroot/middleware"
	"github.com/gocroot/gocroot/repository"
	"github.com/gofiber/fiber/v2"
)

func UserRoute(page *fiber.App) {
	collection := config.Ulbimongoconn
	userRepo := repository.NewUsersRepository(collection, "users")
	roleRepo := repository.NewRoleRepository(collection, "roles")
	authController := controller.NewAuthController(userRepo, roleRepo)

	var path = "/user"

	page.Use(middleware.AuthenticationMiddleware())
	page.Get(path+"/all", authController.GetUsers)
	page.Get(path+"/me", authController.Me)
}
