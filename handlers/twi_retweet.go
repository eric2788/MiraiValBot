package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleReTweet(bot *bot.Bot, data *twitter.TweetStreamData) error {
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 分享了一则推文", data.User.Name))
	if data.RetweetedStatus != nil {
		msg.Append(qq.NewTextfLn("转发推文: "))
		createTweetMessage(msg, data.RetweetedStatus)
	}
	return withRisky(msg)
}

func HandleReTweetWithText(bot *bot.Bot, data *twitter.TweetStreamData) error {
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 分享了一则推文", data.User.Name))
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
