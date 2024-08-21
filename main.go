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

	app.Get("/api/todos", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
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

	// Update a todo
	app.Patch("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		for i, todo := range todos {
			if todo.ID == id {
				todos[i].Completed = true
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	// Delete a todo
	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
		id, err := c.ParamsInt("id")
		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
		}

		for i, todo := range todos {
			if todo.ID == id {
				// Remove the todo from the slice
				todos = append(todos[:i], todos[i+1:]...)
				return c.Status(200).JSON(fiber.Map{"sucess": true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	})

	log.Fatal(app.Listen(":3000"))
}
