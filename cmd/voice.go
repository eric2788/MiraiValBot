package cmd

import (
	"errors"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"strings"
)

func voice(args []string, source *command.MessageSource) error {

	content := strings.Join(args, " ")

	var err error

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic when converting %v to tts: %v", content, r)
			err = errors.New("tts 转换失败")
		}
	}()

	voiceElement, err := qq2.NewTts(content)

	if err != nil {
		logger.Errorf("tts 转换失败: %v", err)
		return err
	}

	logger.Infof("嘗試發送voiceElement: %v", content)

	return qq2.SendGroupMessage(message.NewSendingMessage().Append(voiceElement))
}

var voiceCommand = command.NewNode([]string{"voice", "speak", "语音"}, "语音指令", false, voice, "<讯息>")

func init() {
	command.AddCommand(voiceCommand)
}
