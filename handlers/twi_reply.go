package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleTweetReply(bot *bot.Bot, data *twitter.TweetStreamData) error {
	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 回复了 %s 的一则推文", data.User.Name, *data.InReplyToScreenName))
	msg.Append(qq.NewTextfLn("内容: %s", data.Text))
	msg.Append(qq.NewTextf("回复贴文: https://twitter.com/%s/status/%s", *data.InReplyToScreenName, data.InReplyToStatusIdStr))

	bot.SendGroupMessage(qq.ValGroupInfo.Uin, msg)
	return nil
}

func init() {
	twitter.RegisterDataHandler(twitter.Reply, HandleTweetReply)
}
