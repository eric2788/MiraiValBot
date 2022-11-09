package cmd

import (
	"fmt"
	"net/url"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/common-utils/request"
)

func testSendMp3Voice(args []string, source *command.MessageSource) error {
	data, err := request.GetBytesByUrl(fmt.Sprintf("https://genshin.azurewebsites.net/api/speak?format=mp3&text=%s&id=0", url.QueryEscape("测试测试")))
	if err != nil {
		return err
	}
	voice, err := qq.NewVoiceByBytes(data)
	if err != nil {
		return err
	}
	return qq.SendGroupMessage(message.NewSendingMessage().Append(voice))
}

func testSendWavVoice(args []string, source *command.MessageSource) error {
	data, err := request.GetBytesByUrl(fmt.Sprintf("https://genshin.azurewebsites.net/api/speak?format=wav&text=%s&id=0", url.QueryEscape("测试测试")))
	if err != nil {
		return err
	}
	voice, err := qq.NewVoiceByBytes(data)
	if err != nil {
		return err
	}
	return qq.SendGroupMessage(message.NewSendingMessage().Append(voice))
}

var testCommands = []command.CmdHandler{
	testSendMp3Voice,
	testSendWavVoice,
}

func init() {
	nodes := make([]command.Node, len(testCommands))
	for i, handler := range testCommands {
		name := fmt.Sprintf("%d", i+1)
		nodes[i] = command.NewNode([]string{name}, name, true, handler)
	}
	var testCommand = command.NewParent([]string{"test", "测试"}, "测试指令", nodes...)
	testCommand.AdminOnly = true
	command.AddCommand(testCommand)
}
