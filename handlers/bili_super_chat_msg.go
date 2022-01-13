package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
	"github.com/eric2788/MiraiValBot/utils/qq"
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

	return withRisky(msg)
}

func init() {
	bilibili.RegisterDataHandler(bilibili.SuperChatMessage, HandleSuperChatMsg)
}
