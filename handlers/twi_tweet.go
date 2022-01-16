package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleTweet(_ *bot.Bot, data *twitter.TweetStreamData) error {

	discordMessage := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s 发布了一则贴文", data.User.Name),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "内容",
				Value: twitter.TextWithoutTCLink(data.Text),
			},
		},
	}
	twitter.AddEntitiesByDiscord(discordMessage, data)
	go discord.SendNewsEmbed(discordMessage)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 发布了一则新贴文", data.User.Name))
	return tweetSendQQRisky(msg, data)
}

func init() {
	twitter.RegisterDataHandler(twitter.Tweet, HandleTweet)
}
