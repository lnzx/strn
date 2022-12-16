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

	app.Static("/", "./web/dist")

	app.Get("/healthz", func(c *fiber.Ctx) error {
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

	// fix vue history router
	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusFound).SendFile("./web/dist/index.html")
	})

	port := os.Getenv("PORT")
	if port == "443" {
		log.Fatal(app.ListenTLS(":"+port, "./cert.pem", "./pem.key"))
	} else {
		if port == "" {
			port = "8080"
		}
		log.Fatal(app.Listen(":" + port))
	}
}
