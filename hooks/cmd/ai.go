package cmd

import (
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/huggingface"
)

func aiWaifu(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("Nilaier/Waifu-Diffusers", args, source)
}

func aiWaifu2(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("hakurei/waifu-diffusion", args, source)
}

func aiPaint(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("runwayml/stable-diffusion-v1-5", args, source)
}

func aiMadoka(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("yuk/madoka-waifu-diffusion", args, source)
}

func aiPrompt(args []string, source *command.MessageSource) error {
	return generateHuggingFaceText("Gustavosta/MagicPrompt-Stable-Diffusion", args, source)
}

func aiChinesePaint(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage("IDEA-CCNL/Taiyi-Stable-Diffusion-1B-Chinese-v0.1", args, source)
}

var (
	aiWaifuCommand   = command.NewNode([]string{"waifu"}, "文字生成图像(waifu)", false, aiWaifu, "<文字>")
	aiWaifu2Command  = command.NewNode([]string{"waifu2"}, "文字生成图像(waifu2)", false, aiWaifu2, "<文字>")
	aiPaintCNCommand = command.NewNode([]string{"paintcn", "中文画图", "中文"}, "中文文字生成图像", false, aiChinesePaint, "<文字>")
	aiMadokaCommand  = command.NewNode([]string{"madoka", "円香", "画円香"}, "文字生成图像(円香)", false, aiMadoka, "<文字>")
	aiPaintCommand   = command.NewNode([]string{"paint", "画图", "画画"}, "文字生成图像(普通)", false, aiPaint, "<文字>")
	aiPromptCommand  = command.NewNode([]string{"prompt", "咒语生成"}, "生成文字转图像的咒语", false, aiPrompt, "<开头的字>")
)

var aiCommand = command.NewParent([]string{"ai", "人工智能"}, "AI相关指令",
	aiWaifuCommand,
	aiWaifu2Command,
	aiMadokaCommand,
	aiPaintCommand,
	aiPromptCommand,
	aiPaintCNCommand,
)

func init() {
	command.AddCommand(aiCommand)
}

// hugging face utils

func generateHuggingFaceText(model string, args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在生成文字...."))
	_ = qq.SendGroupMessage(reply)

	inputs := strings.Join(args, " ")

	txt, err := huggingface.GetGeneratedText(model,
		huggingface.NewParam(
			huggingface.Input(inputs),
		),
	)

	if err != nil {
		return err
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("文字生成: %s", txt))

	return qq.SendWithRandomRiskyStrategy(msg)
}

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
