package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleLive(bot *bot.Bot, data *bilibili.LiveData) error {

	handleLiveDiscord(data)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 正在B站直播", data.LiveInfo.Name))
	msg.Append(qq.NewTextfLn("标题: %s", data.LiveInfo.Title))
	msg.Append(qq.NewTextfLn("直播间: %s", biliRoomLink(data.LiveInfo.RoomId)))
	if data.LiveInfo.Cover != nil && *data.LiveInfo.Cover != "" {
		cover := *data.LiveInfo.Cover
		imgElement, err := qq.NewImageByUrl(cover)
		if err != nil {
			logger.Warnf("获取图片 %s 时出现错误: %v", cover, err)
		} else {
			msg.Append(imgElement)
		}
	}
	return withRisky(msg)
}

func handleLiveDiscord(data *bilibili.LiveData) {

	discordMessage := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("[%s](%s) 正在B站直播", data.LiveInfo.Name, biliSpaceLink(data.LiveInfo.UID)),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "标题",
				Value: data.LiveInfo.Title,
			},
			{
				Name:  "直播间",
				Value: biliRoomLink(data.LiveInfo.RoomId),
			},
		},
	}

	if data.LiveInfo.Cover != nil {
		discordMessage.Image = &discordgo.MessageEmbedImage{
			URL: *data.LiveInfo.Cover,
		}
	}

	go discord.SendNewsEmbed(discordMessage)
}

func init() {
	bilibili.RegisterDataHandler(bilibili.Live, HandleLive)
}
