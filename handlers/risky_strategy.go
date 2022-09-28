package handlers

import (
	"github.com/Mrs4s/MiraiGo/message"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/sites/youtube"
)

func withBilibiliRisky(msg *message.SendingMessage) (err error) {
	return qq2.SendWithRandomRiskyStrategy(msg)
}

func tweetSendQQRisky(originalMsg *message.SendingMessage, data *twitter.TweetStreamData) (err error) {

	go qq2.SendRiskyMessage(5, 60, func(try int) error {

		clone := qq2.CloneMessage(originalMsg)

		alt := qq2.GetRandomMessageByTry(try)

		msg, videos := twitter.CreateMessage(clone, data, alt...)

		// 先發送推文內容
		if err := qq2.SendGroupMessage(msg); err != nil {
			return err
		}
		// 後發送視頻訊息
		for _, video := range videos {
			if err := qq2.SendGroupMessage(message.NewSendingMessage().Append(video)); err != nil {
				return err
			}
		}

		return nil
	})
	return
}

func youtubeSendQQRisky(info *youtube.LiveInfo, desc string, blocks ...string) (err error) {

	titles := []string{"标题", "开始时间", "直播间"}

	for i, block := range blocks {
		if i == len(titles) {
			break
		}
		titles[i] = block
	}

	go qq2.SendRiskyMessage(5, 60, func(try int) error {

		noTitle := false

		alt := qq2.GetRandomMessageByTry(try)

		// 风控第四次没有标题
		if try >= 4 {
			noTitle = true
		}

		msg := youtube.CreateQQMessage(desc, info, noTitle, alt, titles...)
		return qq2.SendGroupMessage(msg)
	})

	return
}
