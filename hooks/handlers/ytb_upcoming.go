package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/hooks/sites/youtube"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/discord"
)

func HandleUpcomingEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	if info.Duplicate && file.DataStorage.Youtube.AntiDuplicate {
		return nil
	}

	dmDesc := fmt.Sprintf("[%s](%s) 在油管有预定直播", info.ChannelName, youtube.GetChannelLink(info.ChannelId))
	dm := youtube.CreateDiscordMessage(dmDesc, info, "预定发布时间", "标题", "待机")
	go discord.SendNewsEmbed(dm)

	return youtubeSendQQRisky(info, fmt.Sprintf("%s 在油管有预定直播", info.ChannelName), "标题", "预定发布时间", "待机")
}

func init() {
	youtube.MessageHandler.AddHandler(youtube.UpComing, HandleUpcomingEvent)
}
