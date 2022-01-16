package handlers

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"strings"
)

func withBilibiliRisky(msg *message.SendingMessage) (err error) {
	go qq.SendRiskyMessage(5, 60, func(try int) error {
		clone := message.NewSendingMessage()
		for _, element := range msg.Elements {
			clone.Append(element)
		}

		alt := make([]string, 0)

		// 风控时尝试加随机文字看看会不会减低？

		if try >= 1 {
			alt = append(alt, fmt.Sprintf("[此广播已被风控 %d 次]", try))
		}

		if try >= 2 {
			alt = append(alt, fmt.Sprintf("你好谢谢小笼包再见"))
		}

		if try >= 3 {
			alt = append(alt, fmt.Sprintf("卧槽，这个直播真牛逼!"))
		}

		if try >= 4 {
			alt = append(alt, fmt.Sprintf("哟，风控四次了，这直播是个啥啊？"))
		}

		if len(alt) > 0 {
			logger.Warnf("为被风控的推文新增如下的内容: %s", strings.Join(alt, "\n"))
			clone.Append(qq.NextLn())
			for _, al := range alt {
				clone.Append(qq.NewTextLn(al))
			}
		}

		return qq.SendGroupMessage(clone)

	})
	return
}

func tweetSendQQRisky(originalMsg *message.SendingMessage, data *twitter.TweetStreamData) (err error) {

	go qq.SendRiskyMessage(5, 60, func(try int) error {

		clone := message.NewSendingMessage()

		for _, element := range originalMsg.Elements {
			clone.Append(element)
		}

		alt := make([]string, 0)

		// 风控时尝试加随机文字看看会不会减低？

		if try >= 1 {
			alt = append(alt, fmt.Sprintf("[此推文已被风控 %d 次]", try))
		}

		if try >= 2 {
			alt = append(alt, fmt.Sprintf("你好谢谢小笼包再见"))
		}

		if try >= 3 {
			alt = append(alt, fmt.Sprintf("卧槽，这个推文真牛逼!"))
		}

		if try >= 4 {
			alt = append(alt, fmt.Sprintf("哟，风控四次了，这推文会不会是在GHS啊？"))
		}

		if try > 0 {
			logger.Warnf("为被风控的推文新增如下的内容: %s", strings.Join(alt, "\n"))
		}

		msg, videos := twitter.CreateMessage(clone, data, alt...)

		// 先發送推文內容
		if err := qq.SendGroupMessage(msg); err != nil {
			return err
		}
		// 後發送視頻訊息
		for _, video := range videos {
			if err := qq.SendGroupMessage(message.NewSendingMessage().Append(video)); err != nil {
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

	go qq.SendRiskyMessage(5, 60, func(try int) error {
		alt := make([]string, 0)
		noTitle := false

		// 风控时尝试加随机文字看看会不会减低？

		if try >= 1 {
			alt = append(alt, fmt.Sprintf("[此油管广播已被风控 %d 次]", try))
		}

		if try >= 2 {
			alt = append(alt, fmt.Sprintf("你好谢谢小笼包再见"))
		}

		if try >= 3 {
			alt = append(alt, fmt.Sprintf("看看啊！这谁直播了额!"))
		}

		// 风控第四次没有标题
		if try >= 4 {
			alt = append(alt, fmt.Sprintf("什么油管直播居然被风控了四次这么爽？？你标题没了"))
			noTitle = true
		}

		if try > 0 {
			logger.Warnf("为被风控的推文新增如下的内容: %s", strings.Join(alt, "\n"))
		}

		msg := youtube.CreateQQMessage(desc, info, noTitle, alt, titles...)
		return qq.SendGroupMessage(msg)
	})

	return
}
