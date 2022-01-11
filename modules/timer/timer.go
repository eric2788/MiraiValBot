package timer

import (
	"context"
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"sync"
	"time"
)

const Tag = "valbot.timer"

type Job = func(bot *bot.Bot) error

type handler struct {
	job       Job
	duration  time.Duration
	ctx       context.Context
	canceller context.CancelFunc
	Started   bool
}

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &Timer{
		timerMap: make(map[string]*handler),
	}
	bgCtx = context.Background()
)

type Timer struct {
	timerMap map[string]*handler
}

func (t *Timer) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (t *Timer) Init() {
}

func (t *Timer) PostInit() {
}

func (t *Timer) Serve(bot *bot.Bot) {
	for name, _ := range t.timerMap {
		_, err := t.StartTimer(name, bot)
		if err != nil {
			logger.Warnf("啟動計時器任務 %s 時出現錯誤: %v", name, err)
		}
	}
}

func (t *Timer) Start(bot *bot.Bot) {
	logger.Info("定時器任務模組已啟動")
}

func (t *Timer) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	for name, timer := range t.timerMap {
		_, err := t.StopTimer(name)
		if err != nil {
			logger.Warnf("中止計時器任務 %s 時出現錯誤: %v", name, err)
			continue
		}
		<-timer.ctx.Done()
	}
	logger.Info("定時器任務模組已關閉")
}

func (t *Timer) StopTimer(name string) (bool, error) {
	if timer, ok := t.timerMap[name]; !ok {
		return false, fmt.Errorf("找不到此計時器任務")
	} else {

		if !timer.Started {
			return false, nil
		}

		timer.canceller()
		<-timer.ctx.Done()
		timer.Started = false
		return true, nil
	}
}

func (t *Timer) StartTimer(name string, bot *bot.Bot) (bool, error) {
	if timer, ok := t.timerMap[name]; !ok {
		return false, fmt.Errorf("找不到此計時器任務")
	} else {

		if timer.Started {
			return false, nil
		}

		// new context
		ctx, cancel := context.WithCancel(bgCtx)
		timer.ctx = ctx
		timer.canceller = cancel

		ticker := time.NewTicker(timer.duration)
		go startTimer(name, ctx, ticker, bot, timer.job)
		timer.Started = true
		return true, nil
	}
}

func RegisterTimer(name string, duration time.Duration, handle Job) {
	if _, ok := instance.timerMap[name]; ok {
		logger.Warnf("定時器任務 %s 已存在，已略過。", name)
		return
	}
	instance.timerMap[name] = &handler{
		job:      handle,
		duration: duration,
	}
	logger.Infof("已成功註冊定時器任務 %s", name)
}

func init() {
	bot.RegisterModule(instance)
}
