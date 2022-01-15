package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"strings"
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
		if i == len(titles) {
			break
		}
		titles[i] = block
	}

	go qq.SendRiskyMessage(5, 10, func(try int) error {
		alt := make([]string, 0)
		noTitle := false

		// 风控时尝试加随机文字看看会不会减低？

		if try >= 1 {
			alt = append(alt, fmt.Sprintf("[此油管广播已被风控 %d 次]", try))
		}

		if try >= 2 {
			alt = append(alt, fmt.Sprintf("你好谢谢小笼包再见"))
		}

		if try >= 3 {
			alt = append(alt, fmt.Sprintf("看看啊！这谁直播了额!"))
		}

		// 风控第四次没有标题
		if try >= 4 {
			alt = append(alt, fmt.Sprintf("什么油管直播居然被风控了四次这么爽？？你标题没了"))
			noTitle = true
		}

		if try > 0 {
			logger.Warnf("为被风控的推文新增如下的内容: %s", strings.Join(alt, "\n"))
		}

		msg := youtube.CreateQQMessage(desc, info, noTitle, alt, titles...)
		return qq.SendGroupMessage(msg)
	})

	return
}

func init() {
	youtube.RegisterDataHandler(youtube.Live, HandleLiveEvent)
}
