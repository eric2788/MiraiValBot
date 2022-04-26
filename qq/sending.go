package qq

import (
	"fmt"
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"runtime/debug"
	"time"
)

type Reason uint8

const (
	Muted Reason = iota
	Nil
	Risked
)

type MessageSendError struct {
	Msg    string
	Reason Reason
}

func (m *MessageSendError) Error() string {
	return m.Msg
}

func SendGroupMessage(msg *message.SendingMessage) error {
	if ValGroupInfo == nil {
		return fmt.Errorf("群资料尚未加载。")
	}
	return SendGroupMessageByGroup(ValGroupInfo.Uin, msg)
}

func SendGroupMessageByGroup(gp int64, msg *message.SendingMessage) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("致命错误 => %v from %s", recovered, debug.Stack())
		}
		if err != nil {
			logger.Errorf("向群 %d 發送訊息時出現錯誤: %v", gp, err)
			logger.Errorf("厡訊息: %s", ParseMsgContent(msg.Elements))
		}
	}()

	if msg == nil || bot.Instance == nil {
		err = &MessageSendError{
			Msg:    "讯息或机器人为 NULL",
			Reason: Nil,
		}
		return
	}

	if IsMuted(bot.Instance.Uin) {
		err = &MessageSendError{
			Msg:    fmt.Sprintf("机器人在群 %d 被禁言，无法发送消息", gp),
			Reason: Muted,
		}
		return
	}

	result := bot.Instance.SendGroupMessage(gp, msg)

	if result == nil || result.Id == -1 {
		err = &MessageSendError{
			Msg:    "群消息发送失败，该消息可能被风控",
			Reason: Risked,
		}
	}
	return
}

func SendPrivateMessage(uid int64, msg *message.SendingMessage) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf(fmt.Sprintf("recovered while sending private message: %v", recovered))
		}
		if err != nil {
			logger.Error(err)
		}
	}()

	if msg == nil || bot.Instance == nil {
		err = &MessageSendError{
			Msg:    "讯息或机器人为 NULL",
			Reason: Nil,
		}
		return
	}

	result := bot.Instance.SendPrivateMessage(uid, msg)

	if result == nil || result.Id == -1 {
		err = &MessageSendError{
			Msg:    "私人消息发送失败，帐号可能被风控",
			Reason: Risked,
		}
	}
	return
}

func SendGroupTempMessage(gp int64, uid int64, msg *message.SendingMessage) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf(fmt.Sprintf("recovered while sending group temp message: %v", recovered))
		}
		if err != nil {
			logger.Error(err)
		}
	}()

	if msg == nil || bot.Instance == nil {
		err = &MessageSendError{
			Msg:    "讯息或机器人为 NULL",
			Reason: Nil,
		}
		return
	}

	result := bot.Instance.SendGroupTempMessage(gp, uid, msg)

	if result == nil || result.Id == -1 {
		err = &MessageSendError{
			Msg:    "临时会话消息发送失败，帐号可能被风控",
			Reason: Risked,
		}
	}
	return
}

func retry(maxTry int, seconds int64, do func(int) error, catch func(error) error, stillRiskFunc func()) {
	try, stillRisky := 0, true
	for try < maxTry {
		if err := do(try); err != nil {
			if catch(err) != nil {
				logger.Warnf("執行重試操作時出现錯誤，现正等候 %d 秒后重新发送 (第 %d 次重试)", seconds, try+1)
				<-time.After(time.Second * time.Duration(seconds))
				try += 1
			} else {
				stillRisky = false
				break
			}
		} else {
			stillRisky = false
			break
		}
	}
	if stillRisky {
		logger.Errorf("尝试执行 %d 次后依然有錯誤，放弃执行。", try)
		if stillRiskFunc != nil {
			stillRiskFunc()
		}
	}
}

// SendRiskyMessage 发送风控几率大的消息並实行重试机制
func SendRiskyMessage(maxTry int, seconds int64, f func(currentTry int) error) {
	retry(maxTry, seconds, f, func(err error) error {
		if sendErr, ok := err.(*MessageSendError); ok && sendErr.Reason == Risked {
			logger.Warnf("嘗試发送消息時出現風控: %v", err)
			return err
		} else {
			return nil
		}
	}, nil)
}
