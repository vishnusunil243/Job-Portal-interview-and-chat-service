package initializer

import (
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/adapters"
	handlers "github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/handler"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/usecases"
	"go.mongodb.org/mongo-driver/mongo"
)

func Initialize(mongo *mongo.Database) *handlers.ChatHandlers {
	adapter := adapters.NewChatAdapter(mongo)
	usecase := usecases.NewChatUsecase(adapter)
	insertRoom := usecase.InsertIntoDB()
	handler := handlers.NewChatHandlers(&usecase, insertRoom, "user-service:8081", "company-service:8082")
	return handler

}
