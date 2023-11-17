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

func HandleReTweet(_ *bot.Bot, data *twitter.TweetContent) error {

	go handleRetweetDiscord(data.Tweet, false)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 分享了一则推文", data.NickName))
	msg.Append(qq.NextLn())

	if data.Tweet.RetweetedStatus != nil {
		msg.Append(qq.NewTextLn("================="))
		if data.Tweet.RetweetedStatus.UserID != data.Tweet.UserID {
			msg.Append(qq.NextLn())
			msg.Append(qq.NewTextLn("原作者"))
			msg.Append(qq.NewTextLn(data.Tweet.RetweetedStatus.Name))
		}
		msg.Append(qq.NextLn())
		msg.Append(qq.NewTextLn("内容"))
		return tweetSendQQRisky(msg, data.Tweet.RetweetedStatus)
	} else {
		msg.Append(qq.NewTextLn("[获取转发推文失败]"))
	}

	return withRisky(msg)
}

func handleRetweetDiscord(data *twitter.TweetData, withText bool) {

	first := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("%s 分享了一则推文", data.Name),
		Fields:      []*discordgo.MessageEmbedField{},
	}

	if withText {
		first.Fields = append(first.Fields, &discordgo.MessageEmbedField{
			Name:  "附文",
			Value: data.UnEsacapedText(),
		})
	}

	if data.RetweetedStatus != nil {
		retweetedDiscordMessage := &discordgo.MessageEmbed{
			Description: data.UnEsacapedText(),
		}
		twitter.AddEntitiesByDiscord(retweetedDiscordMessage, &twitter.TweetContent{
			Tweet: data.RetweetedStatus,
			NickName: data.RetweetedStatus.Name,
			// profile: how to get ?
		})
		discord.SendNewsEmbedDouble(first, retweetedDiscordMessage)
	}

}

func HandleReTweetWithText(_ *bot.Bot, data *twitter.TweetContent) error {

	go handleRetweetDiscord(data.Tweet, true)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 转发了一则推文", data.NickName))
	msg.Append(qq.NextLn())
	msg.Append(qq.NewTextLn("附文"))
	msg.Append(qq.NewTextLn(twitter.TextWithoutTCLink(data.Tweet.UnEsacapedText())))
	msg.Append(qq.NextLn())
	if data.Tweet.QuotedStatus != nil {
		msg.Append(qq.NewTextLn("================="))
		if data.Tweet.QuotedStatus.UserID != data.Tweet.UserID {
			msg.Append(qq.NextLn())
			msg.Append(qq.NewTextLn("原作者"))
			msg.Append(qq.NewTextLn(data.Tweet.QuotedStatus.Name))
		}
		msg.Append(qq.NextLn())
		msg.Append(qq.NewTextLn("内容"))
		return tweetSendQQRisky(msg, data.Tweet.QuotedStatus)
	} else {
		msg.Append(qq.NewTextLn("[获取转发推文失败]"))
	}
	return withRisky(msg)
}

func init() {
	twitter.MessageHandler.AddHandler(twitter.ReTweet, HandleReTweet)
	twitter.MessageHandler.AddHandler(twitter.ReTweetWithText, HandleReTweetWithText)
}
