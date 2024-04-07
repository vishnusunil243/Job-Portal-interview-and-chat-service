package adapters

import "github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"

type ChatAdapterInterface interface {
	InsertMessage(msg entities.InsertIntoRoomMessage) error
	LoadMessages(roomId string) ([]entities.Message, error)
}
