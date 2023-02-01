package handlers

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/hooks/sites/twitter"
	"github.com/eric2788/MiraiValBot/hooks/sites/youtube"
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/valorant"
	"github.com/eric2788/common-utils/datetime"
)

func withBilibiliRisky(msg *message.SendingMessage) (err error) {
	return qq.SendWithRandomRiskyStrategy(msg)
}

func tweetSendQQRisky(originalMsg *message.SendingMessage, data *twitter.TweetStreamData) (err error) {

	go qq.SendRiskyMessage(5, 60, func(try int) error {

		clone := qq.CloneMessage(originalMsg)

		alt := qq.GetRandomMessageByTry(try)

		msg, videos := twitter.CreateMessage(try >= 3, clone, data, alt...)

		if try == 1 {
			if err := qq.SendGroupImageText(msg); err != nil {
				return err
			}
		} else {
			// 先發送推文內容
			if err := qq.SendGroupMessage(msg); err != nil {
				return err
			}
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

		noTitle := false

		alt := qq.GetRandomMessageByTry(try)

		// 风控第四次没有标题
		if try >= 4 {
			noTitle = true
		}

		msg := youtube.CreateQQMessage(desc, info, noTitle, alt, titles...)

		// 尝试发送一次图片信息
		if try == 1 {
			return qq.SendGroupImageText(msg)
		}

		return qq.SendGroupMessage(msg)
	})

	return
}

func valorantTrackRisky(displayName, shortHint, cmdId string, metaData *valorant.MatchMetaData) (err error) {

	go qq.SendRiskyMessage(5, 60, func(currentTry int) error {

		// 尝试缩短对战ID
		if currentTry >= 3 && cmdId != metaData.MatchId {
			metaData.MatchId = cmdId
			shortHint = "(已缩短)"
		}

		msg := message.NewSendingMessage()
		msg.Append(qq.NewTextfLn("%s 的最新对战信息已更新👇", displayName))
		msg.Append(qq.NewTextfLn("对战ID: %s%s", metaData.MatchId, shortHint))
		if currentTry <= 4 {
			msg.Append(qq.NewTextfLn("对战模式: %s", metaData.Mode))
			msg.Append(qq.NewTextfLn("对战开始时间: %s", datetime.FormatSeconds(metaData.GameStart)))
			msg.Append(qq.NewTextfLn("对战地图: %s", metaData.Map))
		}
		msg.Append(qq.NewTextfLn("输入 !val match %s 查看详细内容。", cmdId))

		// 尝试发送一次图片信息
		if currentTry == 1 {
			return qq.SendGroupImageText(msg)
		}

		alt := qq.GetRandomMessageByTry(currentTry)

		if len(alt) > 0 {
			msg.Append(qq.NextLn())
		}

		for _, ele := range alt {
			msg.Append(ele)
			msg.Append(qq.NextLn())
		}

		return qq.SendGroupMessage(msg)
	})

	return
}
