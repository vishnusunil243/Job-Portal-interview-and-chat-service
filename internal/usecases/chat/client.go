package chat

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"
)

type Client struct {
	Conn     *websocket.Conn
	ClientID string
	Name     string
	Pool     *Pool
}

func NewClient(conn *websocket.Conn, clientID, name string, pool *Pool) *Client {
	return &Client{
		Conn:     conn,
		ClientID: clientID,
		Name:     name,
		Pool:     pool,
	}
}
func (client *Client) Serve(msgs []entities.Message) {
	client.Pool.JoinChan <- client
	defer func() {
		client.Pool.LeaveChan <- client
		client.Conn.Close()
	}()
	for _, v := range msgs {
		jsonData, err := json.Marshal(v)
		if err != nil {
			log.Println("error at sending ", err)
			continue
		}
		if err := client.Conn.WriteMessage(v.Type, jsonData); err != nil {
			log.Println("error at sending ", err)
			continue
		}
	}
	for {
		msgtype, p, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("error happened, closing connection")
			break
		}
		message := entities.Message{Type: msgtype, Message: string(p), Time: time.Now(), Name: client.Name}
		client.Pool.Broadcast <- message
		log.Printf("message recieved from %s", client.ClientID)
	}
}
