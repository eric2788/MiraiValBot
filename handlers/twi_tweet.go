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

func tweetSendQQRisky(originalMsg *message.SendingMessage, data *twitter.TweetStreamData) (err error) {

	go qq.SendRiskyMessage(5, 10, func(try int) error {

		clone := message.NewSendingMessage()

		for _, element := range originalMsg.Elements {
			clone.Append(element)
		}

		var alt []string

		// 风控时尝试加随机文字看看会不会减低？

		if try > 1 {
			alt = append(alt, fmt.Sprintf("(此推文已被风控 %d 次)", try))
		}

		if try > 2 {
			alt = append(alt, fmt.Sprintf("你好谢谢小笼包再见"))
		}

		if try > 3 {
			alt = append(alt, fmt.Sprintf("卧槽，这个推文真牛逼!"))
		}

		if try > 4 {
			alt = append(alt, fmt.Sprintf("哟，风控四次了，这推文会不会是在GHS啊？"))
		}

		logger.Warnf("为被风控的推文新增如下的内容: %s", strings.Join(alt, "\n"))

		msg := twitter.CreateMessage(clone, data, alt...)
		return qq.SendGroupMessage(msg)
	})
	return
}

func init() {
	twitter.RegisterDataHandler(twitter.Tweet, HandleTweet)
}
