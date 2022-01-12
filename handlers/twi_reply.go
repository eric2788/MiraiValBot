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

func HandleTweetReply(bot *bot.Bot, data *twitter.TweetStreamData) error {

	handleTweetDiscord(data)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 回复了 %s 的一则推文", data.User.Name, *data.InReplyToScreenName))
	msg.Append(qq.NewTextfLn("内容: %s", data.Text))
	msg.Append(qq.NewTextf("回复贴文: https://twitter.com/%s/status/%s", *data.InReplyToScreenName, data.InReplyToStatusIdStr))

	return withRisky(msg)
}

func handleTweetDiscord(data *twitter.TweetStreamData) {
	msg := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    data.User.Name,
			URL:     "https://twitter.com/" + data.User.ScreenName,
			IconURL: data.User.ProfileImageUrlHttps,
		},
		Description: fmt.Sprintf("%s 回复了 %s 的一则推文", data.User.Name, *data.InReplyToScreenName),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "回复贴文",
				Value: fmt.Sprintf("https://twitter.com/%s/status/%s", *data.InReplyToScreenName, data.InReplyToStatusIdStr),
			},
			{
				Name:  "回复内容",
				Value: data.Text,
			},
		},
	}
	addEntitiesTweetDiscord(msg, data)
	go discord.SendNewsEmbed(msg)
}

// withRisky error must be nil
func withRisky(msg *message.SendingMessage) (err error) {
	go qq.SendRiskyMessage(5, 10, func() error {
		return qq.SendGroupMessage(msg)
	})
	return
}

func init() {
	twitter.RegisterDataHandler(twitter.Reply, HandleTweetReply)
}
