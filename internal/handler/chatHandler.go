package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/helper"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/usecases"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/usecases/chat"
	videocall "github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/usecases/server"
	"github.com/vishnusunil243/Job-Portal-proto-files/pb"
)

type ChatHandlers struct {
	InsertChannel chan<- entities.InsertIntoRoomMessage
	Usecase       usecases.ChatUsecaseInterface
	UserConn      pb.UserServiceClient
	CompanyConn   pb.CompanyServiceClient
	Upgrader      websocket.Upgrader
}

func NewChatHandlers(usecase usecases.ChatUsecaseInterface, insertChan chan<- entities.InsertIntoRoomMessage, userAddr, companyAddr string) *ChatHandlers {
	userRes, _ := helper.DialGrpc(userAddr)
	companyRes, _ := helper.DialGrpc(companyAddr)
	return &ChatHandlers{
		InsertChannel: insertChan,
		Usecase:       usecase,
		UserConn:      pb.NewUserServiceClient(userRes),
		CompanyConn:   pb.NewCompanyServiceClient(companyRes),
		Upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}
func (c *ChatHandlers) Handler(w http.ResponseWriter, r *http.Request) {
	r.Header.Del("Sec-WebSocket-Extensions")
	userId := r.Header.Get("userId")
	companyId := r.Header.Get("companyId")
	recieverId := r.Header.Get("recieverId")
	var poolId string
	if userId != "" && recieverId != "" {
		poolId = userId + " " + recieverId
	} else if companyId != "" && recieverId != "" {
		poolId = companyId + " " + recieverId
	} else {
		http.Error(w, "please provide valid headers", http.StatusBadRequest)
		return
	}
	userData, err := c.UserConn.GetUser(context.Background(), &pb.GetUserById{Id: userId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	name := userData.Name
	clientId := userId
	if companyId != "" {
		companyData, err := c.CompanyConn.GetCompany(context.Background(), &pb.GetJobByCompanyId{
			Id: companyId,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		name = companyData.Name
		clientId = companyId
	}
	conn, err := c.Upgrader.Upgrade(w, r, r.Header)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	pool, msgs := c.Usecase.CreatePoolifnotalreadyExists(poolId, c.InsertChannel)
	client := chat.NewClient(conn, clientId, name, pool)
	client.Serve(msgs)
}
func (chat *ChatHandlers) Start() {
	mux := http.NewServeMux()
	videocall.AllRooms.Init()
	go videocall.Broadcaster()
	mux.HandleFunc("/ws", chat.Handler)
	mux.HandleFunc("/create", videocall.CreateRoomRequestHandler)
	mux.HandleFunc("/join", videocall.JoinRoomRequestHandler)
	log.Println("listening on port 8000")
	if err := http.ListenAndServe(":8000", mux); err != nil {
		fmt.Println(err.Error())
	}
}
