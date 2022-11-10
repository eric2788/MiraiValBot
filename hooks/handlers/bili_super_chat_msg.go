package handlers

import (
	"fmt"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/bwmarrin/discordgo"
	"github.com/eric2788/MiraiValBot/hooks/sites/bilibili"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/services/discord"
)

func HandleSuperChatMsg(bot *bot.Bot, data *bilibili.LiveData) error {
	d := data.Content["data"]

	var superchat = &bilibili.SuperChatMessageData{}

	if dict, ok := d.(map[string]interface{}); ok {
		if err := superchat.Parse(dict); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("解析 SuperChat 數據失敗")
	}

	if !bilibili.IsHighlighter(superchat.UID) {
		return nil
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("在 %s 的直播间收到来自 %s 的醒目留言", data.LiveInfo.Name, superchat.UserInfo.UName))
	msg.Append(qq.NewTextfLn("￥ %d", superchat.Price))
	msg.Append(qq.NewTextf("「%s」", superchat.Message))

	go discord.SendNewsEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf(
			"在 [%s](%s) 的直播间收到来自 [%s](%s) 的醒目留言 ",
			data.LiveInfo.Name,
			biliRoomLink(data.LiveInfo.RoomId),
			superchat.UserInfo.UName,
			biliSpaceLink(superchat.UID),
		),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "￥",
				Value: fmt.Sprintf("%v", superchat.Price),
			},
			{
				Name:  "内容",
				Value: superchat.Message,
			},
		},
	})

	return withBilibiliRisky(msg)
}

func init() {
	bilibili.MessageHandler.AddHandler(bilibili.SuperChatMessage, HandleSuperChatMsg)
}
