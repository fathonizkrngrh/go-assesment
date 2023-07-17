package main

import (
	"github.com/gocroot/gocroot/config"
	"github.com/gocroot/gocroot/seeder"
	"log"
)

func main() {
	err := seeder.SeedData(config.Ulbimongoconn)
	if err != nil {
		log.Fatal(err)
	}

}
