package cmd

import (
	"strconv"
	"strings"
	"sync"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/waifu"
)

func getWaifuMultiple(args []string, source *command.MessageSource) error {

	amountStr, tags := args[0], strings.Split(args[1], ",")

	amount, err := strconv.Atoi(amountStr)
	if err != nil {
		return err
	}
	isKeyword := len(tags) == 1

	var search waifu.Searcher
	if isKeyword {
		search = waifu.WithKeyword(tags[0])
	} else {
		search = waifu.WithTags(tags...)
	}

	imgs, err := waifu.GetRandomImages(
		waifu.NewOptions(
			search,
			waifu.WithAmount(amount),
			waifu.WithR18(false), // 为了安全
		),
	)

	if err != nil {
		return err
	}

	forwarder := message.NewForwardMessage()
	wg := &sync.WaitGroup{}

	for _, img := range imgs {
		wg.Add(1)
		go fetchImageToForward(forwarder, img.Url, wg)
	}

	return qq.SendGroupForwardMessage(forwarder)
}

var waifuCommand = command.NewNode([]string{"waifu", "setu", "色图"}, "色图指令", false, getWaifuMultiple, "<数量>", "<标签>")

func init() {
	command.AddCommand(waifuCommand)
}
