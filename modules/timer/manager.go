package timer

import (
	"context"
	"github.com/Logiase/MiraiGo-Template/bot"
	"time"
)

func startTimer(name string, ctx context.Context, ticker *time.Ticker, bot *bot.Bot, handle Job) {
	logger.Infof("計時器任務 %s 開始。", name)
	for {
		select {
		case <-ctx.Done():
			ticker.Stop()
			logger.Infof("計時器任務 %s 已停止。", name)
			return
		case <-ticker.C:
			err := handle(bot)
			if recoverError, ok := recover().(error); ok {
				err = recoverError
			}
			if err != nil {
				logger.Errorf("執行計時器任務 %s 時出現錯誤: %v", name, err)
			}
		}
	}
}
