package timer_tasks

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/timer"
)

func RandomChat(bot *bot.Bot) error {

	rand.Seed(time.Now().UnixNano())

	// 随机略过
	if rand.Intn(2) == 0 {
		return nil
	}

	rand.Seed(time.Now().UnixNano())

	var getMsg func() (*message.SendingMessage, error)

	// 70% 发送群图片, 30% 发送群消息
	if rand.Intn(100) > 70 {
		getMsg = getRandomImage
	} else {
		getMsg = getRandomMessage
	}

	if msg, err := getMsg(); err != nil {
		return err
	} else {
		return qq.SendGroupMessage(msg)
	}
}

func init() {
	timer.RegisterTimer("random.chat", time.Minute*20, RandomChat)
}

func getRandomMessage() (*message.SendingMessage, error) {
	random, err := qq.GetRandomGroupMessage(qq.ValGroupInfo.Code)
	if err != nil {
		return nil, err
	}

	send := message.NewSendingMessage()

	for _, ele := range random.Elements {

		switch ele.(type) {
		case *message.ReplyElement:
			continue
		case *message.ForwardElement:
			continue
		default:
			break
		}
		send.Append(ele)
	}

	// 没有元素也略过
	if len(send.Elements) == 0 {
		return nil, fmt.Errorf("讯息元素为空。")
	}

	return send, nil
}

func getRandomImage() (*message.SendingMessage, error) {
	rand.Seed(time.Now().UnixMicro())
	imgs := qq.GetImageList()

	if len(imgs) == 0 {
		return nil, fmt.Errorf("群图片缓存列表为空。")
	}

	logger.Debugf("成功索取 %d 张群图片缓存。", len(imgs))

	chosen := imgs[rand.Intn(len(imgs))]

	b, err := qq.GetCacheImage(chosen)
	if err != nil {
		return nil, err
	}
	img, err := qq.NewImageByByte(b)
	if err != nil {
		return nil, err
	}
	return message.NewSendingMessage().Append(img), nil
}
