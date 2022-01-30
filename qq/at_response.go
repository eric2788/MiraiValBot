package qq

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/eventhook"
	"github.com/eric2788/common-utils/array"
)

type AtResponse struct {
}

func (a *AtResponse) HookEvent(bot *bot.Bot) {
	bot.OnGroupMessage(func(cl *client.QQClient, msg *message.GroupMessage) {
		content := ParseMsgContent(msg.Elements)

		if array.IndexOfInt64(content.At, cl.Uin) != -1 && len(content.Texts) > 0 {
			random, err := GetRandomGroupMessage(msg.GroupCode)
			if err != nil {
				logger.Warnf("嘗試發送隨機群消息時出現錯誤: %v, 已略過發送。", err)
				return
			}
			send := CreateReply(msg)

			for _, ele := range random.Elements {

				switch ele.(type) {
				case *message.ReplyElement:
					continue
				case *message.ForwardElement:
					continue
				default:
					break
				}

				send.Append(ele)
			}

			_ = SendGroupMessageByGroup(msg.GroupCode, send)
		}
	})
}

func init() {
	eventhook.HookLifeCycle(&AtResponse{})
}
