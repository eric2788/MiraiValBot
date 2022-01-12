package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/datetime"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleLiveEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 正在油管直播", info.ChannelName))

	if info.Info != nil {
		msg.Append(qq.NewTextfLn("标题: %s", info.Info.Title))
		msg.Append(qq.NewTextfLn("开始时间: %s", datetime.FormatSeconds(info.Info.PublishTime)))
		msg.Append(qq.NewTextfLn("直播间: %s", getYTLink(info)))

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
		msg.Append(qq.NewTextf("直播间: %s", getYTLink(info)))
	}

	return withRisky(msg)
}

func getYTLink(info *youtube.LiveInfo) string {
	if info.Info != nil {
		return fmt.Sprintf("https://youtu.be/%s", info.Info.Id)
	} else {
		return fmt.Sprintf("https://youtube.com/channel/%s/live", info.ChannelId)
	}
}

func init() {
	youtube.RegisterDataHandler(youtube.Live, HandleLiveEvent)
}
