package cmd

import (
	"errors"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/aivoice"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/qq"
)

func voiceQQ(args []string, source *command.MessageSource) error {

	content := strings.Join(args, " ")

	var err error

	defer func() {
		if r := recover(); r != nil {
			logger.Errorf("panic when converting %v to tts: %v", content, r)
			err = errors.New("tts 转换失败")
		}
	}()

	voiceElement, err := qq.NewTts(content)

	if err != nil {
		logger.Errorf("tts 转换失败: %v", err)
		return err
	}

	logger.Infof("嘗試發送voiceElement: %v", content)

	return qq.SendGroupMessage(message.NewSendingMessage().Append(voiceElement))
}

func voiceGenshin(args []string, source *command.MessageSource) error {
	actor, content := args[0], strings.Join(args[1:], "，")

	reply := qq.CreateReply(source.Message).Append(qq.NewTextf("正在尝试生成 %s 的语音...", actor))
	_ = qq.SendGroupMessage(reply)

	data, err := aivoice.GetGenshinVoice(content, actor)
	if err != nil {
		return err
	}

	/*
		voice, err := qq.NewVoiceByBytes(data)
		if err != nil {
			return err
		}

	*/

	voice := &message.GroupVoiceElement{Data: data}
	return qq.SendGroupMessage(message.NewSendingMessage().Append(voice))
}

var (
	voiceQQCommand      = command.NewNode([]string{"qq", "腾讯"}, "腾讯QQ的语音", false, voiceQQ, "<讯息>")
	voiceGenshinCommand = command.NewNode([]string{"genshin", "原神", "ys"}, "原神角色语音", false, voiceGenshin, "<角色>", "<讯息>")
)

var voiceCommand = command.NewParent(
	[]string{"voice", "speak", "语音"},
	"语音指令",
	voiceQQCommand,
	voiceGenshinCommand,
)

func init() {
	command.AddCommand(voiceCommand)
}
