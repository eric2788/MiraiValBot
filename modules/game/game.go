package game

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/common-utils/array"
)

const Tag = "game"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &game{}
)

type game struct {
	Gaming bool
}

func (g *game) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {
		if currentGame == nil {
			return
		}

		content := qq.ParseMsgContent(event.Elements)

		if array.Contains(content.At, client.Uin) && len(content.Texts) > 0 {

			result := currentGame.Handle(event)
			if result.EndGame {
				msg := message.NewSendingMessage()
				msg.Append(qq.NewTextfLn("游戏已结束。"))
				if result.Winner != "" {
					msg.Append(qq.NewTextfLn("胜者: %s", result.Winner))
				}
				if result.Score > 0 {
					msg.Append(qq.NewTextfLn("分数: %d", result.Score))
				}
				_ = qq.SendGroupMessage(msg)
				_ = StopGame()
			}

		}

	})
}

func (g *game) StopEvent(bot *bot.Bot) {
	logger.Info(StopGame())
}

func init() {
	eventhook.RegisterAsModule(instance, "文字游戏", Tag, logger)
}
