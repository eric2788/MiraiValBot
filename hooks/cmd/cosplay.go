package cmd

import (
	"errors"
	"sync"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/cosplayer"
	"github.com/eric2788/MiraiValBot/utils/misc"
)

func cosplaySingle(args []string, source *command.MessageSource) error {

	reply := qq.CreateReply(source.Message)
	_ = qq.SendGroupMessage(reply.Append(qq.NewTextf("正在索取 Cosplayer 图片...")))

	b, err := cosplayer.GetImageRandom()
	if err != nil {
		return err
	}
	img, err := qq.NewImageByByte(b)
	if err != nil {
		return err
	}
	msg := message.NewSendingMessage()
	msg.Append(img)
	return qq.SendGroupMessage(msg)
}

func cosplayMultiple(args []string, source *command.MessageSource) error {

	reply := qq.CreateReply(source.Message)
	_ = qq.SendGroupMessage(reply.Append(qq.NewTextf("正在索取 Cosplayer 图片...")))

	data, err := cosplayer.GetImagesRandom()
	if err != nil {
		return err
	}
	if len(data.Urls) == 0 {
		return errors.New("获取到的Cosplayer图片为空，请再尝试一次")
	}

	forwarder := message.NewForwardMessage()
	title := message.NewSendingMessage()
	title.Append(message.NewText(data.Title))
	forwarder.AddNode(qq.NewForwardNode(title))

	wg := &sync.WaitGroup{}
	for _, url := range data.Urls {
		wg.Add(1)
		go misc.FetchImageToForward(forwarder, url, wg)
	}
	wg.Wait()
	return qq.SendGroupForwardMessage(forwarder)
}

var (
	cosplaySingleCommand   = command.NewNode([]string{"single", "一张"}, "一张随机的 Cosplayer 图片", false, cosplaySingle)
	cosplayMultipleCommand = command.NewNode([]string{"multiple", "多张"}, "多张随机的 Cosplayer 图片", false, cosplayMultiple)
)

var cosplayCommand = command.NewParent([]string{"cosplay", "coser", "cos", "角色扮演"}, "Cosplayer 图片指令",
	cosplaySingleCommand,
	cosplayMultipleCommand,
)

func init() {
	command.AddCommand(cosplayCommand)
}
