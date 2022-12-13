package timer_tasks

import (
	"math/rand"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/timer"
	"github.com/eric2788/MiraiValBot/utils/misc"
)

func randomChat(bot *bot.Bot) error {

	rand.Seed(time.Now().UnixNano())

	// 随机略过
	if rand.Intn(2) == 0 {
		return nil
	} else if h := time.Now().Hour(); h > 2 && h < 7 { // 凌晨两点都七点期间暂停屁话发送
		return nil
	}

	var getMsg func() (*message.SendingMessage, error)

	// 70% 发送群图片, 30% 发送群消息
	if rand.Intn(100)+1 > 70 {
		getMsg = misc.NewRandomMessage
	} else {
		getMsg = misc.NewRandomImage
	}

	if msg, err := getMsg(); err != nil {
		return err
	} else {
		return qq.SendGroupMessage(msg)
	}
}

func init() {
	timer.RegisterTimer("random.chat", time.Minute*20, randomChat)
}
