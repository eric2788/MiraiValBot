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

func HandleReTweet(bot *bot.Bot, data *twitter.TweetStreamData) error {

	go handleRetweetDiscord(data, false)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 分享了一则推文", data.User.Name))
	if data.RetweetedStatus != nil {
		msg.Append(qq.NewTextfLn("转发推文: "))
		createTweetMessage(msg, data.RetweetedStatus)
	}

	return withRisky(msg)
}

func handleRetweetDiscord(data *twitter.TweetStreamData, withText bool) {

	msg := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s 分享了一则推文", data.User.Name),
		Fields:      []*discordgo.MessageEmbedField{},
	}

	if withText {
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "附文",
			Value: data.Text,
		})
	}

	discord.SendNewsEmbed(msg)

	if data.RetweetedStatus != nil {
		retweetedDiscordMessage := &discordgo.MessageEmbed{
			Description: data.Text,
		}
		addEntitiesTweetDiscord(retweetedDiscordMessage, data.RetweetedStatus)
		discord.SendNewsEmbed(retweetedDiscordMessage)
	}

}

func HandleReTweetWithText(bot *bot.Bot, data *twitter.TweetStreamData) error {

	go handleRetweetDiscord(data, true)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 转发了一则推文", data.User.Name))
	msg.Append(qq.NewTextfLn("附文: %s", data.Text))
	if data.QuotedStatus != nil {
		msg.Append(qq.NewTextfLn("转发推文: "))
		createTweetMessage(msg, data.QuotedStatus)
	}
	return withRisky(msg)
}

func init() {
	twitter.RegisterDataHandler(twitter.ReTweet, HandleReTweet)
	twitter.RegisterDataHandler(twitter.ReTweetWithText, HandleReTweetWithText)
}
