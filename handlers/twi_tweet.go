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
				Value: data.Text,
			},
		},
	}
	addEntitiesTweetDiscord(discordMessage, data)
	go discord.SendNewsEmbed(discordMessage)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 发布了一则新贴文", data.User.Name))
	createTweetMessage(msg, data)

	return withRisky(msg)
}

func tweetUserLink(screen string) string {
	return fmt.Sprintf("https://twitter.com/%s", screen)
}

func tweetStatusLink(screen, status string) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", screen, status)
}

func addEntitiesTweetDiscord(msg *discordgo.MessageEmbed, data *twitter.TweetStreamData) {
	msg.Author = &discordgo.MessageEmbedAuthor{
		Name:    data.User.Name,
		URL:     tweetUserLink(data.User.ScreenName),
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

	if data.ExtendedEntities.Media != nil && len(*data.ExtendedEntities.Media) > 0 {
		if len(*data.ExtendedEntities.Media) == 1 {
			m := (*data.ExtendedEntities.Media)[0]
			switch m.Type {
			case "photo":
				msg.Image = &discordgo.MessageEmbedImage{
					URL: m.MediaUrlHttps,
				}
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
