package eventhook

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
)

// 當會話恢復失敗需要重新登錄的時候，QQClient所有事件註冊會被清空 (issue?)
// 因此我需要把所有 module 的事件註冊轉移到 登入 後

var logger = utils.GetModuleLogger("valbot.eventhook")

type EventHooker interface {
	HookEvent(bot *bot.Bot)
}

var hookers []EventHooker

func HookLifeCycle(hooker EventHooker) {
	hookers = append(hookers, hooker)
}

// HookBotEvents 必須在 Login 後呼叫
func HookBotEvents() {
	for _, hook := range hookers {
		hook.HookEvent(bot.Instance)
	}
	logger.Info("已成功掛接所有監聽事件。")
}
