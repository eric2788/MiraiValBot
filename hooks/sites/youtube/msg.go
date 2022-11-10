package youtube

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/common-utils/datetime"
)

func CreateQQMessage(desc string, info *LiveInfo, noTitle bool, alt []*message.TextElement, fields ...string) *message.SendingMessage {

	blocks := []string{"标题", "开始时间", "直播间"}

	for i, field := range fields {
		if i == len(blocks) {
			break
		}

		blocks[i] = field
	}

	title, startTime, roomLink := blocks[0], blocks[1], blocks[2]

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextLn(desc))

	if info.Info != nil {

		if !noTitle {
			msg.Append(qq.NewTextfLn("%s: %s", title, info.Info.Title))
		}

		t, err := datetime.ParseISOStr(info.Info.PublishTime)
		if err == nil {
			msg.Append(qq.NewTextfLn("%s: %s", startTime, datetime.FormatMillis(t.UnixMilli())))
		} else {
			logger.Warnf("解析時間文字 %s 時出現錯誤: %v", info.Info.PublishTime, err)
			msg.Append(qq.NewTextfLn("%s: %s", startTime, info.Info.PublishTime)) // 使用原本的 string
		}

		msg.Append(qq.NewTextfLn("%s: %s", roomLink, GetYTLink(info)))

		if info.Info.Cover != nil && *info.Info.Cover != "" {
			cover := *info.Info.Cover
			img, err := qq.NewImageByUrl(cover)
			if err != nil {
				logger.Warnf("获取图片 %s 时出现错误: %v", cover, err)
			} else {
				msg.Append(img)
			}
		}

	} else {
		msg.Append(qq.NewTextf("%s: %s", roomLink, GetYTLink(info)))
	}

	// 随机文字
	if len(alt) > 0 {
		msg.Append(qq.NextLn())
		for _, a := range alt {
			msg.Append(a)
		}
	}

	return msg
}

func CreateDiscordMessage(desc string, info *LiveInfo, fields ...string) *discordgo.MessageEmbed {

	blocks := []string{
		"开始时间",
		"标题",
		"直播间",
	}

	for i, f := range fields {
		if i == len(blocks) {
			break
		}
		blocks[i] = f
	}

	dm := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			URL:  GetChannelLink(info.ChannelId),
			Name: info.ChannelName,
		},
		Description: desc,
		Fields:      []*discordgo.MessageEmbedField{},
	}

	if info.Info != nil {

		var publishTime string
		t, err := datetime.ParseISOStr(info.Info.PublishTime)
		if err == nil {
			publishTime = datetime.FormatMillis(t.UnixMilli())
		} else {
			logger.Warnf("解析時間文字 %s 時出現錯誤: %v", info.Info.PublishTime, err)
			publishTime = info.Info.PublishTime // 使用原本的 string
		}

		dm.Fields = append(dm.Fields,
			&discordgo.MessageEmbedField{
				Name:  blocks[0],
				Value: publishTime,
			}, &discordgo.MessageEmbedField{
				Name:  blocks[1],
				Value: info.Info.Title,
			}, &discordgo.MessageEmbedField{
				Name:  blocks[2],
				Value: GetYTLink(info),
			})

		if info.Info.Cover != nil && *info.Info.Cover != "" {
			cover := *info.Info.Cover
			dm.Image = &discordgo.MessageEmbedImage{
				URL: cover,
			}
		}
	} else {
		dm.Fields = append(dm.Fields, &discordgo.MessageEmbedField{
			Name:  blocks[2],
			Value: GetYTLink(info),
		})
	}

	return dm
}
