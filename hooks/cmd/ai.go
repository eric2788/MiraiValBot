package cmd

import (
	"math/rand"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/huggingface"
)

func aiWaifu(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage(args, source,
		"hakurei/waifu-diffusion",
		"Nilaier/Waifu-Diffusers",
	)
}

func aiPaint(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage(args, source,
		"runwayml/stable-diffusion-v1-5",
		"CompVis/stable-diffusion-v1-4",
	)
}

func aiMadoka(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage(args, source,
		"yuk/madoka-waifu-diffusion",
	)
}

func aiPrompt(args []string, source *command.MessageSource) error {
	return generateHuggingFaceText(args, source,
		"Gustavosta/MagicPrompt-Stable-Diffusion",
		"DrishtiSharma/StableDiffusion-Prompt-Generator-GPT-Neo-125M",
	)
}

func aiChinesePaint(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage(args, source,
		"IDEA-CCNL/Taiyi-Stable-Diffusion-1B-Chinese-v0.1",
		"IDEA-CCNL/Taiyi-Stable-Diffusion-1B-Chinese-EN-v0.1",
	)
}

var (
	aiWaifuCommand   = command.NewNode([]string{"waifu"}, "文字生成图像(waifu)", false, aiWaifu, "<文字>")
	aiPaintCNCommand = command.NewNode([]string{"paintcn", "中文画图", "中文"}, "中文文字生成图像", false, aiChinesePaint, "<文字>")
	aiMadokaCommand  = command.NewNode([]string{"madoka", "円香", "画円香"}, "文字生成图像(円香)", false, aiMadoka, "<文字>")
	aiPaintCommand   = command.NewNode([]string{"paint", "画图", "画画"}, "文字生成图像(普通)", false, aiPaint, "<文字>")
	aiPromptCommand  = command.NewNode([]string{"prompt", "咒语生成"}, "生成文字转图像的咒语", false, aiPrompt, "<开头的字>")
)

var aiCommand = command.NewParent([]string{"ai", "人工智能"}, "AI相关指令",
	aiWaifuCommand,
	aiMadokaCommand,
	aiPaintCommand,
	aiPromptCommand,
	aiPaintCNCommand,
)

func init() {
	command.AddCommand(aiCommand)
}

// hugging face utils

func generateHuggingFaceText(args []string, source *command.MessageSource, models ...string) error {
	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在生成文字...."))
	_ = qq.SendGroupMessage(reply)

	inputs := strings.Join(args, " ")

	var txt string
	var err error

	for _, model := range models {

		api := huggingface.NewInferenceApi(model,
			huggingface.Input(inputs),
		)

		txt, err = api.GetGeneratedText()

		if err == nil {
			break
		} else {
			logger.Errorf("使用model %s 生成文字时出现错误: %v", err)
		}
	}

	if err != nil {
		return err
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("文字生成: %s", txt))

	return qq.SendWithRandomRiskyStrategy(msg)
}

func generateHuggingFaceImage(args []string, source *command.MessageSource, models ...string) error {
	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在生成图像...."))
	_ = qq.SendGroupMessage(reply)

	inputs := strings.Join(args, " ")

	// 打散随机使用
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(models), func(i, j int) {
		models[i], models[j] = models[j], models[i]
	})

	var err error
	var b []byte
	for _, model := range models {

		api := huggingface.NewInferenceApi(model,
			huggingface.Input(inputs),
		)

		b, err = api.GetResultImage()

		if err == nil {
			break
		} else {
			logger.Errorf("使用model %s 生成图像时出现错误: %v", err)
		}

	}

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
