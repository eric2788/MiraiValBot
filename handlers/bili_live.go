package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/sites/bilibili"
)

func HandleLive(bot *bot.Bot, data *bilibili.LiveData) error {
	return nil
}

func init() {
	bilibili.RegisterDataHandler(bilibili.Live, HandleLive)
}
