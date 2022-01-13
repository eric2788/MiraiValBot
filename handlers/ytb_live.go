package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/datetime"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleLiveEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	dmDesc := fmt.Sprintf("[%s](%s) 正在油管直播", info.ChannelName, youtube.GetChannelLink(info.ChannelId))
	dm := youtube.CreateDiscordMessage(dmDesc, info)
	go discord.SendNewsEmbed(dm)

	return youtubeSendQQRisky(info, fmt.Sprintf("%s 正在油管直播", info.ChannelName))
}

func youtubeSendQQRisky(info *youtube.LiveInfo, desc string, blocks ...string) (err error) {

	titles := []string{"标题", "开始时间", "直播间"}

	for i, block := range blocks {
		titles[i] = block
	}

	go qq.SendRiskyMessage(5, 10, func(currentTry int) error {
		fields := make(map[string]string)
		image := true
		switch currentTry {
		case 0: // 风控0次，所有标题
			fields[titles[0]] = info.Info.Title
			fallthrough
		case 1: // 风控一次，没有标题
			fields[titles[2]] = datetime.FormatMillis(info.Info.PublishTime)
			logger.Warnf("油管广播被风控 %d 次，舍弃 %s 重发", currentTry, titles[0])
			fallthrough
		case 2: // 风控两次， 没有开始时间
			fields[titles[1]] = fmt.Sprintf("https://youtu.be/%s", info.Info.Id)
			logger.Warnf("油管广播被风控 %d 次，舍弃 %s 重发", currentTry, titles[1])
			fallthrough
		case 4: // 风控三次，没有图片
			image = false
			logger.Warnf("油管广播被风控 %d 次，舍弃 %s 重发", currentTry, "图片")
		}
		msg := youtube.CreateQQMessage(desc, info, image, titles[2], fields)
		return qq.SendGroupMessage(msg)
	})

	return
}

func init() {
	youtube.RegisterDataHandler(youtube.Live, HandleLiveEvent)
}
