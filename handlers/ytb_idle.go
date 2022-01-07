package handlers

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/eric2788/MiraiValBot/sites/youtube"
)

func HandleIdle(bot *bot.Bot, info *youtube.LiveInfo) error {
	return nil
}

func init() {
	youtube.RegisterDataHandler(youtube.Idle, HandleIdle)
}
