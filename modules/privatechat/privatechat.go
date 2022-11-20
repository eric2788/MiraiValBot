package privatechat

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/qq"
)

const Tag = "valbot.privatechat"

var (
	logger = utils.GetModuleLogger(Tag)
)

type privateChatResponse struct {
}

func (p *privateChatResponse) HookEvent(bot *bot.Bot) {
	bot.PrivateMessageEvent.Subscribe(func(client *client.QQClient, event *message.PrivateMessage) {

		// 非群友
		if info := qq.FindGroupMember(event.Sender.Uin); info == nil {
			// 無視
			return
		}

		// 暂时没啥好发的，就echo吧

		msg := message.NewSendingMessage()

		for _, e := range event.Elements {
			msg.Append(e)
		}

		_ = qq.SendPrivateMessage(event.Sender.Uin, msg)
	})
}

func init() {
	eventhook.RegisterAsModule(&privateChatResponse{}, "私聊回應", Tag, logger)
}
