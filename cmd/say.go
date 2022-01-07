package cmd

import (
	"fmt"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	"strings"
)

func Say(args []string, source *command.MessageSource) error {
	if source == nil { // 测试用
		fmt.Printf("机器人说话: %s\n", strings.Join(args, " "))
	} else {
		source.Client.SendGroupMessage(file.ApplicationYaml.Val.GroupId, message.NewSendingMessage().Append(message.NewText(strings.Join(args, " "))))
	}
	return nil
}

var say = command.NewNode([]string{"say", "说话", "说"}, "说话指令", false, Say, "<话>")

func init() {
	command.AddCommand(say)
}
