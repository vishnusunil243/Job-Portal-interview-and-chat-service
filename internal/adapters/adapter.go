package adapters

import (
	"context"

	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ChatAdapter struct {
	DB *mongo.Database
}

func NewChatAdapter(db *mongo.Database) *ChatAdapter {
	return &ChatAdapter{
		DB: db,
	}
}
func (chat *ChatAdapter) InsertMessage(msg entities.InsertIntoRoomMessage) error {
	collection := chat.DB.Collection("chatCollection")
	filter := bson.D{{"room_id", msg.RoomID}}
	update := bson.D{{"$push", bson.D{{"messages", msg.Messages}}}}
	opts := options.Update().SetUpsert(true)
	_, err := collection.UpdateOne(context.TODO(), filter, update, opts)
	if err != nil {
		return err
	}
	return nil
}
func (chat *ChatAdapter) LoadMessages(roomId string) ([]entities.Message, error) {
	collection := chat.DB.Collection("chatCollection")
	filter := bson.D{{"room_id", roomId}}
	sort := bson.D{{"messages.time", 1}}
	var result struct {
		Messages []entities.Message `bson:"messages"`
	}
	if err := collection.FindOne(context.TODO(), filter, options.FindOne().SetSort(sort)).Decode(&result); err != nil {
		return []entities.Message{}, err
	}
	return result.Messages, nil
}
