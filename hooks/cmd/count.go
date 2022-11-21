package cmd

import (
	"fmt"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/file"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
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

	type WordCount struct {
		Uid   int64
		Times int64
	}

	wc := make([]WordCount, 0)

	for uid, times := range counts {
		wc = append(wc, WordCount{
			Uid:   uid,
			Times: times,
		})
	}

	slices.SortStableFunc(wc, func(a, b WordCount) bool {
		return a.Times > b.Times
	})

	msg.Append(qq.NewTextfLn("字词 %q 的群聊记录次数: (由高到低)", word))

	for d, c := range wc {

		info := qq.FindGroupMember(c.Uid)
		var name string
		if info != nil {
			name = info.DisplayName()
		} else {
			name = fmt.Sprintf("(UID: %d)", c.Uid)
		}

		msg.Append(qq.NewTextfLn("%d. %s 说了 %d 次", d+1, name, c.Times))
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

func rankWords(args []string, source *command.MessageSource) error {

	msg := qq.CreateReply(source.Message)

	if len(file.DataStorage.WordCounts) == 0 {
		msg.Append(message.NewText("没有正在记录的字词。"))
		return qq.SendGroupMessage(msg)
	}

	var counts = make(map[string]int64, len(file.DataStorage.WordCounts))
	for word, users := range file.DataStorage.WordCounts {
		c := int64(0)
		for _, times := range users {
			c += times
		}
		counts[word] = c
	}

	sortedKeys := maps.Keys(counts)

	slices.SortStableFunc(sortedKeys, func(a, b string) bool {
		return a > b
	})

	msg.Append(qq.NewTextf("群聊字词记录排行: (由高到低)"))

	for i, word := range sortedKeys {
		times := counts[word]
		msg.Append(qq.NewTextfLn("%d. %q 说了 %d 次", i+1, word, times))
	}

	return qq.SendWithRandomRiskyStrategy(msg)
}

var (
	addWordCountCommand    = command.NewNode([]string{"add", "新增"}, "启动字词记录", true, addWordCount, "<字词>")
	removeWordCountCommand = command.NewNode([]string{"remove", "移除"}, "移除及清空字词记录", true, removeWordCount, "<字词>")
	listWordCountCommand   = command.NewNode([]string{"list", "列表", "rank", "排行"}, "显示该字词的玩家记录列表(带排行)", false, listWorldCount, "<字词>")
	rankWordsCommand       = command.NewNode([]string{"words", "wordrank", "字词", "字词列表"}, "显示群聊字词列表(带排行)", false, rankWords)
)

var countCommand = command.NewParent([]string{"count", "wordcount", "字词记录"}, "字词记录指令",
	addWordCountCommand,
	removeWordCountCommand,
	listWordCountCommand,
	rankWordsCommand,
)

func init() {
	command.AddCommand(countCommand)
}
