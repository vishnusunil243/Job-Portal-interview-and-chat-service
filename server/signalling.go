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
	type request struct {
		RoomId string `json:"roomId"`
	}
	var req request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	roomID := AllRooms.CreateRoom(req.RoomId)
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
				// Check if connection is still valid before sending
				err := client.Conn.WriteJSON(msg.Message)
				if err != nil {
					log.Println("Error sending message to client:", err)
					// Consider removing the disconnected client from AllRooms
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

	client := wsClient{Conn: ws, RoomID: roomID[0]} // Create a client struct with connection and room ID

	// Add client to AllRooms with a goroutine to handle disconnections gracefully
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
				break // Exit the loop on error or normal closure
			}
			msg.Client = client.Conn
			msg.RoomID = client.RoomID
			broadcast <- msg
		}
		// Remove client from AllRooms upon disconnection
		mutex.Lock()
		defer mutex.Unlock()
		AllRooms.DeleteRoom(roomID[0])
	}()
}

type wsClient struct {
	*websocket.Conn
	RoomID string
}
