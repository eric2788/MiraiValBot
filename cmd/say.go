package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	"strings"
)

func say(args []string, source *command.MessageSource) error {
	source.Client.SendGroupMessage(file.ApplicationYaml.Val.GroupId, message.NewSendingMessage().Append(message.NewText(strings.Join(args, " "))))
	return nil
}

var sayCommand = command.NewNode([]string{"say", "说话", "说"}, "说话指令", false, say, "<讯息>")

func init() {
	command.AddCommand(sayCommand)
}
