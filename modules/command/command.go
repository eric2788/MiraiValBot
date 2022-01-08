package command

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/eventhook"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"sync"
)

const Tag = "valbot.command"

type command struct {
}

func (c *command) HookEvent(bot *bot.Bot) {
	bot.OnGroupMessage(func(ct *client.QQClient, msg *message.GroupMessage) {

		// 非瓦群無視
		if msg.GroupCode != file.ApplicationYaml.Val.GroupId {
			logger.Infof("非瓦群，已略過。")
			return
		}
		source := &MessageSource{ct, msg}
		content := msg.ToString()
		member := qq.FindGroupMember(msg.Sender.Uin)
		if member == nil {
			logger.Infof("%s (%d) 不是瓦群成員，已略過。", msg.Sender.Nickname, msg.Sender.Uin)
			return
		}

		admin := member.Permission >= client.Administrator
		response, err := InvokeCommand(content, admin, source)

		if e, ok := recover().(error); ok {
			err = e
		}

		if err != nil {
			logger.Warnf("處理指令 %s 時出現錯誤: %v", content, err)
			errorMsg := message.NewSendingMessage().Append(message.NewText(fmt.Sprintf("处理此指令时出现错误: %v", err))).Append(message.NewReply(msg))
			ct.SendGroupMessage(msg.GroupCode, errorMsg)
			return
		}

		logger.Debugf("%+v", *response)

		if response.Ignore {
			return
		} else if response.ShowHelp {

			var responseContent string

			if response.Content == "" {

				logger.Warnf("無法發送指令幫助，指令幫助為空。")
				responseContent = "指令帮助为空"

			} else {

				// 發送私人或临时会话訊息的指令幫助
				helpMessage := message.NewSendingMessage().Append(message.NewText(response.Content))

				if msg.Sender.IsFriend {
					ct.SendPrivateMessage(msg.Sender.Uin, helpMessage)
					logger.Infof("已向 %s (%d) 發送指令幫助的私人訊息。", msg.Sender.Nickname, msg.Sender.Uin)
				} else {
					ct.SendGroupTempMessage(msg.GroupCode, msg.Sender.Uin, helpMessage)
					logger.Infof("已向 %s (%d) 發送指令幫助的臨時會話訊息。", msg.Sender.Nickname, msg.Sender.Uin)
				}
				responseContent = "未知参数，已私聊你指令列表。"

			}

			// 發送群組訊息提示
			hintMessage := message.NewSendingMessage().Append(message.NewReply(msg)).Append(message.NewText(responseContent))
			ct.SendGroupMessage(msg.GroupCode, hintMessage)
		} else if response.Content != "" {
			m := message.NewSendingMessage().Append(message.NewReply(msg)).Append(message.NewText(response.Content))
			ct.SendGroupMessage(msg.GroupCode, m)
		}
	})

	// 瓦群成员自动接受好友邀请
	bot.OnNewFriendRequest(func(ct *client.QQClient, req *client.NewFriendRequest) {
		if qq.FindGroupMember(req.RequesterUin) == nil {
			logger.Infof("%s (%d) 非瓦群成員，已無視好友邀請。", req.RequesterNick, req.RequesterUin)
			req.Reject()
			return
		}
		req.Accept()
	})
}

var (
	instance = &command{}
	logger   = utils.GetModuleLogger(Tag)
)

func (c *command) MiraiGoModule() bot.ModuleInfo {
	return bot.ModuleInfo{
		ID:       Tag,
		Instance: instance,
	}
}

func (c *command) Init() {
}

func (c *command) PostInit() {
}

func (c *command) Serve(bot *bot.Bot) {
}

func (c *command) Start(bot *bot.Bot) {
	logger.Info("指令管理模組已啟動。")
}

func (c *command) Stop(bot *bot.Bot, wg *sync.WaitGroup) {
	logger.Info("指令管理模組已關閉")
	wg.Done()
}

func init() {
	bot.RegisterModule(instance)
	eventhook.HookLifeCycle(instance)
}
