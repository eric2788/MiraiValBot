package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleLive(bot *bot.Bot, data *bilibili.LiveData) error {

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 正在B站直播", data.LiveInfo.Name))
	msg.Append(qq.NewTextfLn("标题: %s", data.LiveInfo.Title))
	msg.Append(qq.NewTextfLn("直播间: https://live.bilibili.com/%d", data.LiveInfo.RoomId))
	if data.LiveInfo.Cover != nil {
		cover := *data.LiveInfo.Cover
		imgElement, err := qq.NewImageByUrl(cover)
		if err != nil {
			logger.Warnf("获取图片 %s 时出现错误: %v", cover, err)
		} else {
			msg.Append(imgElement)
		}
	}
	bot.SendGroupMessage(qq.ValGroupInfo.Uin, msg)
	return nil
}

func init() {
	bilibili.RegisterDataHandler(bilibili.Live, HandleLive)
}
