package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/youtube"
)

func HandleLiveEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	dmDesc := fmt.Sprintf("[%s](%s) 正在油管直播", info.ChannelName, youtube.GetChannelLink(info.ChannelId))
	dm := youtube.CreateDiscordMessage(dmDesc, info)
	go discord.SendNewsEmbed(dm)

	return youtubeSendQQRisky(info, fmt.Sprintf("%s 正在油管直播", info.ChannelName))
}

func init() {
	youtube.RegisterDataHandler(youtube.Live, HandleLiveEvent)
}
