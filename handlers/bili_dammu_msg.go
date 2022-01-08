package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
)

func HandleDanmuMsg(bot *bot.Bot, data *bilibili.LiveData) error {

	room := data.LiveInfo.RoomId

	info := data.Content["info"].([]interface{})
	userInfo := info[2].([]interface{})

	danmu := info[1].(string)
	uname := userInfo[1].(string)
	uid := int64(userInfo[0].(float64))

	// debug only
	logger.Infof("從房間 %d 收到來自 %s (%d) 的彈幕: %s\n", room, uname, uid, danmu)

	return nil
}

func init() {
	bilibili.RegisterDataHandler(bilibili.DanmuMsg, HandleDanmuMsg)
}
