package verbose

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/eventhook"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"sync"
)

type verbose struct {
}

const Tag = "valbot.verbose"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &verbose{}
)

func (v *verbose) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (v *verbose) Init() {
}

func (v *verbose) PostInit() {
}

func (v *verbose) Serve(bot *bot.Bot) {
}

func (v *verbose) Start(bot *bot.Bot) {
	logger.Infof("Verbose 模组已启动")
}

func (v *verbose) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	logger.Infof("verbose 模组已关闭")
}

func (v *verbose) HookEvent(qqBot *bot.Bot) {

	verboseLiveRoomStatus()

	qqBot.OnGroupMessageRecalled(func(c *client.QQClient, event *client.GroupMessageRecalledEvent) {

		if !file.DataStorage.Setting.VerboseDelete {
			return
		}

		msg := message.NewSendingMessage()
		msg.Append(qq.NewTextfLn("%s 所撤回的消息: "))
		m, err := qq.GetGroupMessage(event.GroupCode, int64(event.MessageId))
		if err != nil {
			msg.Append(qq.NewTextf("获取消息失败: %v", err))
		} else {
			for _, element := range m.Elements {
				msg.Append(element)
			}
		}
		_ = qq.SendGroupMessage(msg)
	})

}

func init() {
	bot.RegisterModule(instance)
	eventhook.HookLifeCycle(instance)
}
