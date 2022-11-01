package chat_reply

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/eventhook"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/common-utils/array"
)

var logger = utils.GetModuleLogger("valbot.chat_reply")

type (
	AtResponse struct {
		strategies []ResponseStrategy
	}

	ResponseStrategy interface {
		Response(msg *message.GroupMessage) (*message.SendingMessage, error)
	}
)

func (a *AtResponse) HookEvent(bot *bot.Bot) {
	bot.OnGroupMessage(func(cl *client.QQClient, msg *message.GroupMessage) {
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
	eventhook.HookLifeCycle(&AtResponse{
		strategies: []ResponseStrategy{
			&aiChatResponse{},
			&randomResponse{},
		},
	})
}
