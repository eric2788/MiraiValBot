package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/hooks/sites/youtube"
	"github.com/eric2788/MiraiValBot/internal/file"
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/discord"
)

func HandleIdle(bot *bot.Bot, info *youtube.LiveInfo) error {

	// if true, don't broadcast stream end
	if !file.DataStorage.Youtube.BroadcastIdle {
		return nil
	}

	go discord.SendNewsEmbed(&discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:  youtube.GetChannelLink(info.ChannelId),
			Name: info.ChannelName,
		},
		Description: fmt.Sprintf("%s 的油管直播已结束。", info.ChannelName),
	})

	msg := message.NewSendingMessage().Append(qq.NewTextf("%s 的油管直播已结束。", info.ChannelName))
	return qq.SendGroupMessage(msg) // 一句，无需管理风控
}

func init() {
	youtube.MessageHandler.AddHandler(youtube.Idle, HandleIdle)
}
