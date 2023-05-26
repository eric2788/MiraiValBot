package valorant

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/file"
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
	listening := file.DataStorage.Listening.Valorant.ToSlice()
	topics := make([]string, len(listening))
	for i, line := range listening {
		parts := strings.Split(line, "//")
		if len(parts) != 2 {
			logger.Warnf("Invalid line in listening: %s", line)
			continue
		}
		topics[i] = fmt.Sprintf("valorant:%s", parts[0])
	}
	return topics
}

func (m *messageHandler) ToLiveData(message *redis.Message) (*MatchMetaDataSub, error) {
	var matchData = &MatchMetaDataSub{}
	err := json.Unmarshal([]byte(message.Payload), matchData)
	return matchData, err
}

func (m *messageHandler) HandleError(bot *bot.Bot, error error) {
}

func (m *messageHandler) GetCommand(data *MatchMetaDataSub) string {
	return MatchesUpdated // only one
}

var MessageHandler = broadcaster.BuildHandle[MatchMetaDataSub](logger, &messageHandler{})

func init() {
	broadcaster.RegisterHandler("valorant", MessageHandler)
}
