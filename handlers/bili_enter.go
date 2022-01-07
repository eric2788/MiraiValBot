package handlers

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
)

func HandleEnterRoom(bot *bot.Bot, data *bilibili.LiveData) error {
	entered := data.Content["data"].(map[string]interface{})
	uname := entered["uname"].(string)
	//uid := int64(entered["uid"].(float64))
	fmt.Printf("%s 進入了 %s 的直播間 (%d)\n", uname, data.LiveInfo.Name, data.LiveInfo.RoomId)
	return nil
}

func init() {
	bilibili.RegisterDataHandler(bilibili.InteractWord, HandleEnterRoom)
}
