package youtube

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/utils/datetime"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func CreateQQMessage(desc string, info *LiveInfo, image bool, roomTitle string, fields map[string]string) *message.SendingMessage {

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextLn(desc))

	if info.Info != nil {

		for field, value := range fields {
			msg.Append(qq.NewTextfLn("%s: %s", field, value))
		}

		if info.Info.Cover != nil && image {
			cover := *info.Info.Cover
			img, err := qq.NewImageByUrl(cover)
			if err != nil {
				logger.Warnf("获取图片 %s 时出现错误: %v", cover, err)
			} else {
				msg.Append(img)
			}
		}

	} else {
		msg.Append(qq.NewTextf("%s: %s", roomTitle, GetYTLink(info)))
	}

	return msg
}

func CreateDiscordMessage(desc string, info *LiveInfo, fields ...string) *discordgo.MessageEmbed {

	blocks := []string{
		"开始时间",
		"标题",
		"描述",
		"直播间",
	}

	for i, f := range fields {
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
		if err != nil {
			publishTime = datetime.FormatMillis(t.UnixMilli())
		} else {
			logger.Warnf("解析時間文字 %s 時出現錯誤: %v", info.Info.PublishTime, err)
			publishTime = info.Info.PublishTime // 使用原本的 string
		}

		dm.Fields = append(dm.Fields,
			&discordgo.MessageEmbedField{
				Name:   blocks[0],
				Value:  info.Info.Title,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:   blocks[1],
				Value:  info.Info.Description,
				Inline: true,
			}, &discordgo.MessageEmbedField{
				Name:  blocks[2],
				Value: publishTime,
			}, &discordgo.MessageEmbedField{
				Name:  blocks[3],
				Value: GetYTLink(info),
			})
	} else {
		dm.Fields = append(dm.Fields, &discordgo.MessageEmbedField{
			Name:  blocks[3],
			Value: GetYTLink(info),
		})
	}

	return dm
}
