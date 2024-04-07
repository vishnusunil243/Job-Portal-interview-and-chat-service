package server

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap
var mutex = &sync.Mutex{} // Mutex for concurrent access to AllRooms.Map
var broadcast = make(chan broadcastMsg)

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()
	type resp struct {
		RoomID string `json:"room_id"`
	}
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Broadcaster() {
	for {
		msg := <-broadcast

		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client && client.Conn != nil {
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					log.Println("Error sending message to client:", err)
					go func() {
						mutex.Lock()
						defer mutex.Unlock()
						AllRooms.DeleteRoom(msg.RoomID)
					}()
				}
			}
		}
		log.Println(msg)
	}
}

func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomId"]
	if !ok {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	client := wsClient{Conn: ws, RoomID: roomID[0]}

	go func() {
		defer ws.Close()
		AllRooms.InsertIntoRoom(roomID[0], false, client.Conn)
		for {
			var msg broadcastMsg
			err := ws.ReadJSON(&msg.Message)
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
					log.Println("Client disconnected:", client.Conn.RemoteAddr())
				} else {
					log.Println("Read Error:", err)
				}
				break
			}
			msg.Client = client.Conn
			msg.RoomID = client.RoomID
			broadcast <- msg
		}

		mutex.Lock()
		defer mutex.Unlock()
		AllRooms.DeleteRoom(roomID[0])
	}()
}

type wsClient struct {
	*websocket.Conn
	RoomID string
}
