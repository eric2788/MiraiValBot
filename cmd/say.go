package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
	"strings"
)

func say(args []string, source *command.MessageSource) error {

	elements := source.Message.Elements

	textElements := qq.NewTextLn(strings.Join(args, " "))

	msg := message.NewSendingMessage()
	msg.Append(textElements)

	// find all possible elements to add
	for _, element := range elements {
		switch e := element.(type) {
		case *message.AtElement:
			msg.Append(e)
		case *message.FaceElement:
			msg.Append(e)
		case *message.GroupImageElement:
			msg.Append(e)
		}
	}

	return qq.SendGroupMessage(msg)
}

var sayCommand = command.NewNode([]string{"say", "说话", "说"}, "说话指令", false, say, "<讯息>")

func init() {
	command.AddCommand(sayCommand)
}
