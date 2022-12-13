package cmd

import (
	"math/rand"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/copywriting"
)

func writeCP(args []string, source *command.MessageSource) error {
	atk, def := args[0], args[1]
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
	atk := args[0]
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
	atk := args[0]

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

func writeZuowen(args []string, source *command.MessageSource) error {
	atk := args[0]
	rand.Seed(time.Now().UnixNano())
	msg := message.NewSendingMessage()
	list, err := copywriting.GetRanranList()
	if err != nil {
		return err
	}
	random := list[rand.Intn(len(list))]
	msg.Append(message.NewText(strings.ReplaceAll(random.Text, random.Person, atk)))
	return qq.SendWithRandomRiskyStrategy(msg)
}

var (
	writeCPCommand     = command.NewNode([]string{"cp", "组合"}, "帮两个目标写CP", false, writeCP, "<攻>", "<受>")
	writeFaBingCommand = command.NewNode([]string{"fabing", "发病"}, "对着一个目标发病", false, writeFaBing, "<对象>")
	writeFaDianCommand = command.NewNode([]string{"fadian", "发癫"}, "对着一个目标发癫", false, writeFaDian, "<对象>")
	writeZuowenCommand = command.NewNode([]string{"zuowen", "作文"}, "对着一个目标写小作文", false, writeZuowen, "<对象>")
)

var writeCommand = command.NewParent([]string{"write", "写"}, "写点什么",
	writeCPCommand,
	writeFaBingCommand,
	writeFaDianCommand,
	writeZuowenCommand,
)

func init() {
	command.AddCommand(writeCommand)
}
