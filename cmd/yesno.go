package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"regexp"
)

var yesNoPattern = regexp.MustCompile("^.+是.+吗[\\?？]*$")

func setYesNo(args []string, source *command.MessageSource) error {
	question, ans := args[0], args[1] == "true"

	reply := qq.CreateReply(source.Message)

	if !yesNoPattern.MatchString(question) {
		reply.Append(message.NewText("不是一个有效的问题"))
	} else {
		file.UpdateStorage(func() {
			file.DataStorage.Answers[question] = ans
		})
		reply.Append(qq.NewTextf("已成功设置 %s 的答案为 %v", question, ans))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

func removeYesNo(args []string, source *command.MessageSource) error {
	question := args[0]

	reply := qq.CreateReply(source.Message)

	if _, ok := file.DataStorage.Answers[question]; !ok {
		reply.Append(message.NewText("找不到此问题"))
	} else {
		file.UpdateStorage(func() {
			delete(file.DataStorage.Answers, question)
		})
		reply.Append(qq.NewTextf("已成功移除 %s 的答案", question))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

func checkYesNo(args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)
	reply.Append(message.NewText("问题列表: "))

	for question, ans := range file.DataStorage.Answers {
		reply.Append(qq.NewTextf("%s = %v\n", question, ans))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

var (
	setYesNoCommand    = command.NewNode([]string{"set"}, "设置yes no答案", false, setYesNo, "<问题>", "<true | false>")
	removeYesNoCommand = command.NewNode([]string{"remove"}, "移除问题", false, removeYesNo, "<问题>")
	checkYesNoCommand  = command.NewNode([]string{"check"}, "移除问题", false, checkYesNo)
)

var yesNoCommand = command.NewParent([]string{"yesno"}, "YesNo指令",
	setYesNoCommand,
	removeYesNoCommand,
	checkYesNoCommand,
)

func init() {
	command.AddCommand(yesNoCommand)
}
