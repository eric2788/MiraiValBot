package cmd

import (
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/aidraw"
)

func register(args []string, source *command.MessageSource) error {
	return auth(args, source, true)
}

func login(args []string, source *command.MessageSource) error {
	return auth(args, source, false)
}

var (
	loginCommand    = command.NewNode([]string{"login", "登录"}, "登录到 Sexy AI", true, login, "<电邮>", "[一次性密码]")
	registerCommand = command.NewNode([]string{"register", "注册"}, "注册 Sexy AI", true, register, "<电邮>", "[一次性密码]")
)

var sexyAICommand = command.NewParent([]string{"sexyai", "sai"}, "Sexy AI",
	loginCommand,
	registerCommand,
)

func init() {
	command.AddCommand(sexyAICommand)
}

func auth(args []string, source *command.MessageSource, register bool) error {
	email := args[0]

	var otp string = ""
	if len(args) > 1 {
		otp = args[1]
	}

	reply := qq.CreateReply(source.Message)

	if otp == "" {
		// authentication step 1: send otp
		success, err := aidraw.SaiRequestOTP(email, register)
		if err != nil {
			return err
		}
		if !success {
			reply.Append(qq.NewTextf("发送一次性密码到电邮 %v 失败", email))
		} else {
			reply.Append(qq.NewTextf("已成功发送一次性密码到电邮 %v, 请检查你的垃圾邮件", email))
		}
	} else {
		// authentication step 2: verify otp
		result, err := aidraw.SaiAuth(email, otp, register)
		if err != nil {
			return err
		}
		reply.Append(qq.NewTextf("成功以 %v 的身份登录。", result))
	}

	return qq.SendGroupMessage(reply)

}
