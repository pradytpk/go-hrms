// Package main is the entry point for the Fiber HRMS (Human Resource Management System) application.
package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoInstance holds the MongoDB client and database instances.
type MongoInstance struct {
	Client *mongo.Client
	Db     *mongo.Database
}

// mg represents the global MongoInstance.
var mg MongoInstance

// dbName holds the name of the MongoDB database.
const dbName = "fiber-hrms"

// mongoURI holds the connection URI for MongoDB.
const mongoURI = "mongodb://localhost:27017/" + dbName

// Employee represents the structure of an employee.
type Employee struct {
	ID     string  `json:"id,omitempty" bson:"_id,omitempty"`
	Name   string  `json:"name"`
	Salary float64 `json:"salary"`
	Age    float64 `json:"age"`
}

// Connect establishes a connection to MongoDB.
func Connect() error {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		return err
	}
	db := client.Database(dbName)
	mg = MongoInstance{
		Client: client,
		Db:     db,
	}
	return nil
}

func main() {
	// Attempt to connect to MongoDB.
	if err := Connect(); err != nil {
		log.Error(err.Error())
	}

	// Initialize Fiber app.
	app := fiber.New()

	// GET endpoint to fetch all employees.
	app.Get("/employee", func(c *fiber.Ctx) error {
		query := bson.D{{}}
		cursor, err := mg.Db.Collection("employees").Find(c.Context(), query)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		var employees []Employee = make([]Employee, 0)
		if err := cursor.All(c.Context(), &employees); err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(employees)
	})

	// POST endpoint to create a new employee.
	app.Post("/employee", func(c *fiber.Ctx) error {
		collection := mg.Db.Collection("employees")
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		employee.ID = ""
		insertResult, err := collection.InsertOne(c.Context(), employee)
		if err != nil {
			return c.Status(500).SendString(err.Error())
		}
		filter := bson.D{{Key: "_id", Value: insertResult.InsertedID}}
		createdRecord := collection.FindOne(c.Context(), filter)

		createdEmployee := &Employee{}
		createdRecord.Decode(createdEmployee)
		return c.Status(201).JSON(createdEmployee)
	})

	// PUT endpoint to update an existing employee.
	app.Put("/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		employeedId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.SendStatus(400)
		}
		employee := new(Employee)
		if err := c.BodyParser(employee); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		query := bson.D{{Key: "_id", Value: employeedId}}
		update := bson.D{
			{Key: "$set",
				Value: bson.D{
					{Key: "name", Value: employee.Name},
					{Key: "age", Value: employee.Age},
					{Key: "salary", Value: employee.Salary},
				},
			},
		}
		err = mg.Db.Collection("employees").FindOneAndUpdate(c.Context(), query, update).Err()
		if err != nil {
			if err == mongo.ErrNoDocuments {
				return c.Status(404).SendString(err.Error())
			}
			return c.SendStatus(500)
		}
		employee.ID = idParam
		return c.Status(200).JSON(employee)
	})

	// DELETE endpoint to delete an employee.
	app.Delete("/employee/:id", func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		employeedId, err := primitive.ObjectIDFromHex(idParam)
		if err != nil {
			return c.SendStatus(400)
		}
		query := bson.D{{Key: "_id", Value: employeedId}}

		result, err := mg.Db.Collection("employees").DeleteOne(c.Context(), &query)
		if err != nil {
			return c.SendStatus(500)
		}
		if result.DeletedCount < 1 {
			return c.SendStatus(404)
		}
		return c.SendStatus(200)
	})

	// Start Fiber app.
	log.Fatal(app.Listen(":3000"))
}
