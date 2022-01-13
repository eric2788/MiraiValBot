package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"strings"
	"time"
)

func HandleTweet(bot *bot.Bot, data *twitter.TweetStreamData) error {

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

func tweetSendQQRisky(msg *message.SendingMessage, data *twitter.TweetStreamData) (err error) {
	go qq.SendRiskyMessage(5, time.Second*10, func(try int) error {
		shows := []bool{
			true, // 視頻
			true, // 圖片
			true, // 鏈接
			true, // 內文
		}

		/*
			风控一次，没有视频
			风控两次，没有图片
			风控三次，没有链接
		*/
		if try > 0 {
			for i := 0; i < try-1; i++ {
				shows[i] = false
			}
			logger.Warnf("推特广播被风控 %d 次，舍弃 %v 重发", try, strings.Join([]string{"視頻", "圖片", "鏈接", ""}[0:try-1], ", "))
		}
		msg := twitter.CreateMessage(msg, data, shows...)
		return qq.SendGroupMessage(msg)
	})
	return
}

func init() {
	twitter.RegisterDataHandler(twitter.Tweet, HandleTweet)
}
