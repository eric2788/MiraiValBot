package timer_tasks

import (
	"math/rand"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/timer"
)

func RandomChat(bot *bot.Bot) error {

	rand.Seed(time.Now().UnixMicro())

	// 随机略过
	if rand.Intn(2) == 0 {
		return nil
	}

	random, err := qq.GetRandomGroupMessage(qq.ValGroupInfo.Code)
	if err != nil {
		return err
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
		return nil
	}

	return qq.SendGroupMessage(send)
}

func init() {
	timer.RegisterTimer("random.chat", time.Minute*20, RandomChat)
}
