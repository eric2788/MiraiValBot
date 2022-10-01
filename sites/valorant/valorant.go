package valorant

import (
	"encoding/json"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/broadcaster"
	"github.com/go-redis/redis/v8"
)

var logger = utils.GetModuleLogger("sites.valorant")

const MatchesUpdated = "matches_updated"

type messageHandler struct {
}

func (m *messageHandler) PubSubPrefix() string {
	return "valorant:"
}

func (m *messageHandler) GetOfflineListening() []string {
	listening := file.DataStorage.Listening.Valorant.ToArr()
	topics := make([]string, len(listening))
	for i, v := range listening {
		topics[i] = fmt.Sprintf("valorant:%s", v)
	}
	return topics
}

func (m *messageHandler) ToLiveData(message *redis.Message) (interface{}, error) {
	var matchData = &MatchMetaDataSub{}
	err := json.Unmarshal([]byte(message.Payload), matchData)
	return matchData, err
}

func (m *messageHandler) HandleError(bot *bot.Bot, error error) {
}

func (m *messageHandler) GetCommand(data interface{}) string {
	return MatchesUpdated // only one
}

var MessageHandler = broadcaster.BuildHandle(logger, &messageHandler{})

func init() {
	broadcaster.RegisterHandler("valorant", MessageHandler)
}
