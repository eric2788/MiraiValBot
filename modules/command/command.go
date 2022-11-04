package command

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Logiase/MiraiGo-Template/utils"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/eventhook"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/qq"
	"runtime/debug"
	"sync"
)

const Tag = "valbot.command"

type command struct {
}

func (c *command) HookEvent(bot *bot.Bot) {
	bot.GroupMessageEvent.Subscribe(func(ct *client.QQClient, msg *message.GroupMessage) {

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

		// avoid panic from InvokeCommand
		defer func() {

			if e := recover(); e != nil {
				err := fmt.Errorf(fmt.Sprintf("%v", e))
				logger.Errorf("處理指令 %s 時出現严重錯誤: %v", content, err)
				debug.PrintStack()
				_ = qq.SendGroupMessageByGroup(msg.GroupCode, qq.CreateReply(msg).Append(qq.NewTextf("处理此指令时出现严重错误: %v", err)))
			}

		}()

		admin := member.Permission <= client.Administrator
		response, err := InvokeCommand(content, admin, source)

		if err != nil {
			logger.Warnf("處理指令 %s 時出現錯誤: %v", content, err)
			errorMsg := qq.CreateReply(msg).Append(qq.NewTextf("处理此指令时出现错误: %v", err))
			_ = qq.SendGroupMessageByGroup(msg.GroupCode, errorMsg)
			return
		}

		logger.Debugf("%+v", *response)

		if response.Ignore {
			return
		} else if response.ShowHelp {

			if response.Content == "" {

				logger.Warnf("無法發送指令幫助，指令幫助為空。")

				hintMessage := qq.CreateReply(msg).Append(message.NewText("指令帮助为空"))
				_ = qq.SendGroupMessageByGroup(msg.GroupCode, hintMessage)

			} else {

				// 發送私人或临时会话訊息的指令幫助
				helpMessage := message.NewSendingMessage().Append(message.NewText(response.Content))

				_ = qq.SendWithRandomRiskyFunc(helpMessage, func() {

					if msg.Sender.IsFriend {
						if err = qq.SendPrivateMessage(msg.Sender.Uin, helpMessage); err == nil {
							logger.Infof("已向 %s (%d) 發送指令幫助的私人訊息。", msg.Sender.Nickname, msg.Sender.Uin)
						} else {
							logger.Warnf("無法向 %s (%d) 發送指令幫助的私人訊息: %v", msg.Sender.Nickname, msg.Sender.Uin, err)
						}
					} else {
						if err = qq.SendGroupTempMessage(msg.GroupCode, msg.Sender.Uin, helpMessage); err == nil {
							logger.Infof("已向 %s (%d) 發送指令幫助的臨時會話訊息。", msg.Sender.Nickname, msg.Sender.Uin)
						} else {
							logger.Warnf("無法向 %s (%d) 發送指令幫助的臨時會話訊息: %v", msg.Sender.Nickname, msg.Sender.Uin, err)
						}
					}

					hintMessage := qq.CreateReply(msg).Append(message.NewText("向群发送指令帮助讯息时由于遭到屡次风控，已改为私聊你指令列表。"))
					_ = qq.SendGroupMessageByGroup(msg.GroupCode, hintMessage)
				})

			}

		} else if response.Content != "" {
			m := qq.CreateReply(msg).Append(message.NewText(response.Content))
			_ = qq.SendGroupMessageByGroup(msg.GroupCode, m)
		}
	})

	// 瓦群成员自动接受好友邀请
	bot.NewFriendRequestEvent.Subscribe(func(ct *client.QQClient, req *client.NewFriendRequest) {
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
	defer wg.Done()
	logger.Info("指令管理模組已關閉")
}

func init() {
	bot.RegisterModule(instance)
	eventhook.HookLifeCycle(instance)
}
