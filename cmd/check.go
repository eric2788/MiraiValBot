package cmd

import (
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"time"
)

func check(args []string, source *command.MessageSource) error {

	ats := qq.ParseMsgContent(source.Message.Elements).At

	for _, at := range ats {
		if member := qq.FindGroupMember(at); member != nil {
			msg := message.NewSendingMessage()

			msg.Append(qq.NewTextfLn("UID: %d", member.Uin))
			msg.Append(qq.NewTextfLn("名称: %s", member.Nickname))
			msg.Append(qq.NewTextfLn("显示名称: %s", member.DisplayName()))
			msg.Append(qq.NewTextfLn("卡片名称: %s", member.CardName))
			msg.Append(qq.NewTextfLn("性别: %s", genderName(member.Gender)))
			msg.Append(qq.NewTextfLn("加入日期: %s", toTime(member.JoinTime)))
			msg.Append(qq.NewTextfLn("权限: %s", permissionName(member.Permission)))
			msg.Append(qq.NewTextfLn("最后发言时间: %s", toTime(member.LastSpeakTime)))
			msg.Append(qq.NewTextfLn("等级: %d", member.Level))
			msg.Append(qq.NewTextf("特别头衔: %s", member.SpecialTitle))

			source.Client.SendGroupMessage(source.Message.GroupCode, msg)
		}
	}

	return nil
}

func toTime(ts int64) string {
	return time.UnixMilli(ts * 1000).Format("2006-01-02 15:04:05")
}
func genderName(b byte) string {
	switch b {
	case 0:
		return "男"
	case 1:
		return "女"
	default:
		return "未知"
	}
}

func permissionName(permission client.MemberPermission) string {
	switch permission {
	case client.Member:
		return "群成员"
	case client.Administrator:
		return "群管理员"
	case client.Owner:
		return "群主"
	default:
		return "未知"
	}
}

var checkCommand = command.NewNode([]string{"check", "查成分", "查"}, "查成分", false, check, "<用户>")

func init() {
	command.AddCommand(checkCommand)
}
