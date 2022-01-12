package timer

import (
	"context"
	"github.com/Logiase/MiraiGo-Template/bot"
	"sync"
	"time"
)

func startTimer(name string, ctx context.Context, ticker *time.Ticker, bot *bot.Bot, handle Job, wg *sync.WaitGroup) {
	logger.Infof("計時器任務 %s 開始。", name)
	defer wg.Done()
	defer ticker.Stop()
	run(name, bot, handle) // 開頭運行一次
	for {
		select {
		case <-ctx.Done():
			logger.Infof("計時器任務 %s 已停止。", name)
			return
		case <-ticker.C:
			run(name, bot, handle)
		}
	}
}

func run(name string, bot *bot.Bot, handle Job) {
	err := handle(bot)
	if recoverError, ok := recover().(error); ok {
		err = recoverError
	}
	if err != nil {
		logger.Errorf("執行計時器任務 %s 時出現錯誤: %v", name, err)
	}
}
