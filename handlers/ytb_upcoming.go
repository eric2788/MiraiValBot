package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/datetime"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleUpcomingEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 在油管有预定直播", info.ChannelName))

	if info.Info != nil {

		msg.Append(qq.NewTextfLn("标题: %s", info.Info.Title))
		msg.Append(qq.NewTextfLn("预定发布时间: %s", datetime.FormatSeconds(info.Info.PublishTime)))
		msg.Append(qq.NewTextfLn("待机: %s", getYTLink(info)))

		if info.Info.Cover != nil {
			cover := *info.Info.Cover
			img, err := qq.NewImageByUrl(cover)
			if err != nil {
				logger.Warnf("获取图片 %s 时出现错误: %v", cover, err)
			} else {
				msg.Append(img)
			}
		}

	} else {
		msg.Append(qq.NewTextf("待机: %s", getYTLink(info)))
	}

	return withRisky(msg)
}

func init() {
	youtube.RegisterDataHandler(youtube.UpComing, HandleUpcomingEvent)
}
