package cmd

import (
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/huggingface"
)

// putting hugging face api here...

func aiPaint(args []string, source *command.MessageSource) error {

	reply := qq.CreateReply(source.Message)

	if len(args) == 0 {
		reply.Append(message.NewText("参数不能为空!"))
		return qq.SendGroupMessage(reply)
	}

	reply.Append(qq.NewTextf("正在生成图像...."))
	_ = qq.SendGroupMessage(reply)

	inputs := strings.Join(args, " ")

	b, err := huggingface.GetResultImage("Nilaier/Waifu-Diffusers",
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

var (
	aiPaintCommand = command.NewNode([]string{"paint", "画图", "画画", "画"}, "文字生成图像", false, aiPaint, "<文字>")
)

var aiCommand = command.NewParent([]string{"ai", "人工智能"}, "AI相关指令",
	aiPaintCommand,
)

func init() {
	command.AddCommand(aiCommand)
}
