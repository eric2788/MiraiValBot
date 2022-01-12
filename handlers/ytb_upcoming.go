package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/youtube"
)

func HandleUpcomingEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	dmDesc := fmt.Sprintf("[%s](%s) 在油管有预定直播", info.ChannelName, youtubeChannelLink(info.ChannelId))
	dm := createYoutubeMessageDiscord(dmDesc, info, "预定发布时间", "标题", "描述", "待机")
	go discord.SendNewsEmbed(dm)

	msg := createYoutubeMessage(fmt.Sprintf("%s 在油管有预定直播", info.ChannelName), info, "预定发布时间", "标题", "待机")
	return withRisky(msg)
}

func init() {
	youtube.RegisterDataHandler(youtube.UpComing, HandleUpcomingEvent)
}
