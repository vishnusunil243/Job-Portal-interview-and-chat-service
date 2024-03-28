package main

import (
	"log"
	"net/http"

	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/server"
)

func main() {
	server.AllRooms.Init()
	http.Handle("/join", http.HandlerFunc(server.JoinRoomRequestHandler))
	http.Handle("/create", http.HandlerFunc(server.CreateRoomRequestHandler))
	go server.Broadcaster()
	log.Println("listening on port 8000")
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		log.Fatal(err)
	}

}
