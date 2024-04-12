package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/db"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/initializer"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("error loading env")
	}
	addr := os.Getenv("MONGO_KEY")
	db, err := db.InitMongoDB(addr)
	if err != nil {
		log.Fatal("error connecting to database, ", err)
	}
	handler := initializer.Initialize(db)
	handler.Start()
}
