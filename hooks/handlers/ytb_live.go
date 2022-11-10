package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/hooks/sites/youtube"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/services/discord"
)

func HandleLiveEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	if info.Duplicate && file.DataStorage.Youtube.AntiDuplicate {
		return nil
	}

	dmDesc := fmt.Sprintf("[%s](%s) 正在油管直播", info.ChannelName, youtube.GetChannelLink(info.ChannelId))
	dm := youtube.CreateDiscordMessage(dmDesc, info)
	go discord.SendNewsEmbed(dm)

	return youtubeSendQQRisky(info, fmt.Sprintf("%s 正在油管直播", info.ChannelName))
}

func init() {
	youtube.MessageHandler.AddHandler(youtube.Live, HandleLiveEvent)
}
