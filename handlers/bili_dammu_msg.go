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

func HandleDanmuMsg(bot *bot.Bot, data *bilibili.LiveData) error {

	room := data.LiveInfo.RoomId

	info := data.Content["info"].([]interface{})
	userInfo := info[2].([]interface{})

	danmu := info[1].(string)
	uname := userInfo[1].(string)
	uid := int64(userInfo[0].(float64))

	if !bilibili.IsHighlighter(uid) {
		return nil
	}

	//debug only
	logger.Debugf("從房間 %d 收到來自 %s (%d) 的彈幕: %s\n", room, uname, uid, danmu)

	go discord.SendNewsEmbed(&discordgo.MessageEmbed{
		Description: fmt.Sprintf("[%s](%s) 在 [%s](%s) 的直播间发送了一则讯息: ", uname, biliSpaceLink(uid), biliRoomLink(room), data.LiveInfo.Name),
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  "弹幕",
				Value: danmu,
			},
		},
	})

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextfLn("%s 在 %s 的直播间发送了一则消息", uname, data.LiveInfo.Name))
	msg.Append(qq.NewTextfLn("弹幕: %s", danmu))

	return withBilibiliRisky(msg)
}

func biliSpaceLink(uid int64) string {
	return fmt.Sprintf("https://space.bilibili.com/%d", uid)
}

func biliRoomLink(room int64) string {
	return fmt.Sprintf("https://live.bilibili.com/%d", room)
}

func init() {
	bilibili.RegisterDataHandler(bilibili.DanmuMsg, HandleDanmuMsg)
}
