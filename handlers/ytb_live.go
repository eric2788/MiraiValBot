package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/youtube"
	"github.com/eric2788/MiraiValBot/utils/datetime"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleLiveEvent(bot *bot.Bot, info *youtube.LiveInfo) error {

	dmDesc := fmt.Sprintf("[%s](%s) 正在油管直播", info.ChannelName, youtubeChannelLink(info.ChannelId))
	dm := createYoutubeMessageDiscord(dmDesc, info)
	go discord.SendNewsEmbed(dm)

	msg := createYoutubeMessage(fmt.Sprintf("%s 正在油管直播", info.ChannelName), info)
	return withRisky(msg)
}

func createYoutubeMessageDiscord(desc string, info *youtube.LiveInfo, fields ...string) *discordgo.MessageEmbed {

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
			URL:  youtubeChannelLink(info.ChannelId),
			Name: info.ChannelName,
		},
		Description: desc,
		Fields:      []*discordgo.MessageEmbedField{},
	}

	if info.Info != nil {
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
				Value: datetime.FormatMillis(info.Info.PublishTime),
			}, &discordgo.MessageEmbedField{
				Name:  blocks[3],
				Value: getYTLink(info),
			})
	} else {
		dm.Fields = append(dm.Fields, &discordgo.MessageEmbedField{
			Name:  blocks[3],
			Value: getYTLink(info),
		})
	}

	return dm
}

func createYoutubeMessage(desc string, info *youtube.LiveInfo, fields ...string) *message.SendingMessage {

	blocks := []string{
		"开始时间",
		"标题",
		"直播间",
	}

	for i, f := range fields {
		blocks[i] = f
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextLn(desc))

	if info.Info != nil {
		msg.Append(qq.NewTextfLn("%s: %s", blocks[0], info.Info.Title))
		msg.Append(qq.NewTextfLn("%s: %s", blocks[1], datetime.FormatMillis(info.Info.PublishTime)))
		msg.Append(qq.NewTextfLn("%s: %s", blocks[2], getYTLink(info)))

		if info.Info.Cover != nil {
			cover := *info.Info.Cover
			img, err := qq.NewImageByUrl(cover)
			if err != nil {
				logger.Warnf("获取图片 %s 时出现错误: %v", cover, err)
			} else {
				msg.Append(img)
			}
		}

	} else {
		msg.Append(qq.NewTextf("%s: %s", blocks[2], getYTLink(info)))
	}

	return msg
}

func youtubeChannelLink(id string) string {
	return fmt.Sprintf("https://youtube.com/channel/%s", id)
}

func getYTLink(info *youtube.LiveInfo) string {
	if info.Info != nil {
		return fmt.Sprintf("https://youtu.be/%s", info.Info.Id)
	} else {
		return fmt.Sprintf("https://youtube.com/channel/%s/live", info.ChannelId)
	}
}

func init() {
	youtube.RegisterDataHandler(youtube.Live, HandleLiveEvent)
}
