package counting

import (
	"strings"
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"golang.org/x/exp/maps"
)

const Tag = "valbot.counting"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &wordCounting{}
)

type wordCounting struct {
}

func (w *wordCounting) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (w *wordCounting) Init() {
}

func (w *wordCounting) PostInit() {
}

func (w *wordCounting) Serve(bot *bot.Bot) {
}

func (w *wordCounting) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(client *client.QQClient, event *message.GroupMessage) {

		// 非瓦群无视
		if event.GroupCode != qq.ValGroupInfo.Code {
			return
		}

		// 机器人无视
		if event.Sender.Uin == client.Uin {
			return
		}

		msg := strings.TrimSpace(command.ExtractPrefix(strings.Join(qq.ParseMsgContent(event.Elements).Texts, " ")))
		counts, ok := file.DataStorage.WordCounts[msg]

		// 没有这个字词记录
		if !ok {
			return
		}

		times, ok := counts[event.Sender.Uin]

		if !ok {
			times = 0
		}

		times++

		file.UpdateStorage(func() {
			file.DataStorage.WordCounts[msg][event.Sender.Uin] = times
		})

		logger.Infof("%s 字词已在本群说了合共 %d 次。", msg, w.sum(msg))

		if times%int64(file.DataStorage.Setting.TimesPerNotify) != 0 {
			return
		}

		send := message.NewSendingMessage()
		send.Append(qq.NewTextf("%s 已经在本群说了 %q 合共 %d 次", event.Sender.DisplayName(), msg, times))
		_ = qq.SendGroupMessage(send)
	})
}

func (w *wordCounting) sum(msg string) int64 {
	sum := int64(0)
	for _, t := range maps.Values(file.DataStorage.WordCounts[msg]) {
		sum += t
	}
	return sum
}

func (w *wordCounting) Start(bot *bot.Bot) {
	logger.Infof("字词计算模组已启动。")
}

func (w *wordCounting) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("字词计算模组已关闭。")
}

func init() {
	bot.RegisterModule(instance)
	eventhook.HookLifeCycle(instance)
}
