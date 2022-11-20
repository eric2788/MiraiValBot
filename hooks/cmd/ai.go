package cmd

import (
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/huggingface"
)

// putting hugging face api here...

func generateHuggingFaceImage(model string, args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在生成图像...."))
	_ = qq.SendGroupMessage(reply)

	inputs := strings.Join(args, " ")

	b, err := huggingface.GetResultImage(model,
		huggingface.NewParam(
			huggingface.Input(inputs),
		),
	)

	if err != nil {
		return err
	}

	img, err := qq.NewImageByByte(b)

	if err != nil {
		return err
	}

	msg := message.NewSendingMessage()
	msg.Append(img)
	return qq.SendGroupMessage(msg)
}

func aiWaifu(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("Nilaier/Waifu-Diffusers", args, source)
}

func aiWaifu2(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("hakurei/waifu-diffusion", args, source)
}

func aiMadoka(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("yuk/madoka-waifu-diffusion", args, source)
}

var (
	aiWaifuCommand  = command.NewNode([]string{"waifu"}, "文字生成图像(waifu)", false, aiWaifu, "<文字>")
	aiWaifu2Command = command.NewNode([]string{"waifu2"}, "文字生成图像(waifu2)", false, aiWaifu2, "<文字>")
	aiMadokaCommand = command.NewNode([]string{"madoka", "円香", "画円香"}, "文字生成图像(円香)", false, aiMadoka, "<文字>")
)

var aiCommand = command.NewParent([]string{"ai", "人工智能"}, "AI相关指令",
	aiWaifuCommand,
	aiWaifu2Command,
	aiMadokaCommand,
)

func init() {
	command.AddCommand(aiCommand)
}
