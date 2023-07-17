package middleware

import (
	"fmt"
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/utils"
	"github.com/gofiber/fiber/v2"
	"log"
	"net/http"
	"strings"
)

func authenticationMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": "UNAUTHORIZED", "message": "Empty Bearer token "})
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	tokenMaker, err := utils.NewPasetoMaker([]byte(config.PrivateKey))
	if err != nil {
		log.Println(" cannot create token maker: ", err.Error())
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"status": "INTERNAL_SERVER_ERROR", "message": "Failed when generate token"})
	}

	payload, err := tokenMaker.ValidateToken(token)
	if err != nil {
		log.Println("unautthorized: ", err.Error())
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{"status": "UNAUTHORIZED", "message": "Invalid token "})
	}

	fmt.Println(payload)
	// Token is valid, proceed to the next middleware or route handler
	return c.Next()
}
