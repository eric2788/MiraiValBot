package twitter

import (
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/internal/qq"
)

// CreateMessage 短視頻要單獨發送，否則無法發送原文
func CreateMessage(risk bool, msg *message.SendingMessage, data *TweetStreamData, alt ...*message.TextElement) (*message.SendingMessage, []*message.ShortVideoElement) {

	extraUrls := ExtractExtraLinks(data)

	videos := make([]*message.ShortVideoElement, 0)

	// 内文 (need to change html &gt back to >)
	msg.Append(qq.NewTextLn(data.UnEsacapedText()))

	if len(alt) > 0 {
		msg.Append(qq.NextLn())

		// 額外的中文字來減低風控機率
		for _, txt := range alt {
			msg.Append(txt)
		}
	}

	// 額外連結 (仅限非风控状态发送)
	if len(extraUrls) > 0 && !risk {
		msg.Append(qq.NextLn())
		msg.Append(qq.NewTextLn("链接"))
		for _, url := range extraUrls {
			msg.Append(qq.NewTextfLn("- %s", url))
		}
	}

	// 媒體
	if data.ExtendedEntities != nil && data.ExtendedEntities.Media != nil {
		media := *data.ExtendedEntities.Media
		for i, m := range media {
			switch m.Type {
			// 圖片
			case "photo":
				img, err := qq.NewImageByUrl(m.MediaUrlHttps)
				if err != nil {
					logger.Warnf("加载推特图片 %s 时出现错误: %v, 将尝试发送链接", m.MediaUrlHttps, err)
					msg.Append(qq.NewTextfLn("\n图片%d: %s", i+1, m.MediaUrlHttps))
				} else {
					msg.Append(img)
				}
			// Gif 圖片
			case "animated_gif":
				fallthrough
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
					videos = append(videos, video)
					success = true
					break
				}
				if !success {
					logger.Warnf("推特視頻加載失敗，將改用圖片推送。")
					img, err := qq.NewImageByUrl(m.MediaUrlHttps)
					if err != nil {
						logger.Warnf("加载推特图片 %s 时出现错误: %v, 将尝试发送链接", m.MediaUrlHttps, err)
						msg.Append(qq.NewTextf("\n视频封面\n%s", m.MediaUrlHttps))
					} else {
						msg.Append(img)
					}
				}
			default:
				logger.Warnf("未知的媒體類型: %s", m.Type)
			}
		}
	}
	return msg, videos
}

func AddEntitiesByDiscord(msg *discordgo.MessageEmbed, data *TweetStreamData) {
	msg.Author = &discordgo.MessageEmbedAuthor{
		Name:    data.User.Name,
		URL:     GetUserLink(data.User.ScreenName),
		IconURL: data.User.ProfileImageUrlHttps,
	}
	if msg.Fields == nil {
		msg.Fields = make([]*discordgo.MessageEmbedField, 0)
	}
	if data.Entities.Urls != nil && len(data.Entities.Urls) > 0 {
		urls := make([]string, 0)
		for _, url := range data.Entities.Urls {
			urls = append(urls, url.ExpandedUrl)
		}
		msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
			Name:  "链接",
			Value: strings.Join(urls, "\n"),
		})
	}

	if data.ExtendedEntities != nil && len(*data.ExtendedEntities.Media) > 0 {
		if len(*data.ExtendedEntities.Media) == 1 {
			m := (*data.ExtendedEntities.Media)[0]
			switch m.Type {
			case "photo":
				msg.Image = &discordgo.MessageEmbedImage{
					URL: m.MediaUrlHttps,
				}
			case "animated_gif":
				fallthrough
			case "video":
				for _, variant := range m.VideoInfo.Variants {
					if variant.ContentType == "video/mp4" {
						msg.Video = &discordgo.MessageEmbedVideo{
							URL: variant.Url,
						}
						break
					}
				}
			}
		} else {
			videoUrls := make([]string, 0)
			photoUrls := make([]string, 0)
			for _, m := range *data.ExtendedEntities.Media {
				switch m.Type {
				case "photo":
					photoUrls = append(photoUrls, m.MediaUrlHttps)
				case "animated_gif":
					fallthrough
				case "video":
					for _, variant := range m.VideoInfo.Variants {
						if variant.ContentType == "video/mp4" {
							videoUrls = append(videoUrls, variant.Url)
							break
						}
					}
				}
			}
			if len(videoUrls) > 0 {
				msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
					Name:  "视频",
					Value: strings.Join(videoUrls, "\n"),
				})
			}
			if len(photoUrls) > 0 {
				msg.Fields = append(msg.Fields, &discordgo.MessageEmbedField{
					Name:  "图片",
					Value: strings.Join(photoUrls, "\n"),
				})
			}
		}
	}
}
