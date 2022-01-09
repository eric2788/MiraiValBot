package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
	"github.com/eric2788/MiraiValBot/utils/qq"
)

func HandleEnterRoom(bot *bot.Bot, data *bilibili.LiveData) error {
	entered := data.Content["data"].(map[string]interface{})
	uname := entered["uname"].(string)
	uid := int64(entered["uid"].(float64))

	if !bilibili.IsHighlighter(uid) {
		return nil
	}

	logger.Debugf("%s 進入了 %s 的直播間 (%d)\n", uname, data.LiveInfo.Name, data.LiveInfo.RoomId)

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("噔噔咚！你所关注的用户 %s 进入了 %s 的直播间。", uname, data.LiveInfo.Name))

	bot.SendGroupMessage(qq.ValGroupInfo.Uin, msg)
	return nil
}

func init() {
	bilibili.RegisterDataHandler(bilibili.InteractWord, HandleEnterRoom)
}
