package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/sites/youtube"
)

func HandleIdle(bot *bot.Bot, info *youtube.LiveInfo) error {

	go discord.SendNewsEmbed(&discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:  youtube.GetChannelLink(info.ChannelId),
			Name: info.ChannelName,
		},
		Description: fmt.Sprintf("%s 的油管直播已结束。", info.ChannelName),
	})

	msg := message.NewSendingMessage().Append(qq2.NewTextf("%s 的油管直播已结束。", info.ChannelName))
	return qq2.SendGroupMessage(msg) // 一句，无需管理风控
}

func init() {
	youtube.RegisterDataHandler(youtube.Idle, HandleIdle)
}
