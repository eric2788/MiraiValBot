package chat_reply

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/game"
	"github.com/eric2788/common-utils/array"
)

const Tag = "valbot.chat_reply"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &atResponse{
		strategies: []ResponseStrategy{
			AIChat,
			&RandomResponse{},
		},
	}
)

type (
	atResponse struct {
		strategies []ResponseStrategy
	}

	ResponseStrategy interface {
		Response(msg *message.GroupMessage) (*message.SendingMessage, error)
	}
)

func (a *atResponse) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(cl *client.QQClient, msg *message.GroupMessage) {

		if game.IsInGame() {
			return
		}

		content := qq.ParseMsgContent(msg.Elements)

		if array.IndexOfInt64(content.At, cl.Uin) != -1 && len(content.Texts) > 0 {

			for _, strategy := range a.strategies {
				send, err := strategy.Response(msg)
				if err == nil {
					_ = qq.SendGroupMessageByGroup(msg.GroupCode, send)
					break
				}
			}

		}
	})
}

func init() {
	eventhook.RegisterAsModule(instance, "自動回復", Tag, logger)
}
