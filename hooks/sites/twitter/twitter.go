package twitter

import (
	"encoding/json"
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

type messageHandler struct {
}

func (m *messageHandler) PubSubPrefix() string {
	return "twitter:"
}

func (m *messageHandler) ToLiveData(message *redis.Message) (*TweetContent, error) {
	var twitterStream = &TweetContent{}
	err := json.Unmarshal([]byte(message.Payload), twitterStream)
	return twitterStream, err
}

func (m *messageHandler) GetCommand(data *TweetContent) string {
	if data.Tweet == nil {
		return ""
	}
	return data.Tweet.GetCommand()
}

func (m *messageHandler) GetOfflineListening() []string {
	listening := file.DataStorage.Listening.Twitter.ToSlice()
	topics := make([]string, len(listening))
	for i, v := range listening {
		topics[i] = fmt.Sprintf("twitter:%s", v)
	}
	return topics
}

func (m *messageHandler) HandleError(bot *bot.Bot, error error) {
}

var MessageHandler = broadcaster.BuildHandle[TweetContent](logger, &messageHandler{})

func init() {
	broadcaster.RegisterHandler("twitter", MessageHandler)
}
