package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	qq2 "github.com/eric2788/MiraiValBot/qq"
)

func checkRes(args []string, source *command.MessageSource) error {
	reply := qq2.CreateReply(source.Message)
	reply.Append(message.NewText("回应列表: ")).Append(message.NewText("\n"))
	for content, res := range file.DataStorage.Responses {
		reply.Append(qq2.NewTextf("%s: %s\n", content, res))
	}

	return qq2.SendGroupMessage(reply)
}

func setRes(args []string, source *command.MessageSource) error {

	content, res := args[0], args[1]

	file.UpdateStorage(func() {
		file.DataStorage.Responses[content] = res
	})

	reply := qq2.CreateReply(source.Message).Append(qq2.NewTextf("已成功设置 %s 的回应为 %s。", content, res))
	return qq2.SendGroupMessage(reply)
}

func removeRes(args []string, source *command.MessageSource) error {
	content := args[0]

	reply := qq2.CreateReply(source.Message)

	if _, ok := file.DataStorage.Responses[content]; !ok {
		reply.Append(message.NewText("找不到这个文字。"))
	} else {
		file.UpdateStorage(func() {
			delete(file.DataStorage.Responses, content)
		})

		reply.Append(message.NewText("移除成功。"))
	}

	return qq2.SendGroupMessage(reply)
}

var (
	checkResCommand  = command.NewNode([]string{"check", "检查"}, "检查所有自定义回应", false, checkRes)
	setResCommand    = command.NewNode([]string{"set", "设置"}, "设置自定义回应", false, setRes, "<文字>", "<回应>")
	removeResCommand = command.NewNode([]string{"remove", "移除"}, "移除回应", false, removeRes, "<文字>")
)

var resCommand = command.NewParent([]string{"res"}, "自定义回应",
	checkResCommand,
	setResCommand,
	removeResCommand,
)

func init() {
	command.AddCommand(resCommand)
}
