package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/lnzx/strn/api"
	"github.com/lnzx/strn/cron"
)

func main() {
	app := fiber.New()
	app.Use(cors.New())

	app.Static("/", "./web")

	app.Get("/keepalive", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Get("/cron", func(c *fiber.Ctx) error {
		err := cron.Start()
		fmt.Println()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.SendString("ok")
	})

	app.Get("/api", func(c *fiber.Ctx) error {
		data, err := api.GetData()
		if err != nil {
			return c.SendString(err.Error())
		}
		return c.JSON(data)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// fix vue history router
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusFound).SendFile("./web/index.html")
	})

	log.Fatal(app.Listen(":" + port))
}
