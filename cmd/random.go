package cmd

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/modules/command"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/common-utils/datetime"
	"math/rand"
	"time"
)

func randomMember(args []string, source *command.MessageSource) error {
	rand.Seed(time.Now().UnixMicro())
	members := qq2.ValGroupInfo.Members

	if len(members) == 0 {
		reply := qq2.CreateReply(source.Message).Append(message.NewText("群成员列表为空。"))
		_ = qq2.SendGroupMessage(reply)
		return nil
	}

	chosen := members[rand.Intn(len(members))]
	reply := message.NewSendingMessage().Append(message.NewAt(chosen.Uin))
	return qq2.SendGroupMessage(reply)
}

func randomMessage(args []string, source *command.MessageSource) error {

	msg, err := qq2.GetRandomGroupMessage(source.Message.GroupCode)
	if err != nil {
		return err
	} else if msg == nil {
		return fmt.Errorf("随机消息为空")
	}

	reply := message.NewSendingMessage()
	var nick string
	if msg.Sender.CardName == "" {
		nick = msg.Sender.Nickname
	} else {
		nick = msg.Sender.CardName
	}
	reply.Append(qq2.NewTextfLn("%s 在 %s 说过: ", nick, datetime.FormatSeconds(int64(msg.Time))))
	for _, element := range msg.Elements {
		switch element.(type) {
		case *message.ReplyElement:
			continue
		case *message.ForwardElement:
			continue
		default:
			break
		}
		reply.Append(element)
	}

	return qq2.SendGroupMessage(reply)
}

func randomEssence(args []string, source *command.MessageSource) error {

	rand.Seed(time.Now().UnixMicro())

	gpDist, err := source.Client.GetGroupEssenceMsgList(source.Message.GroupCode)

	// why empty ? not sure but let's try other method
	if len(gpDist) == 0 {
		logger.Warnf("群消息為空，正在使用第 2 種方式獲取")
		gpDist, err = source.Client.GetGroupEssenceMsgList(qq2.ValGroupInfo.Uin)
	}

	// why empty ? not sure but let's try other method
	if len(gpDist) == 0 {
		logger.Warnf("群消息為空，正在使用第 3 種方式獲取")
		gpDist, err = bot.Instance.GetGroupEssenceMsgList(source.Message.GroupCode)
	}

	// why empty ? not sure but let's try other method
	if len(gpDist) == 0 {
		logger.Warnf("群消息為空，正在使用第 4 種方式獲取")
		gpDist, err = bot.Instance.GetGroupEssenceMsgList(qq2.ValGroupInfo.Uin)
	}

	if err != nil {
		logger.Warnf("获取群精华消息列表失败: %v", source.Message.GroupCode)
		return err
	}

	if len(gpDist) == 0 {
		reply := qq2.CreateReply(source.Message).Append(message.NewText("群精华消息列表为空。"))
		_ = qq2.SendGroupMessage(reply)
		return nil
	}

	chosen := gpDist[rand.Intn(len(gpDist))]

	seq := int64(chosen.MessageID)
	essence, err := qq2.GetGroupMessage(source.Message.GroupCode, seq)

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

	return qq2.SendGroupMessage(msg)
}

var (
	randomEssenceCommand = command.NewNode([]string{"essence", "群精华"}, "获取随机一条群精华消息", false, randomEssence)
	randomMemberCommand  = command.NewNode([]string{"member", "成员"}, "随机群成员指令", false, randomMember)
	randomMessageCommand = command.NewNode([]string{"message", "msg", "群消息"}, "随机群消息指令", false, randomMessage)
)

var randomCommand = command.NewParent([]string{"random", "随机"}, "随机指令",
	randomMemberCommand,
	randomEssenceCommand,
	randomMessageCommand,
)

func init() {
	command.AddCommand(randomCommand)
}
