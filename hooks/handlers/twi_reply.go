package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/hooks/sites/twitter"
	"github.com/eric2788/MiraiValBot/internal/file"
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/discord"
)

func HandleTweetReply(_ *bot.Bot, data *twitter.TweetContent) error {

	// 设置了不推送推文回复
	if !file.DataStorage.Twitter.ShowReply {
		return nil
	}

	discordMsg := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("[%s](%s) 回复了 [%s](%s) 的一则推文", data.NickName, twitter.GetUserLink(data.Profile.Username), data.Tweet.InReplyToStatus.Name, twitter.GetUserLink(data.Tweet.InReplyToStatus.Username)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "回复贴文",
				Value: twitter.GetStatusLink(data.Tweet.InReplyToStatus.Username, data.Tweet.InReplyToStatusID),
			},
			{
				Name:  "回复内容",
				Value: twitter.TextWithoutTCLink(data.Tweet.UnEsacapedText()),
			},
		},
	}
	twitter.AddEntitiesByDiscord(discordMsg, data)
	go discord.SendNewsEmbed(discordMsg)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 回复了 %s 的一则推文", data.NickName, data.Tweet.InReplyToStatus.Name))
	msg.Append(qq.NextLn())
	msg.Append(qq.NewTextLn("回复贴文"))
	msg.Append(qq.NewTextfLn("https://twitter.com/%s/status/%s", data.Tweet.InReplyToStatus.Username, data.Tweet.InReplyToStatusID))
	msg.Append(qq.NextLn())
	msg.Append(qq.NewTextLn("内容"))
	return tweetSendQQRisky(msg, data.Tweet)
}

// withRisky error must be nil
func withRisky(msg *message.SendingMessage) (err error) {
	go qq.SendRiskyMessage(5, 60, func(try int) error {
		return qq.SendGroupMessage(msg)
	})
	return
}

func init() {
	twitter.MessageHandler.AddHandler(twitter.Reply, HandleTweetReply)
}
