package cmd

import (
	"crypto/md5"
	"errors"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"strings"
)

func voice(args []string, source *command.MessageSource) error {

	content := strings.Join(args, " ")

	var err error

	defer func() {
		if r := recover(); r != nil {
			err = errors.New("tts 转换失败")
		}
	}()

	data, err := source.Client.GetTts(content)

	if err != nil {
		return err
	}

	voice := &message.VoiceElement{
		Name: content,
		Data: data,
		Size: int32(len(data)),
		Md5:  md5.New().Sum(data),
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, message.NewSendingMessage().Append(voice))

	return err
}

var voiceCommand = command.NewNode([]string{"voice", "speak", "语音"}, "语音指令", false, voice, "<讯息>")

func init() {
	command.AddCommand(voiceCommand)
}
