package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id, omitempty" bson:"_id, omitempty"`
	Title     string             `json:"title" bson:"title"`
	Completed bool               `json:"completed" bson:"completed"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("Hello, World!")

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")
	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB")

	collection = client.Database("react-go-tutorial").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Put("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")

	if port == "" {
		port = "3001"
	}

	log.Fatal(app.Listen("0.0.0.0" + port))

}

// Create todo
func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": err.Error()})
		}

		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

func createTodo(c *fiber.Ctx) error {
	// Implement createTodo function here
	var todo Todo
	if err := c.BodyParser(&todo); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if todo.Title == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Title is required"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to create todo"})
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo)
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to update todo"})
	}

	return c.SendStatus(204)
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	}

	filter := bson.M{"_id": objectId}
	_, err = collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to delete todo"})
	}

	return c.SendStatus(204)
}
