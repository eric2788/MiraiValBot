package chat_reply

import (
	"sync"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/eventhook"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/common-utils/array"
)

const Tag = "valbot.chat_reply"

var (
	logger   = utils.GetModuleLogger(Tag)
	instance = &atResponse{
		strategies: []ResponseStrategy{
			&aiChatResponse{},
			&randomResponse{},
		},
	}
)

type (
	atResponse struct {
		strategies []ResponseStrategy
	}

	ResponseStrategy interface {
		Response(msg *message.GroupMessage) (*message.SendingMessage, error)
	}
)

func (a *atResponse) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (a *atResponse) Init() {
}

func (a *atResponse) PostInit() {
}

func (a *atResponse) Serve(bot *bot.Bot) {
}

func (a *atResponse) Start(bot *bot.Bot) {
	logger.Infof("聊天回复模组已启动。")
}

func (a *atResponse) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	defer wg.Done()
	logger.Infof("聊天回复模组已关闭。")
}

func (a *atResponse) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(cl *client.QQClient, msg *message.GroupMessage) {
		content := qq.ParseMsgContent(msg.Elements)

		if array.IndexOfInt64(content.At, cl.Uin) != -1 && len(content.Texts) > 0 {

			for _, strategy := range a.strategies {
				send, err := strategy.Response(msg)
				if err == nil {
					_ = qq.SendGroupMessageByGroup(msg.GroupCode, send)
					break
				}
			}

		}
	})
}

func init() {
	bot.RegisterModule(instance)
	eventhook.HookLifeCycle(instance)
}
