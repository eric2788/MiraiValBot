package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"strings"
)

func say(args []string, source *command.MessageSource) error {
	return qq.SendGroupMessage(message.NewSendingMessage().Append(message.NewText(strings.Join(args, " "))))
}

var sayCommand = command.NewNode([]string{"say", "说话", "说"}, "说话指令", false, say, "<讯息>")

func init() {
	command.AddCommand(sayCommand)
}
