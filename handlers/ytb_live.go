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

	go qq.SendRiskyMessage(5, 10, func(try int) error {
		fields := make(map[string]string)
		image := true
		if info.Info != nil { // 防止 NPE

			if try < 1 { // 風控一次，沒有標題
				fields[titles[0]] = info.Info.Title
			}

			if try < 2 { // 風控兩次，沒有開始時間
				t, err := datetime.ParseISOStr(info.Info.PublishTime)
				if err != nil {
					fields[titles[1]] = datetime.FormatMillis(t.UnixMilli())
				} else {
					logger.Warnf("解析時間文字 %s 時出現錯誤: %v", info.Info.PublishTime, err)
					fields[titles[1]] = info.Info.PublishTime // 使用原本的 string
				}
			}

			if try < 3 { // 風控三次，沒有圖片
				image = false
			}

			fields[titles[2]] = fmt.Sprintf("https://youtu.be/%s", info.Info.Id)

		}
		msg := youtube.CreateQQMessage(desc, info, image, titles[2], fields)
		return qq.SendGroupMessage(msg)
	})

	return
}

func init() {
	youtube.RegisterDataHandler(youtube.Live, HandleLiveEvent)
}
