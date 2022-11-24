package cmd

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/huggingface"
	"github.com/eric2788/MiraiValBot/services/imgtag"
)

func aiWaifu(args []string, source *command.MessageSource) error {

	// model should sort by best quality
	return generateHuggingFaceImage(args, source, false,
		"Linaqruf/anything-v3.0",
		"eimiss/EimisAnimeDiffusion_1.0v",
		"hakurei/waifu-diffusion",
		"Nilaier/Waifu-Diffusers",
	)
}

func aiPaint(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage(args, source, true,
		"prompthero/openjourney",
		"runwayml/stable-diffusion-v1-5",
		"CompVis/stable-diffusion-v1-4",
	)
}

func aiMadoka(args []string, source *command.MessageSource) error {
	return generateHuggingFaceImage(args, source, false,
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
	return generateHuggingFaceImage(args, source, true,
		"IDEA-CCNL/Taiyi-Stable-Diffusion-1B-Chinese-EN-v0.1",
		"IDEA-CCNL/Taiyi-Stable-Diffusion-1B-Chinese-v0.1",
	)
}

func aiTags(args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)

	content := qq.ParseMsgContent(source.Message.Elements)
	imgs := content.Images

	// 支援 reply 圖片輸入指令
	if len(imgs) == 0 && len(content.Replies) > 0 {
		for _, ele := range source.Message.Elements {
			if reply, ok := ele.(*message.ReplyElement); ok {
				imgs = qq.ParseMsgContent(reply.Elements).Images
			}
		}
	}

	if len(imgs) == 0 {
		reply.Append(message.NewText("找不到图片, 请附带图片!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在识别图片...."))
	_ = qq.SendGroupMessage(reply)

	reply = qq.CreateReply(source.Message)

	img := imgs[0]
	tag, nsfw, err := imgtag.GetTagsFromImage(img)
	if err != nil {
		return err
	}

	reply.Append(qq.NewTextfLn("图片识别标签: %s", strings.Join(tag, ", ")))
	reply.Append(qq.NewTextf("NSFW: %t", nsfw))
	return qq.SendGroupMessage(reply)
}

func aiSearchTags(args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在搜索标签...."))
	_ = qq.SendGroupMessage(reply)

	tags, err := imgtag.SearchTags(args[0])
	if err != nil {
		return err
	}

	reply = qq.CreateReply(source.Message)
	reply.Append(qq.NewTextfLn("%s 的搜索结果:", args[0]))

	for tag, cn := range tags {
		reply.Append(qq.NewTextfLn("%s: %s", tag, cn))
	}

	return qq.SendWithRandomRiskyStrategy(reply)
}

// aiWaifu2 generate high quality waifu image (but still nsfw filtered)
func aiWaifu2(args []string, source *command.MessageSource) error {
	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在生成图像...."))
	_ = qq.SendGroupMessage(reply)

	inputs := strings.Join(args, " ")

	var err error
	var bb [][]byte

	apis := []*huggingface.SpaceApi{
		huggingface.NewSpaceApi("akhaliq-anything-v3-0",
			"anything v3",
			inputs,
			7.5,
			35,
			720,
			720,
			0,
			nil,
			0.5,
			huggingface.BadPrompt,
		),
		huggingface.NewSpaceApi("fkunn1326-animestyle-diffusionmodels",
			"EimisAnimeDiffusion_1.0v",
			inputs,
			7.5,
			35,
			720,
			720,
			0,
			nil,
			0.5,
			huggingface.BadPrompt,
		),
	}

	for _, api := range apis {
		bb, err = api.UseWebsocketHandler().GetResultImages()
		if err == nil {
			break
		} else {
			logger.Errorf("使用model %s 生成图像时出现错误: %v", api.Id, err)
		}
	}

	if err != nil {
		return err
	} else if len(bb) == 0 {
		return errors.New("没有图片被生成")
	}

	img, err := qq.NewImageByByte(bb[0])
	if err != nil {
		return err
	}

	msg := qq.CreateReply(source.Message)
	msg.Append(img)
	return qq.SendGroupMessage(msg)
}

var (
	aiWaifuCommand      = command.NewNode([]string{"waifu"}, "文字生成二次元图", false, aiWaifu, "<文字>")
	aiWaifu2Command     = command.NewNode([]string{"waifu2"}, "文字生成二次元图(高质量,但超慢)", false, aiWaifu2, "<文字>")
	aiPaintCNCommand    = command.NewNode([]string{"paintcn", "中文画图", "中文"}, "中文文字生成图像", false, aiChinesePaint, "<文字>")
	aiMadokaCommand     = command.NewNode([]string{"madoka", "円香", "画円香"}, "文字生成图像(円香)", false, aiMadoka, "<文字>")
	aiPaintCommand      = command.NewNode([]string{"paint", "画图", "画画"}, "文字生成图像(普通)", false, aiPaint, "<文字>")
	aiPromptCommand     = command.NewNode([]string{"prompt", "咒语生成"}, "生成文字转图像的咒语", false, aiPrompt, "<开头的字>")
	aiTagCommand        = command.NewNode([]string{"tags", "标签", "分析"}, "分析图片获取标签", false, aiTags)
	aiSearchTagsCommand = command.NewNode([]string{"searchtags", "search", "搜索标签"}, "中文搜索图片标签", false, aiSearchTags, "<中文关键词>")
)

var aiCommand = command.NewParent([]string{"ai", "人工智能"}, "AI相关指令",
	aiWaifuCommand,
	aiWaifu2Command,
	aiMadokaCommand,
	aiPaintCommand,
	aiPromptCommand,
	aiPaintCNCommand,
	aiTagCommand,
	aiSearchTagsCommand,
)

func init() {
	command.AddCommand(aiCommand)
}

// hugging face utils

const badPrompt = `bad feet, bad foot, lowres, bad anatomy, bad hands, text, error, missing fingers, extra digit, fewer digits, cropped, worst quality, low quality, normal quality, jpeg artifacts, signature, watermark, username, blurry`

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
			huggingface.InputWithoutBracket(inputs),
		)

		txt, err = api.GetGeneratedText()

		if err == nil {
			break
		} else {
			logger.Errorf("使用model %s 生成文字时出现错误: %v", model, err)
		}
	}

	if err != nil {
		return err
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("文字生成: %s", txt))

	return qq.SendWithRandomRiskyStrategy(msg)
}

func generateHuggingFaceImage(args []string, source *command.MessageSource, random bool, models ...string) error {
	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在生成图像...."))
	_ = qq.SendGroupMessage(reply)

	inputs := strings.Join(args, " ")

	if random {
		// shuffle
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(models), func(i, j int) {
			models[i], models[j] = models[j], models[i]
		})
	}

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
			logger.Errorf("使用model %s 生成图像时出现错误: %v", model, err)
		}

	}

	if err != nil {
		return err
	}

	img, err := qq.NewImageByByte(b)

	if err != nil {
		return err
	}

	msg := qq.CreateReply(source.Message)
	msg.Append(img)
	return qq.SendGroupMessage(msg)
}
