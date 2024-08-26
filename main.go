package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type User struct {
	ID          int
	AccountType string
	Email       string
}

type TimeStamp struct {
	Day       int
	Month     int
	Year      int
	Hour      int
	Minute    int
	Second    int
	StampType string
}

func main() {
	app := fiber.New()

	err := godotenv.Load((".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	users := []User{}

	app.Get("/", func(c *fiber.Ctx) error { // Get all users
		return c.Status(200).JSON(users)
	})

	app.Post("/api/users", func(c *fiber.Ctx) error { // Create user
		user := &User{}
		if err := c.BodyParser(user); err != nil {
			return err
		}
		if user.AccountType == "" || user.Email == "" {
			return c.Status(400).JSON(fiber.Map{"error": "User email and account type is required"})
		}

		user.ID = len(users) + 1
		users = append(users, *user)

		return c.Status(201).JSON(user)
	})

	app.Patch("/api/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, user := range users {
			if fmt.Sprint(user.ID) == id {
				users[i].AccountType = "employee"
				users[i].Email = "email@example.com"
				return c.Status(200).JSON(users[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	})

	app.Delete("/api/users/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		for i, user := range users {
			if fmt.Sprint(user.ID) == id {
				users = append(users[:i], users[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"success": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
