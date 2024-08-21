package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

type Todo struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func main() {
	fmt.Println("Hello, World!")
	app := fiber.New()

	todos := []Todo{}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{"message": "Hello, World!"})
	})

	// Create a new todo
	app.Post("/api/todos", func(c *fiber.Ctx) error {
		todo := &Todo{}
		if err := c.BodyParser(todo); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
		}
		if todo.Title == "" {
			return c.Status(400).JSON(fiber.Map{"error": "Todo title is required"})
		}
		todo.ID = len(todos) + 1
		todos = append(todos, *todo)

		

		return c.Status(201).JSON(todo)
	})

	log.Fatal(app.Listen(":3000"))
}
