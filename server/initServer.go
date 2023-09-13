package server

import (
	"cpu-mon/database"
	"cpu-mon/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func InitializeServer() error {
	InitEnv()

	mongoErr := database.InitMongoDB()
	if mongoErr != nil {
		return mongoErr
	}

	defer database.DisconnectMongoDB()

	server := fiber.New()

	server.Use(recover.New())
	server.Use(logger.New(logger.Config{
		Format: "[${ip}]:${port} ${status} - ${method} ${path} ${latency}\n",
	}))

	routes.InitRoutes(server)
	log.Fatal(server.Listen(":" + os.Getenv("PORT")))
	return nil
}

func InitEnv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
		return err
	}
	return nil
}
