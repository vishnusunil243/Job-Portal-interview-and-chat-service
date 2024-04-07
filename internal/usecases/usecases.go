package usecases

import (
	"fmt"
	"log"
	"strings"

	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/entities"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/adapters"
	"github.com/vishnusunil243/Job-Portal-interview-and-chat-service/internal/usecases/chat"
)

type ChatUsecase struct {
	adapters adapters.ChatAdapterInterface
	Chat     *chat.ChatPool
}

func NewChatUsecase(adapters adapters.ChatAdapterInterface) ChatUsecase {
	return ChatUsecase{
		adapters: adapters,
		Chat:     chat.NewChatPool(),
	}
}
func (c *ChatUsecase) SpinupPoolifnotalreadyExists(poolid string, insertChan chan<- entities.InsertIntoRoomMessage) (*chat.Pool, []entities.Message) {
	res, err := c.adapters.LoadMessages(poolid)
	if err != nil {
		log.Println("error while loading messages", err)
	}
	if c.Chat.Pool[poolid] == nil {
		ids := strings.Split(poolid, " ")
		if c.Chat.Pool[ids[1]+" "+ids[0]] == nil {
			pool := chat.NewPool(ids[1] + " " + ids[0])
			go pool.Serve(insertChan)
			c.Chat.Pool[ids[1]+" "+ids[0]] = pool
			fmt.Println("eeeeee")
			return pool, res
		}
		fmt.Println(("iiiiii"))
		return c.Chat.Pool[ids[1]+" "+ids[0]], res

	}
	fmt.Println("111")
	return c.Chat.Pool[poolid], res
}
