package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/hooks/sites/twitter"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/discord"
)

func HandleTweet(_ *bot.Bot, c *twitter.TweetContent) error {

	data := c.Tweet

	discordMessage := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s 发布了一则推文", c.NickName),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "内容",
				Value: twitter.TextWithoutTCLink(data.UnEsacapedText()),
			},
		},
	}
	twitter.AddEntitiesByDiscord(discordMessage, c)
	go discord.SendNewsEmbed(discordMessage)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 发布了一则新推文", c.NickName))
	return tweetSendQQRisky(msg, data)
}

func init() {
	twitter.MessageHandler.AddHandler(twitter.Tweet, HandleTweet)
}
