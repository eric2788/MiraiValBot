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

func HandleEnterRoom(bot *bot.Bot, data *bilibili.LiveData) error {
	entered := data.Content["data"].(map[string]interface{})
	uname := entered["uname"].(string)
	uid := int64(entered["uid"].(float64))

	if !bilibili.IsHighlighter(uid) {
		return nil
	}

	logger.Debugf("%s 進入了 %s 的直播間 (%d)\n", uname, data.LiveInfo.Name, data.LiveInfo.RoomId)

	discordMessage := &discordgo.MessageEmbed{
		Description: fmt.Sprintf("噔噔咚！你所关注的用户 [%s](%s) 进入了 [%s](%s) 的直播间。", uname, biliSpaceLink(uid), data.LiveInfo.Name, biliRoomLink(data.LiveInfo.RoomId)),
	}

	go discord.SendNewsEmbed(discordMessage)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("噔噔咚！你所关注的用户 %s 进入了 %s 的直播间。", uname, data.LiveInfo.Name))

	return withBilibiliRisky(msg)
}

func init() {
	bilibili.MessageHandler.AddHandler(bilibili.InteractWord, HandleEnterRoom)
}
