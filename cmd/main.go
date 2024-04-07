package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/db"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/initializer"
)

// func main() {
// 	server.AllRooms.Init()
// 	http.Handle("/join", http.HandlerFunc(server.JoinRoomRequestHandler))
// 	http.Handle("/create", http.HandlerFunc(server.CreateRoomRequestHandler))
// 	go server.Broadcaster()
// 	log.Println("listening on port 8000")
// 	err := http.ListenAndServe(":8000", nil)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// }
func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal("error loading env")
	}
	addr := os.Getenv("MONGO_KEY")
	db, err := db.InitMongoDB(addr)
	if err != nil {
		log.Fatal("error connecting to database")
	}
	handler := initializer.Initialize(db)
	handler.Start()
}
