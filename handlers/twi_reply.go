package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/file"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/sites/twitter"
)

func HandleTweetReply(_ *bot.Bot, data *twitter.TweetStreamData) error {

	// 设置了不推送推文回复
	if !file.DataStorage.Twitter.ShowReply {
		return nil
	}

	discordMsg := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("[%s](%s) 回复了 [%s](%s) 的一则推文", data.User.Name, twitter.GetUserLink(data.User.ScreenName), *data.InReplyToScreenName, twitter.GetUserLink(*data.InReplyToScreenName)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "回复贴文",
				Value: twitter.GetStatusLink(*data.InReplyToScreenName, data.InReplyToStatusIdStr),
			},
			{
				Name:  "回复内容",
				Value: twitter.TextWithoutTCLink(data.Text),
			},
		},
	}
	twitter.AddEntitiesByDiscord(discordMsg, data)
	go discord.SendNewsEmbed(discordMsg)

	msg := message.NewSendingMessage()
	msg.Append(qq2.NewTextfLn("%s 回复了 %s 的一则推文", data.User.Name, *data.InReplyToScreenName))
	msg.Append(qq2.NextLn())
	msg.Append(qq2.NewTextLn("回复贴文"))
	msg.Append(qq2.NewTextfLn("https://twitter.com/%s/status/%s", *data.InReplyToScreenName, data.InReplyToStatusIdStr))
	msg.Append(qq2.NextLn())
	msg.Append(qq2.NewTextLn("内容"))
	return tweetSendQQRisky(msg, data)
}

// withRisky error must be nil
func withRisky(msg *message.SendingMessage) (err error) {
	go qq2.SendRiskyMessage(5, 60, func(try int) error {
		return qq2.SendGroupMessage(msg)
	})
	return
}

func init() {
	twitter.RegisterDataHandler(twitter.Reply, HandleTweetReply)
}
