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
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title"`
	Completed bool               `json:"completed"`
}

var collection *mongo.Collection

func main() {
	if os.Getenv("ENV") != "production" {
		// Load the .env file if not in production
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file:", err)
		}
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

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:5173",
	// 	AllowHeaders: "Origin,Content-Type,Accept",
	// }))

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Patch("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8081"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen("0.0.0.0:"  +  port))

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
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
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

	return c.Status(200).JSON(fiber.Map{"success": true})
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

	return c.Status(200).JSON(fiber.Map{"success": true})
}
