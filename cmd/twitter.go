package cmd

import (
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/MiraiValBot/sites/twitter"
)

func tListen(args []string, source *command.MessageSource) error {
	screenId := args[0]
	reply := qq2.CreateReply(source.Message)

	result, err := twitter.StartListen(screenId)

	if err != nil {
		reply.Append(qq2.NewTextf("启动监听时出现错误: %v", err))
	} else if result {
		reply.Append(qq2.NewTextf("开始监听推特用户 %s", screenId))
	} else {
		reply.Append(qq2.NewTextf("该用户 %s 已启动监听", screenId))
	}

	return qq2.SendGroupMessage(reply)

}

func tshowReply(args []string, source *command.MessageSource) error {

	file.UpdateStorage(func() {
		file.DataStorage.Twitter.ShowReply = !file.DataStorage.Twitter.ShowReply
	})

	reply := qq2.CreateReply(source.Message)
	reply.Append(qq2.NewTextf("已成功设置推送推文回复为: %v", file.DataStorage.Twitter.ShowReply))

	return qq2.SendGroupMessage(reply)
}

func tTerminate(args []string, source *command.MessageSource) error {
	screenId := args[0]
	reply := qq2.CreateReply(source.Message)

	result, err := twitter.StopListen(screenId)

	if err != nil {
		reply.Append(qq2.NewTextf("中止监听时出现错误: %v", err))
	} else if result {
		reply.Append(qq2.NewTextf("已中止监听推特用户 %s", screenId))
	} else {
		reply.Append(qq2.NewTextf("你尚未开始监听此推特用户。"))
	}

	return qq2.SendGroupMessage(reply)
}

func tListening(args []string, source *command.MessageSource) error {
	listening := file.DataStorage.Listening.Twitter
	reply := qq2.CreateReply(source.Message)

	if listening.Size() > 0 {
		reply.Append(qq2.NewTextf("正在监听的推特用户: %v", strings.Join(listening.ToArr(), ", ")))
	} else {
		reply.Append(message.NewText("没有在监听的推特用户"))
	}

	return qq2.SendWithRandomRiskyStrategy(reply)
}

var (
	tListenCommand    = command.NewNode([]string{"listen", "监听"}, "启动监听用户", true, tListen, "<用户ID>")
	tTerminateCommand = command.NewNode([]string{"terminate", "中止", "中止监听"}, "中止监听用户", true, tTerminate, "<用户ID>")
	tListeningCommand = command.NewNode([]string{"listening", "正在监听", "监听列表"}, "获取监听列表", false, tListening)
	tShowReplyCommand = command.NewNode([]string{"showReply", "推送回复", "推送推文回复"}, "切换推文回复推送", true, tshowReply)
)

var twitterCommand = command.NewParent([]string{"twitter", "推特"}, "推特指令",
	tListenCommand,
	tTerminateCommand,
	tListeningCommand,
	tShowReplyCommand,
)

func init() {
	command.AddCommand(twitterCommand)
}
