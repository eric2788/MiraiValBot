package timer_tasks

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/file"
	"github.com/eric2788/MiraiValBot/modules/timer"
	qq2 "github.com/eric2788/MiraiValBot/qq"
	"github.com/eric2788/common-utils/datetime"
	"time"
)

var (
	setting = &file.DataStorage.Setting
)

func EssenceTask(bot *bot.Bot) (err error) {

	if (*setting).LastChecked == 0 {

		file.UpdateStorage(func() {
			(*setting).LastChecked = time.Now().Unix()
		})

	} else {
		duration := datetime.Duration((*setting).LastChecked, time.Now().Unix())

		// 檢測間隔少於一天
		if duration.Hours() < 24 {
			return
		}

		file.UpdateStorage(func() {
			(*setting).LastChecked = time.Now().Unix()
		})

	}

	logger.Infof("正在檢查 %s 的今天有無群精華消息被設置...", tellTime())

	essences, err := bot.GetGroupEssenceMsgList(qq2.ValGroupInfo.Uin)

	if err != nil {
		return
	}

	if len(essences) == 0 {
		logger.Infof("群精華消息列表為空，已略過")
		return
	}

	todayEssences := make([]client.GroupDigest, 0)

	for _, ess := range essences {
		if getCompareFunc()(ess.AddDigestTime) {
			todayEssences = append(todayEssences, ess)
		}
	}

	if len(todayEssences) == 0 {
		logger.Info("今天沒有被設置的群精華消息，已略過。")
		return
	}

	msg := message.NewSendingMessage()
	msg.Append(qq2.NewTextf("%s 的今天，共有 %d 则群精华消息被设置", tellTime(), len(todayEssences)))
	_ = qq2.SendGroupMessage(msg)

	for _, essence := range todayEssences {
		msg = message.NewSendingMessage()
		msg.Append(qq2.NewTextfLn("%s 设置了一则由 %s 所发送的消息为群精华消息: ", essence.AddDigestNick, essence.SenderNick))
		essenceMsg, msgErr := qq2.GetGroupMessage(qq2.ValGroupInfo.Uin, int64(essence.MessageID))

		if msgErr != nil || essenceMsg == nil {
			msg.Append(qq2.NewTextf("获取消息失败: %v", msgErr))
		} else if len(essenceMsg.Elements) == 0 {
			msg.Append(qq2.NewTextf("获取消息 %d 失败: 消息为空", essence.MessageID))
		} else {
			for _, element := range essenceMsg.Elements {
				// 不要回復元素
				if _, ok := element.(*message.ReplyElement); ok {
					continue
				}
				msg.Append(element)
			}
		}
		_ = qq2.SendGroupMessage(msg)
	}

	return
}

func getCompareFunc() func(int64) bool {
	if (*setting).YearlyCheck {
		return compareYearly
	} else {
		return compareMonthly
	}
}

func compareYearly(ts int64) bool {
	that := datetime.FromSeconds(ts)
	now := time.Now()
	return that.Day() == now.Day() && that.Month() == now.Month() && that.Year() != now.Year()
}

func compareMonthly(ts int64) bool {
	that := datetime.FromSeconds(ts)
	now := time.Now()
	return that.Day() == now.Day() && !(that.Year() == now.Year() && that.Month() == now.Month())
}

func tellTime() string {
	if (*setting).YearlyCheck {
		return "上年度"
	} else {
		return "上个月"
	}
}

func init() {
	timer.RegisterTimer("essence.ref", time.Hour*24, EssenceTask)
}
