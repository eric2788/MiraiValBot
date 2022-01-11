package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleIdle(bot *bot.Bot, info *youtube.LiveInfo) error {
	msg := message.NewSendingMessage().Append(qq.NewTextf("%s 的油管直播已结束。", info.ChannelName))
	bot.SendGroupMessage(qq.ValGroupInfo.Uin, msg)
	return nil
}

func init() {
	youtube.RegisterDataHandler(youtube.Idle, HandleIdle)
}
