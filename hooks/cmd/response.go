package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/file"
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"strings"
)

func checkRes(args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)
	reply.Append(message.NewText("回应列表: ")).Append(message.NewText("\n"))
	for content, res := range file.DataStorage.Responses {
		reply.Append(qq.NewTextf("%s: %s\n", content, res))
	}

	return qq.SendGroupMessage(reply)
}

func setRes(args []string, source *command.MessageSource) error {

	content, res := args[0], args[1]

	logger.Debugf("content: %s, res: %s, args: %s", content, res, strings.Join(args, ", "))

	if content == "" {
		return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("触发参数不能为空")))
	} else if res == "" {
		return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText("回应参数不能为空")))
	}

	file.UpdateStorage(func() {
		file.DataStorage.Responses[content] = res
	})

	reply := qq.CreateReply(source.Message).Append(qq.NewTextf("已成功设置 %s 的回应为 %s。", content, res))
	return qq.SendGroupMessage(reply)
}

func removeRes(args []string, source *command.MessageSource) error {
	content := args[0]

	reply := qq.CreateReply(source.Message)

	if _, ok := file.DataStorage.Responses[content]; !ok {
		reply.Append(message.NewText("找不到这个文字。"))
	} else {
		file.UpdateStorage(func() {
			delete(file.DataStorage.Responses, content)
		})

		reply.Append(message.NewText("移除成功。"))
	}

	return qq.SendGroupMessage(reply)
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
