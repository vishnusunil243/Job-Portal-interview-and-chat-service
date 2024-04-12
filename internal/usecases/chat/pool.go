package chat

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"
)

type Pool struct {
	ID        string
	JoinChan  chan *Client
	LeaveChan chan *Client
	Broadcast chan entities.Message
	Clients   map[string]*Client
}

func NewPool(id string) *Pool {
	return &Pool{
		ID:        id,
		JoinChan:  make(chan *Client),
		LeaveChan: make(chan *Client),
		Broadcast: make(chan entities.Message),
		Clients:   make(map[string]*Client),
	}
}

type Register struct {
	Message string    `json:"Message"`
	Time    time.Time `json:"Time"`
}

func (pool *Pool) Serve(insertChan chan<- entities.InsertIntoRoomMessage) {
	defer func() {
		close(pool.JoinChan)
		close(pool.LeaveChan)
		close(pool.Broadcast)
	}()
	for {
		select {
		case client := <-pool.JoinChan:
			for _, v := range pool.Clients {
				reg := Register{
					Time:    time.Now(),
					Message: fmt.Sprintf("%s is online", client.Name),
				}
				if err := v.Conn.WriteJSON(reg); err != nil {
					log.Println("error happened at sending ", err)
					continue
				}
			}
			pool.Clients[client.ClientID] = client
		case client := <-pool.LeaveChan:
			for _, v := range pool.Clients {
				unReg := Register{
					Time:    time.Now(),
					Message: fmt.Sprintf("%s is offline ", client.Name),
				}
				if err := v.Conn.WriteJSON(unReg); err != nil {
					log.Println("error at sending ", err)
					continue
				}
			}
			delete(pool.Clients, client.ClientID)
		case message := <-pool.Broadcast:
			for _, v := range pool.Clients {
				jsonData, err := json.Marshal(message)
				if err != nil {
					log.Println("error at sending")
					continue
				}
				if err := v.Conn.WriteMessage(message.Type, jsonData); err != nil {
					log.Println("error at sending")
					continue
				}
			}
			msg := entities.InsertIntoRoomMessage{
				RoomID:   pool.ID,
				Messages: message,
			}
			fmt.Println("message sent ", msg)
			insertChan <- msg
		}
	}
}
