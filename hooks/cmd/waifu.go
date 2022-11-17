package cmd

import (
	"errors"
	"fmt"
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
	} else if amount > 40 {
		return errors.New("最高每次获取40张。")
	}
	isKeyword := len(tags) == 1

	var search waifu.Searcher
	if isKeyword {
		search = waifu.WithKeyword(tags[0])
	} else {
		search = waifu.WithTags(tags...)
	}

	reply := qq.CreateReply(source.Message)
	_ = qq.SendGroupMessage(reply.Append(qq.NewTextf("正在索取 %s 的相关图片...", strings.Join(tags, ","))))

	imgs, err := waifu.GetRandomImages(
		waifu.NewOptions(
			search,
			waifu.WithAmount(amount),
			waifu.WithR18(false), // 为了安全
		),
	)

	if err != nil {
		return err
	} else if len(imgs) == 0 {
		return fmt.Errorf("搜索 %s 的结果为空。", strings.Join(tags, ","))
	}

	forwarder := message.NewForwardMessage()
	wg := &sync.WaitGroup{}

	for _, img := range imgs {
		wg.Add(1)
		go fetchImageToForward(forwarder, img.Url, wg)
	}

	wg.Wait()
	return qq.SendGroupForwardMessage(forwarder)
}

var waifuCommand = command.NewNode([]string{"waifu", "setu", "色图"}, "色图指令", false, getWaifuMultiple, "<数量>", "<标签>")

func init() {
	command.AddCommand(waifuCommand)
}
