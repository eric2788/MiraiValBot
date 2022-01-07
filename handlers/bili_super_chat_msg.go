package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
)

func HandleSuperChatMsg(bot *bot.Bot, data *bilibili.LiveData) error {

	superchat := data.Content["data"].(*bilibili.SuperChatMessageData)

	fmt.Printf("從房間 %d 收到來自 %s 的 ￥%d 醒目留言: %s\n", data.LiveInfo.RoomId, superchat.UserInfo.UName, superchat.Price, superchat.Message)
	return nil
}

func init() {
	bilibili.RegisterDataHandler(bilibili.SuperChatMessage, HandleSuperChatMsg)
}
