package games

import (
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
)

var logger = utils.GetModuleLogger("valbot.games")

func risky(err error) error {
	if err == nil {
		return nil
	}
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("发送游戏信息时出现错误: %v", err))
	return qq.SendGroupMessage(msg)
}

func sendGameMsg(msg *message.SendingMessage) (err error) {
	err = qq.SendGroupMessage(msg)
	if err != nil {
		logger.Warnf("发送游戏信息时出现错误: %v, 将改用文字图片", err)
		err = risky(qq.SendGroupImageText(msg))
	}
	return
}
