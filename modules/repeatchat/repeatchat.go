package repeatchat

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/qq"
)

const Tag = "valbot.repeatchat"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &repeatChat{}
)

// 复读操作 -> 打断/復讀
type repeatChat struct {
	repeatRaw string
}

// 参考了 FloatTech/ZeroBot-Plugin 的复读判断
func (r *repeatChat) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		
		// 無視群機器人的消息
		if event.Sender.Uin == client.Uin {
			return
		}

		content := event.ToString()
		if r.repeatRaw == "" {
			r.repeatRaw = "0:" + content
		} else {
			lastContent := r.repeatRaw[2:]
			// 复读被打断，重新计算
			if lastContent != content {
				logger.Debugf("群消息与上一则不一样: %s, 已重新计算。", lastContent)
				r.repeatRaw = "0:" + content
			} else {
				times := int(r.repeatRaw[0] - '0')
				c := times + 1
				logger.Debugf("群消息与上一则相同: %s (%d + 1)", lastContent, times)
				r.repeatRaw = fmt.Sprintf("%d:%s", c, lastContent)

				// 3 的倍数时候开始进行操作
				if c%3 == 0 && c != 0 {

					rand.Seed(time.Now().UnixNano())

					msg := message.NewSendingMessage()

					// 60% 复读, 40% 打断复读
					if rand.Intn(100)+1 > 60 {
						logger.Debugf("操作为: 复读")
						msg.Append(message.NewText(lastContent))
					} else {
						logger.Debugf("操作为: 打断复读")
						lcrune := []rune(lastContent)
						rand.Shuffle(len(lcrune), func(i, j int) {
							lcrune[i], lcrune[j] = lcrune[j], lcrune[i]
						})

						msg.Append(message.NewText(string(lcrune)))
					}

					_ = qq.SendGroupMessageByGroup(event.GroupCode, msg)
				}

			}
		}
	})
}

func init() {
	eventhook.RegisterAsModule(instance, "復讀操作", Tag, logger)
}
