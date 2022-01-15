package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/discord"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"strings"
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
	return withBilibiliRisky(msg)
}

func withBilibiliRisky(msg *message.SendingMessage) (err error) {
	go qq.SendRiskyMessage(5, 10, func(try int) error {
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
