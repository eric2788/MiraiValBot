package twitter

import (
	"encoding/json"
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

type messageHandler struct {
}

func (m *messageHandler) PubSubPrefix() string {
	return "twitter:"
}

func (m *messageHandler) ToLiveData(message *redis.Message) (*TweetStreamData, error) {
	var twitterStream = &TweetStreamData{}
	err := json.Unmarshal([]byte(message.Payload), twitterStream)
	return twitterStream, err
}

func (m *messageHandler) GetCommand(data *TweetStreamData) string {
	return data.GetCommand()
}

func (m *messageHandler) GetOfflineListening() []string {
	listening := file.DataStorage.Listening.Twitter.ToArr()
	topics := make([]string, len(listening))
	for i, v := range listening {
		topics[i] = fmt.Sprintf("twitter:%s", v)
	}
	return topics
}

func (m *messageHandler) HandleError(bot *bot.Bot, error error) {
}

var MessageHandler = broadcaster.BuildHandle[TweetStreamData](logger, &messageHandler{})

func init() {
	broadcaster.RegisterHandler("twitter", MessageHandler)
}
