package cmd

import (
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/command"
	qq "github.com/eric2788/MiraiValBot/qq"
)

func yearlyCheck(args []string, source *command.MessageSource) error {
	file.UpdateStorage(func() {
		file.DataStorage.Setting.YearlyCheck = !file.DataStorage.Setting.YearlyCheck
	})

	reply := qq.CreateReply(source.Message)
	var s string
	if file.DataStorage.Setting.YearlyCheck {
		s = "每年"
	} else {
		s = "每月"
	}
	reply.Append(qq.NewTextf("已设置群精华消息检查间隔为 %s", s))
	return qq.SendGroupMessage(reply)
}

func verbose(args []string, source *command.MessageSource) error {
	file.UpdateStorage(func() {
		file.DataStorage.Setting.Verbose = !file.DataStorage.Setting.Verbose
	})
	reply := qq.CreateReply(source.Message)
	reply.Append(qq.NewTextf("成功切换广播状态为 %v", file.DataStorage.Setting.Verbose))
	return qq.SendGroupMessage(reply)

}

func verboseDelete(args []string, source *command.MessageSource) error {
	file.UpdateStorage(func() {
		file.DataStorage.Setting.VerboseDelete = !file.DataStorage.Setting.VerboseDelete
	})
	reply := qq.CreateReply(source.Message)
	reply.Append(qq.NewTextf("已成功设置显示撤回消息为 %v", file.DataStorage.Setting.VerboseDelete))
	return qq.SendGroupMessage(reply)
}

var (
	verboseCommand       = command.NewNode([]string{"verbose", "切换广播"}, "切换是否广播监听状态", true, verbose)
	verboseDeleteCommand = command.NewNode([]string{"telldelete"}, "显示撤回的消息", true, verboseDelete)
	yearlyCheckCommand   = command.NewNode([]string{"yearly"}, "设置群精华消息检查间隔", true, yearlyCheck)
)

var settingCommand = command.NewParent([]string{"setting", "设定"}, "设定指令",
	verboseCommand,
	verboseDeleteCommand,
	yearlyCheckCommand,
)

func init() {
	command.AddCommand(settingCommand)
}
