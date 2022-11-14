package cmd

import (
	"fmt"
	"strings"

	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"golang.org/x/exp/maps"
)

func addWordCount(args []string, source *command.MessageSource) error {

	msg := args[0]

	reply := qq.CreateReply(source.Message)

	_, ok := file.DataStorage.WordCounts[msg]

	if ok {
		reply.Append(qq.NewTextf("字词 %q 已经启动记录。", msg))
		return qq.SendGroupMessage(reply)
	}

	file.UpdateStorage(func() {
		file.DataStorage.WordCounts[msg] = make(map[int64]int64)
	})

	reply.Append(qq.NewTextf("开始记录字词 %q", msg))
	return qq.SendGroupMessage(reply)
}

func removeWordCount(args []string, source *command.MessageSource) error {
	msg := args[0]

	reply := qq.CreateReply(source.Message)

	_, ok := file.DataStorage.WordCounts[msg]

	if !ok {
		reply.Append(qq.NewTextf("该字词 %q 没有启动记录。", msg))
		return qq.SendGroupMessage(reply)
	}

	file.UpdateStorage(func() {
		delete(file.DataStorage.WordCounts, msg)
	})

	reply.Append(qq.NewTextf("成功中止及清空字词记录 %q", msg))
	return qq.SendGroupMessage(reply)
}

func listWorldCount(args []string, source *command.MessageSource) error {
	word := args[0]
	msg := qq.CreateReply(source.Message)

	counts, ok := file.DataStorage.WordCounts[word]

	if !ok {
		msg.Append(qq.NewTextf("未知字词 %q, 可用字词: %s", word, strings.Join(maps.Keys(file.DataStorage.WordCounts), ", ")))
		return qq.SendGroupMessage(msg)
	}

	msg.Append(qq.NewTextLn("字词 %q 的群聊记录次数:"))

	for uid, times := range counts {

		info := qq.FindGroupMember(uid)
		var name string
		if info != nil {
			name = info.Nickname
		} else {
			name = fmt.Sprintf("(UID: %d)", uid)
		}

		msg.Append(qq.NewTextfLn("%s 说了 %d 次", name, times))
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

var (
	addWordCountCommand    = command.NewNode([]string{"add", "新增"}, "启动字词记录", true, addWordCount, "<字词>")
	removeWordCountCommand = command.NewNode([]string{"remove", "移除"}, "移除字词记录", true, removeWordCount, "<字词>")
	listWorldCountCommand  = command.NewNode([]string{"list", "列表"}, "显示字词记录列表", false, listWorldCount, "<字词>")
)

var countCommand = command.NewParent([]string{"count", "wordcount", "字词记录"}, "字词记录指令",
	addWordCountCommand,
	removeWordCountCommand,
	listWorldCountCommand,
)

func init() {
	command.AddCommand(countCommand)
}
