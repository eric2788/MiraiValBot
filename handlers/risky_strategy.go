package handlers

import (
	"github.com/Mrs4s/MiraiGo/message"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"time"
)

func cloneMessage(msg *message.SendingMessage) *message.SendingMessage {
	clone := message.NewSendingMessage()
	for _, element := range msg.Elements {
		clone.Append(element)
	}
	return clone
}

func getRandomMessageByTry(try int) []*message.TextElement {

	extras := make([]*message.TextElement, 0)

	// 新增随机发过的群消息

	if try > 0 {

		random, err := qq2.GetRandomGroupMessage(qq2.ValGroupInfo.Uin)

		if try > 2 { // 發送多一則隨機消息
			for _, element := range getRandomMessageByTry(1) { // 使用 1 確保不無限套娃
				extras = append(extras, element)
			}
		}

		if err == nil && random != nil {

			for _, element := range random.Elements {
				switch e := element.(type) {
				case *message.TextElement:
					extras = append(extras, e)
				case *message.AtElement:
					extras = append(extras, message.NewText(e.Display))
				case *message.FaceElement:
					extras = append(extras, message.NewText(e.Name))
				default:
					break
				}
			}

			// 随机消息没有文本
			if len(extras) == 0 {

				logger.Warnf("为被风控的广播插入一条新消息再发送: %s", random.ToString())

				sendFirst := message.NewSendingMessage()
				for _, element := range random.Elements {

					switch element.(type) {
					case *message.ReplyElement:
						continue
					case *message.ForwardElement:
						continue
					default:
						break
					}

					sendFirst.Append(element)
				}
				_ = qq2.SendGroupMessage(sendFirst)
				<-time.After(time.Second * 5)     // 发送完等待五秒
				return getRandomMessageByTry(try) // 再獲取一則隨機消息

			} else {

				logger.Warnf("为被风控的广播新增如下的内容: %s", random.ToString())

			}

		} else { // 随机消息获取失败

			if err != nil {
				logger.Warnf("获取随机消息时出现错误: %v, 将改为发送风控次数", err)
			} else if random == nil {
				logger.Warnf("获取随机消息时出现错误: 訊息為 nil , 将改为发送风控次数")
			}

			// 则发送风控次数?
			extras = append(extras, qq2.NewTextf("此广播已被风控 %d 次 QAQ!!", try))

		}

	}

	return extras
}

func withBilibiliRisky(msg *message.SendingMessage) (err error) {
	go qq2.SendRiskyMessage(5, 60, func(try int) error {

		clone := cloneMessage(msg)

		alt := getRandomMessageByTry(try)

		if len(alt) > 0 {
			clone.Append(qq2.NextLn())
			for _, element := range alt {
				clone.Append(element)
			}
		}

		return qq2.SendGroupMessage(clone)

	})
	return
}

func tweetSendQQRisky(originalMsg *message.SendingMessage, data *twitter.TweetStreamData) (err error) {

	go qq2.SendRiskyMessage(5, 60, func(try int) error {

		clone := cloneMessage(originalMsg)

		alt := getRandomMessageByTry(try)

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

		alt := getRandomMessageByTry(try)

		// 风控第四次没有标题
		if try >= 4 {
			noTitle = true
		}

		msg := youtube.CreateQQMessage(desc, info, noTitle, alt, titles...)
		return qq2.SendGroupMessage(msg)
	})

	return
}
