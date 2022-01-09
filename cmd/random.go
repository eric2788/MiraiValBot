package cmd

import (
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/utils/qq"
	"math/rand"
	"time"
)

func randomMember(args []string, source *command.MessageSource) error {
	rand.Seed(time.Now().UnixMicro())
	members := qq.ValGroupInfo.Members

	if len(members) == 0 {
		reply := qq.CreateReply(source.Message).Append(message.NewText("群成员列表为空。"))
		source.Client.SendGroupMessage(source.Message.GroupCode, reply)
		return nil
	}

	chosen := members[rand.Intn(len(members))]
	reply := message.NewSendingMessage().Append(message.NewAt(chosen.Uin))
	source.Client.SendGroupMessage(source.Message.GroupCode, reply)
	return nil
}

func randomEssence(args []string, source *command.MessageSource) error {

	rand.Seed(time.Now().UnixMicro())

	gpDist, err := source.Client.GetGroupEssenceMsgList(source.Message.GroupCode)

	if err != nil {
		logger.Warnf("获取群精华消息列表失败: %v", source.Message.GroupCode)
		return err
	}

	if len(gpDist) == 0 {
		reply := qq.CreateReply(source.Message).Append(message.NewText("群精华消息列表为空。"))
		source.Client.SendGroupMessage(source.Message.GroupCode, reply)
		return nil
	}

	chosen := gpDist[rand.Intn(len(gpDist))]

	seq := int64(chosen.MessageID)
	essence, err := qq.GetGroupMessage(source.Message.GroupCode, seq)

	if err != nil {
		logger.Warnf("获取群精华消息失败: %+v", chosen)
		return err
	}
	msg := message.NewSendingMessage()

	if essence != nil {
		for _, element := range essence.Elements {
			msg.Append(element)
		}
	} else {
		msg.Append(message.NewText("没有群精华消息"))
	}

	source.Client.SendGroupMessage(source.Message.GroupCode, msg)
	return nil
}

var (
	randomEssenceCommand = command.NewNode([]string{"essence", "群精华"}, "获取随机一条群精华消息", false, randomEssence)
	randomMemberCommand  = command.NewNode([]string{"member", "成员"}, "随机群成员指令", false, randomMember)
)

var randomCommand = command.NewParent([]string{"random", "随机"}, "随机指令",
	randomMemberCommand,
	randomEssenceCommand,
)

func init() {
	command.AddCommand(randomCommand)
}
