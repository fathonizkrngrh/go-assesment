package main

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/router"
	"github.com/gocroot/gocroot/url"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/whatsauth/whatsauth"
	"log"
)

func main() {
	go whatsauth.RunHub()
	site := fiber.New(config.Iteung)
	site.Use(cors.New(config.Cors))
	url.Web(site)
	router.AuthRoute(site)
	router.UserRoute(site)

	log.Fatal(site.Listen("0.0.0.0:8000"))

}
