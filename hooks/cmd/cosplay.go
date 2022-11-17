package cmd

import (
	"errors"
	"sync"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/services/cosplayer"
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
		go fetchImageToForward(forwarder, url, wg)
	}
	wg.Wait()
	return qq.SendGroupForwardMessage(forwarder)
}

var (
	cosplaySingleCommand   = command.NewNode([]string{"single", "一张"}, "一张随机的 Cosplayer 图片", false, cosplaySingle)
	cosplayMultipleCommand = command.NewNode([]string{"multiple", "多张"}, "多张随机的 Cosplayer 图片", false, cosplayMultiple)
)

var cosplayCommand = command.NewParent([]string{"cosplay", "角色扮演"}, "Cosplayer 图片指令",
	cosplaySingleCommand,
	cosplayMultipleCommand,
)

func init() {
	command.AddCommand(cosplayCommand)
}

func fetchImageToForward(forwarder *message.ForwardMessage, url string, wg *sync.WaitGroup) {
	defer wg.Done()
	msg := message.NewSendingMessage()
	img, err := qq.NewImageByUrl(url)
	if err != nil {
		logger.Errorf("尝试获取图片 %s 失败: %v, 将使用URL链接。", url, err)
		msg.Append(qq.NewTextf("[图片获取失败, 原链接: %s]", url))
	} else {
		msg.Append(img)
	}
	forwarder.AddNode(qq.NewForwardNode(msg))
}
