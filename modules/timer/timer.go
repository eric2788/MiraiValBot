package timer

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
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
		wg:       &sync.WaitGroup{},
	}
	bgCtx = context.Background()
)

type Timer struct {
	timerMap map[string]*handler
	wg       *sync.WaitGroup
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
}

func (t *Timer) HookEvent(bot *bot.Bot) {
	for name := range t.timerMap {
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
	for name := range t.timerMap {
		go func(name string) {
			logger.Debugf("正在中止計時器任務 %s ...", name)
			_, err := t.StopTimer(name)
			if err != nil {
				logger.Warnf("中止計時器任務 %s 時出現錯誤: %v", name, err)
			}
		}(name)
	}
	logger.Debugf("正在等待所有計時器任務關閉...")
	t.wg.Wait()
	logger.Info("定時器任務模組已關閉")
}

func (t *Timer) StopTimer(name string) (bool, error) {
	if timer, ok := t.timerMap[name]; !ok {
		return false, fmt.Errorf("找不到此計時器任務")
	} else {

		if !timer.Started {
			logger.Debugf("計時器任務 %s 未開始，停止失敗。", name)
			return false, nil
		}

		timer.canceller()
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
		t.wg.Add(1)
		logger.Debugf("wg added 1")
		go startTimer(name, ctx, time.NewTicker(timer.duration), bot, timer.job, func() {
			t.wg.Done()
			logger.Debugf("wg removed 1")
			timer.Started = false
		})
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
	eventhook.HookLifeCycle(instance)
}
