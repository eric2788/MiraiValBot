package youtube

import (
	"github.com/Logiase/MiraiGo-Template/bot"
)

type LiveDataHandler = func(*bot.Bot, *LiveInfo) error

type Handle = func(bot *bot.Bot, info *LiveInfo)

func RegisterDataHandler(cmd string, handle LiveDataHandler) {
	MessageHandler.AddHandler(cmd, handle)
}
