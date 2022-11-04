package youtube

import (
	"encoding/json"
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

var logger = utils.GetModuleLogger("sites.youtube")

type messageHandler struct {
}

func (m *messageHandler) PubSubPrefix() string {
	return "ylive:"
}

func (m *messageHandler) ToLiveData(message *redis.Message) (*LiveInfo, error) {
	var liveInfo = &LiveInfo{}
	err := json.Unmarshal([]byte(message.Payload), liveInfo)
	return liveInfo, err
}

func (m *messageHandler) GetCommand(data *LiveInfo) string {
	return data.Status
}

func (m *messageHandler) GetOfflineListening() []string {
	listening := file.DataStorage.Listening.Youtube.ToArr()
	topics := make([]string, len(listening))
	for i, v := range listening {
		topics[i] = fmt.Sprintf("ylive:%s", v)
	}
	return topics
}

func (m *messageHandler) HandleError(bot *bot.Bot, error error) {
}

var MessageHandler = broadcaster.BuildHandle[LiveInfo](logger, &messageHandler{})

func init() {
	broadcaster.RegisterHandler("youtube", MessageHandler)
}
