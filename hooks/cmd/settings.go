package cmd

import (
	"fmt"
	"strconv"

	"github.com/eric2788/MiraiValBot/internal/file"
	qq "github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
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

func timesPerNotify(args []string, source *command.MessageSource) error {
	t, err := strconv.Atoi(args[0])
	if err != nil {
		return err
	}
	file.UpdateStorage(func() {
		file.DataStorage.Setting.TimesPerNotify = t
	})
	reply := qq.CreateReply(source.Message)
	reply.Append(qq.NewTextf("已成功设置字词记录提醒间隔为 %d 次", t))
	return qq.SendGroupMessage(reply)
}

func msgSeqAfter(args []string, source *command.MessageSource) error {
	msg := qq.CreateReply(source.Message)

	if len(args) == 0 {
		msg.Append(qq.NewTextf("目前信息获取序列是 %d", file.DataStorage.Setting.MsgSeqAfter))
		return qq.SendGroupMessage(msg)
	}

	seq, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return err
	}
	if seq <= 0 {
		return fmt.Errorf("序列必须大于0")
	}
	file.UpdateStorage(func() {
		file.DataStorage.Setting.MsgSeqAfter = seq
	})

	msg.Append(qq.NewTextf("成功设置信息获取序列为 %d", seq))
	return qq.SendGroupMessage(msg)
}

func tagClassifyLimit(args []string, source *command.MessageSource) error {
	msg := qq.CreateReply(source.Message)

	if len(args) == 0 {
		msg.Append(qq.NewTextf("目前标签鉴别强度是 %d", file.DataStorage.Setting.TagClassifyLimit))
		return qq.SendGroupMessage(msg)
	}

	limit, err := strconv.ParseFloat(args[0], 64)
	if err != nil {
		return err
	}
	if limit < 0 || limit > 1 {
		return fmt.Errorf("强度必须在0-1之间")
	}

	file.UpdateStorage(func() {
		file.DataStorage.Setting.TagClassifyLimit = limit
	})

	msg.Append(qq.NewTextf("成功设置标签鉴别强度为 %f", limit))
	return qq.SendGroupMessage(msg)
}

var (
	verboseCommand          = command.NewNode([]string{"verbose", "切换广播"}, "切换是否广播监听状态", true, verbose)
	verboseDeleteCommand    = command.NewNode([]string{"telldelete"}, "显示撤回的消息", true, verboseDelete)
	yearlyCheckCommand      = command.NewNode([]string{"yearly"}, "设置群精华消息检查间隔", true, yearlyCheck)
	timerPerNotifyCommand   = command.NewNode([]string{"notifytimes", "提醒间隔"}, "设置字词记录提醒间隔", true, timesPerNotify)
	msgSeqAfterCommand      = command.NewNode([]string{"msgseq", "信息序列"}, "设置信息获取序列", true, msgSeqAfter, "[序列]")
	tagClassifyLimitCommand = command.NewNode([]string{"taglimit", "标签强度"}, "设置标签鉴别强度", true, tagClassifyLimit, "[强度]")
)

var settingCommand = command.NewParent([]string{"setting", "设定"}, "设定指令",
	verboseCommand,
	verboseDeleteCommand,
	yearlyCheckCommand,
	fetchEssenceCommand,
	timerPerNotifyCommand,
	msgSeqAfterCommand,
	tagClassifyLimitCommand,
)

func init() {
	command.AddCommand(settingCommand)
}
