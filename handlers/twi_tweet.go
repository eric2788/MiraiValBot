package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/twitter"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleTweet(bot *bot.Bot, data *twitter.TweetStreamData) error {

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 发布了一则新贴文", data.User.Name))
	createTweetMessage(msg, data)

	bot.SendGroupMessage(qq.ValGroupInfo.Uin, msg)
	return nil
}

func createTweetMessage(msg *message.SendingMessage, data *twitter.TweetStreamData) {
	// 内文
	msg.Append(qq.NewTextLn(data.Text))
	// 連結
	if data.Entities.Urls != nil && len(data.Entities.Urls) > 0 {
		msg.Append(qq.NewTextLn("链接: "))
		for _, url := range data.Entities.Urls {
			msg.Append(qq.NewTextfLn("- %s", url.ExpandedUrl))
		}
	}
	// 媒體
	if data.ExtendedEntities != nil && data.ExtendedEntities.Media != nil {
		media := *data.ExtendedEntities.Media
		for _, m := range media {
			switch m.Type {
			// 圖片
			case "photo":
				img, err := qq.NewImageByUrl(m.MediaUrlHttps)
				if err != nil {
					logger.Warnf("加载推特图片 %s 时出现错误: %v", m.MediaUrlHttps, err)
				} else {
					msg.Append(img)
				}
			// 視頻
			case "video":
				videoInfo := m.VideoInfo
				success := false
				for _, variant := range videoInfo.Variants {
					if variant.ContentType != "video/mp4" {
						continue
					}
					video, err := qq.NewVideoByUrl(variant.Url, m.MediaUrlHttps)
					if err != nil {
						logger.Warnf("加載推特視頻 %s 時出現錯誤: %v, 尋找下一個線路。", variant.Url, err)
						continue
					}
					msg.Append(video)
					success = true
					break
				}
				if !success {
					logger.Warnf("推特視頻加載失敗，將改用圖片推送。")
					img, err := qq.NewImageByUrl(m.MediaUrlHttps)
					if err != nil {
						logger.Warnf("加载推特图片 %s 时出现错误: %v", m.MediaUrlHttps, err)
					} else {
						msg.Append(img)
					}
				}
			default:
				logger.Warnf("未知的媒體類型: %s", m.Type)
			}
		}
	}
}

func init() {
	twitter.RegisterDataHandler(twitter.Tweet, HandleTweet)
}
