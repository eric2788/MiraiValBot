package games

import (
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
)

var logger = utils.GetModuleLogger("valbot.games")

func risky(err error) {
	if err == nil {
		return
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("发送游戏信息时出现错误: %v", err))
	_ = qq.SendGroupMessage(msg)
}
