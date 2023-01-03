package cmd

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eric2788/MiraiValBot/utils/misc"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/waifu"
)

func getWaifuR18(args []string, source *command.MessageSource) error {

	tags := []string{""}

	if len(args) > 0 {
		tags = strings.Split(strings.Join(args, " "), ",")
	}

	isKeyword := len(tags) == 1

	var search waifu.Searcher
	if isKeyword {
		search = waifu.WithKeyword(tags[0])
	} else {
		search = waifu.WithTags(tags...)
	}

	reply := qq.CreateReply(source.Message)
	_ = qq.SendGroupMessage(reply.Append(qq.NewTextf("正在索取 %s 的相关图片...", strings.Join(tags, ", "))))

	imgs, err := waifu.GetRandomImages(
		waifu.NewOptions(
			search,
			waifu.WithAmount(1),
			waifu.WithR18(true),
		),
	)

	if err != nil {
		return err
	} else if len(imgs) == 0 {
		return fmt.Errorf("搜索 %s 的结果为空。", strings.Join(tags, ","))
	}

	data := imgs[0]

	var img *message.GroupImageElement
	if data.Image != nil {
		img, err = qq.NewImageByByte(data.Image)
	} else if data.Url != "" {
		img, err = qq.NewImageByUrl(data.Url)
	}

	if err != nil {
		return err
	} else if img == nil {
		return errors.New("图片获取失败")
	}

	reply = qq.CreateReply(source.Message)
	reply.Append(img)

	// 三十秒后撤回消息
	return qq.SendGroupMessageAndRecall(reply, time.Second*30)
}

func getWaifuMultiple(args []string, source *command.MessageSource) error {

	amountStr, tags := args[0], []string{""}

	if len(args) > 1 {
		tags = strings.Split(strings.Join(args[1:], " "), ",")
	}

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
	_ = qq.SendGroupMessage(reply.Append(qq.NewTextf("正在索取 %s 的相关图片...", strings.Join(tags, ", "))))

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
		if len(img.Image) > 0 {
			go misc.FetchImageByteToForward(forwarder, img.Image, wg)
		} else if img.Url != "" {
			go misc.FetchImageToForward(forwarder, img.Url, wg)
		} else {
			// nothing
			logger.Warnf("图片 %d (%q) 的 Url 和 Image 均为空。", img.Pid, img.Title)
			wg.Done()
		}
	}

	wg.Wait()
	return qq.SendGroupForwardMessage(forwarder)
}

var (
	waifuMultipleCommand = command.NewNode([]string{"multi", "multiple", "多张"}, "索取多张色图(没有r18)", false, getWaifuMultiple, "<数量>", "[标签]")
	waifuR18Command      = command.NewNode([]string{"r18", "瑟瑟"}, "索取一张可以是r18的色图", false, getWaifuR18, "[标签]")
)

var waifuCommand = command.NewParent([]string{"waifu", "setu", "色图"}, "色图指令",
	waifuMultipleCommand,
	waifuR18Command,
)

func init() {
	command.AddCommand(waifuCommand)
}
