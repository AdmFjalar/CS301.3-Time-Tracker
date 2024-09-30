package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type User struct {
	ID          uint32
	AccountType uint8
	Email       string
}

type TimeStamp struct {
	StampType   uint8
	UserID      uint32
	TimeStampID uint32
	Year        int16
	Month       uint8
	Day         uint8
	Hour        uint8
	Minute      uint8
	Second      uint8
}

func main() {
	app := fiber.New()

	err := godotenv.Load((".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	users := []User{}
	timestamps := []TimeStamp{}

	// TIME STAMP API

	// GET ALL TIME STAMPS OF SPECIFIC USER

	app.Get("/api/users/:id/timestamps", func(c *fiber.Ctx) error {

		userTimestamps := []TimeStamp{}
		for _, timestamp := range timestamps { // Loops through all timestamps, ignores index and accesses value in the timestamps array with timestamp
			if fmt.Sprint(timestamp.UserID) == c.Params("id") {
				userTimestamps = append(userTimestamps, timestamp)
			}
		}

		if len(userTimestamps) > 0 {
			return c.Status(200).JSON(userTimestamps)
		}

		return c.Status(404).JSON(fiber.Map{"error": "User not found or no timestamps registered to user"})
	})

	// CREATE TIME STAMP

	app.Post("/api/users/:id/timestamps/:type", func(c *fiber.Ctx) error {

		timestamp := TimeStamp{}
		if err := c.BodyParser(&timestamp); err != nil {
			return err
		}

		timestamp.Year = int16(time.Now().Year())
		timestamp.Month = uint8(time.Now().Month())
		timestamp.Day = uint8(time.Now().Day())
		timestamp.Hour = uint8(time.Now().Hour())
		timestamp.Minute = uint8(time.Now().Minute())
		timestamp.Second = uint8(time.Now().Second())

		var stampType int
		stampType, err = strconv.Atoi(c.Params("type"))
		timestamp.StampType = uint8(stampType)

		if timestamp.Day == 0 || timestamp.Month == 0 || timestamp.Year == 0 || timestamp.Hour == 0 || timestamp.Minute == 0 || timestamp.Second == 0 || timestamp.StampType == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Day, month, year, hour, minute, second and stamp type are required"})
		}

		var userID int
		userID, err = strconv.Atoi(c.Params("id"))
		timestamp.UserID = uint32(userID)

		timestamp.TimeStampID = uint32(len(timestamps) + 1)
		timestamps = append(timestamps, timestamp)

		return c.Status(201).JSON(timestamp)
	})

	// USER API

	app.Get("/api/users/:id", func(c *fiber.Ctx) error { // Get user
		id := c.Params("id")
		for _, user := range users {
			if fmt.Sprint(user.ID) == id {
				return c.Status(200).JSON(user)
			}
		}
		return c.Status(404).JSON(fiber.Map{"error": "User not found"})
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
