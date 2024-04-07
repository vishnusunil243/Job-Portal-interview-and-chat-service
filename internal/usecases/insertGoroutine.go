package usecases

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"
)

func (chat *ChatUsecase) InsertIntoDB() chan<- entities.InsertIntoRoomMessage {
	insertChan := make(chan entities.InsertIntoRoomMessage, 100)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	var run = true
	go func() {
		defer func() {
			for v := range insertChan {
				if err := chat.adapters.InsertMessage(v); err != nil {
					log.Println("error while inserting message", err)
				}
			}
			close(insertChan)
			close(sigChan)
		}()
		for run {
			select {
			case <-sigChan:
				run = false
			case msg := <-insertChan:
				fmt.Println("message recieved from channel")
				if err := chat.adapters.InsertMessage(msg); err != nil {
					log.Println("error happened at insertMessage adapter ", err)
				}
			}
		}
	}()
	return insertChan

}
