package usecases

import (
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/usecases/chat"
)

type ChatUsecaseInterface interface {
	CreatePoolifnotalreadyExists(string, chan<- entities.InsertIntoRoomMessage) (*chat.Pool, []entities.Message)
}
