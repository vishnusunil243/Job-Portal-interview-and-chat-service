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
func (c *ChatUsecase) CreatePoolifnotalreadyExists(poolid string, insertChan chan<- entities.InsertIntoRoomMessage) (*chat.Pool, []entities.Message) {
	res, err := c.adapters.LoadMessages(poolid)
	ids := strings.Split(poolid, " ")
	if err != nil {
		log.Println("error while loading messages", err)
		res, err = c.adapters.LoadMessages(ids[1] + " " + ids[0])
		if err != nil {
			log.Println("error retrieving message ", err)
		}
	}
	if c.Chat.Pool[poolid] == nil {
		fmt.Println("poolId for second try ", ids)
		if c.Chat.Pool[ids[1]+" "+ids[0]] == nil {
			fmt.Println("no message found for id ", ids[1], ids[0])
			pool := chat.NewPool(ids[1] + " " + ids[0])
			go pool.Serve(insertChan)
			c.Chat.Pool[ids[1]+" "+ids[0]] = pool
			return pool, res
		}
		return c.Chat.Pool[ids[1]+" "+ids[0]], res

	}
	fmt.Println("pool id is ", poolid)
	return c.Chat.Pool[poolid], res
}
