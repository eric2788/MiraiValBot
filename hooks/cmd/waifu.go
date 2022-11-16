package cmd

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/waifu"
	"strconv"
	"strings"
	"sync"
)

func getWaifuMultiple(args []string, source *command.MessageSource) error {

	amountStr, tags, r18 := args[0], strings.Split(args[1], ","), args[2] == "true"

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
			waifu.WithR18(r18),
		),
	)

	if err != nil {
		return err
	}

	forwarder := message.NewForwardMessage()
	wg := &sync.WaitGroup{}

	for _, img := range imgs {
		wg.Add(1)
		if r18 {
			go fetchImageToForwardR18(forwarder, img.Url, wg)
		} else {
			go fetchImageToForward(forwarder, img.Url, wg)
		}
	}

	msg := message.NewSendingMessage()
	msg.Append(forwarder)

	return qq.SendGroupMessage(msg)
}

func fetchImageToForwardR18(forwarder *message.ForwardMessage, url string, wg *sync.WaitGroup) {
	defer wg.Done()
	msg := message.NewSendingMessage()
	img, err := qq.NewImageByUrl(url)
	if err != nil {
		logger.Errorf("尝试获取图片 %s 失败: %v, 将使用URL链接。", url, err)
		msg.Append(qq.NewTextf("[图片]: %s", url))
	} else {
		img.Flash = true
		url, err := bot.Instance.GetGroupImageDownloadUrl(img.FileId, qq.ValGroupInfo.Code, img.Md5)
		if err != nil {
			logger.Errorf("构建闪照失败: %v", err)
		} else {
			img.Url = url
		}
		msg.Append(img)
	}
	forwarder.AddNode(qq.NewForwardNode(msg))
}

var waifuCommand = command.NewNode([]string{"waifu", "setu", "色图"}, "色图指令", false, getWaifuMultiple, "<数量>", "<标签>", "<R18>")

func init() {
	command.AddCommand(waifuCommand)
}
