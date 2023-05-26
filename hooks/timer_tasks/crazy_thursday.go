package timer_tasks

import (
	"math/rand"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/timer"
	"github.com/eric2788/MiraiValBot/services/copywriting"
)

func crazyThursday(bot *bot.Bot) (err error) {

	now := time.Now()

	if now.Weekday() != time.Thursday {
		return
	}

	rand.Seed(now.UnixNano())

	list, err := copywriting.GetCrazyThursdayList()
	if err != nil {
		return err
	}

	random := list[rand.Intn(len(list))]

	msg := message.NewSendingMessage()
	msg.Append(message.NewText(random))

	return qq.SendWithRandomRiskyStrategy(msg)
}

func init() {
	timer.RegisterTimer("crazy.thursday", time.Hour*24, crazyThursday)
}
