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

func HandleReTweet(_ *bot.Bot, data *twitter.TweetStreamData) error {

	go handleRetweetDiscord(data, false)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 分享了一则推文", data.User.Name))
	msg.Append(qq.NextLn())

	if data.RetweetedStatus != nil {
		msg.Append(qq.NewTextLn("================="))
		if data.RetweetedStatus.User.Id != data.User.Id {
			msg.Append(qq.NextLn())
			msg.Append(qq.NewTextLn("原作者"))
			msg.Append(qq.NewTextLn(data.RetweetedStatus.User.Name))
		}
		msg.Append(qq.NextLn())
		msg.Append(qq.NewTextLn("内容"))
		return tweetSendQQRisky(msg, data.RetweetedStatus)
	} else {
		msg.Append(qq.NewTextLn("[获取转发推文失败]"))
	}

	return withRisky(msg)
}

func handleRetweetDiscord(data *twitter.TweetStreamData, withText bool) {

	first := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s 分享了一则推文", data.User.Name),
		Fields:      []*discordgo.MessageEmbedField{},
	}

	if withText {
		first.Fields = append(first.Fields, &discordgo.MessageEmbedField{
			Name:  "附文",
			Value: data.Text,
		})
	}

	if data.RetweetedStatus != nil {
		retweetedDiscordMessage := &discordgo.MessageEmbed{
			Description: data.Text,
		}
		twitter.AddEntitiesByDiscord(retweetedDiscordMessage, data.RetweetedStatus)
		discord.SendNewsEmbedDouble(first, retweetedDiscordMessage)
	}

}

func HandleReTweetWithText(_ *bot.Bot, data *twitter.TweetStreamData) error {

	go handleRetweetDiscord(data, true)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 转发了一则推文", data.User.Name))
	msg.Append(qq.NextLn())
	msg.Append(qq.NewTextLn("附文"))
	msg.Append(qq.NewTextLn(twitter.TextWithoutTCLink(data.Text)))
	msg.Append(qq.NextLn())
	if data.QuotedStatus != nil {
		msg.Append(qq.NewTextLn("================="))
		if data.QuotedStatus.User.Id != data.User.Id {
			msg.Append(qq.NextLn())
			msg.Append(qq.NewTextLn("原作者"))
			msg.Append(qq.NewTextLn(data.QuotedStatus.User.Name))
		}
		msg.Append(qq.NextLn())
		msg.Append(qq.NewTextLn("内容"))
		return tweetSendQQRisky(msg, data.QuotedStatus)
	} else {
		msg.Append(qq.NewTextLn("[获取转发推文失败]"))
	}
	return withRisky(msg)
}

func init() {
	twitter.MessageHandler.AddHandler(twitter.ReTweet, HandleReTweet)
	twitter.MessageHandler.AddHandler(twitter.ReTweetWithText, HandleReTweetWithText)
}
