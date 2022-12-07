package cmd

import (
	"errors"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/copywriting"
	"math/rand"
	"strings"
	"time"
)

func writeCP(args []string, source *command.MessageSource) error {
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) != 2 {
		return errors.New("请@两个人: 前<攻>后<受>")
	}
	atk, def := strings.ReplaceAll(ats[0].Display, "@", ""), strings.ReplaceAll(ats[1].Display, "@", "")
	if atk == def {
		return errors.New("请@两个不同的人")
	}

	rand.Seed(time.Now().UnixNano())
	msg := message.NewSendingMessage()

	list, a, d, err := copywriting.GetCPList()
	if err != nil {
		return err
	}
	random := list[rand.Intn(len(list))]
	replacer := strings.NewReplacer(a, atk, d, def)
	msg.Append(message.NewText(replacer.Replace(random)))

	return qq.SendWithRandomRiskyStrategy(msg)
}

func writeFaBing(args []string, source *command.MessageSource) error {
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) == 0 {
		return errors.New("请@一个人")
	}
	atk := strings.ReplaceAll(ats[0].Display, "@", "")

	rand.Seed(time.Now().UnixNano())
	msg := message.NewSendingMessage()
	list, a, err := copywriting.GetFabingList()
	if err != nil {
		return err
	}
	random := list[rand.Intn(len(list))]
	replacer := strings.NewReplacer(a, atk)
	msg.Append(message.NewText(replacer.Replace(random)))

	return qq.SendWithRandomRiskyStrategy(msg)
}

func writeFaDian(args []string, source *command.MessageSource) error {
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) == 0 {
		return errors.New("请@一个人")
	}
	atk := strings.ReplaceAll(ats[0].Display, "@", "")

	rand.Seed(time.Now().UnixNano())
	msg := message.NewSendingMessage()
	list, a, err := copywriting.GetFadianList()
	if err != nil {
		return err
	}
	random := list[rand.Intn(len(list))]
	replacer := strings.NewReplacer(a, atk)
	msg.Append(message.NewText(replacer.Replace(random)))

	return qq.SendWithRandomRiskyStrategy(msg)
}

var (
	writeCPCommand     = command.NewNode([]string{"cp", "组合"}, "帮两个群友写CP", false, writeCP)
	writeFaBingCommand = command.NewNode([]string{"fabing", "发病"}, "对着一个群友发病", false, writeFaBing)
	writeFaDianCommand = command.NewNode([]string{"fadian", "发癫"}, "对着一个群友发癫", false, writeFaDian)
)

var writeCommand = command.NewParent([]string{"write", "写"}, "写点什么",
	writeCPCommand,
	writeFaBingCommand,
	writeFaDianCommand,
)

func init() {
	command.AddCommand(writeCommand)
}
