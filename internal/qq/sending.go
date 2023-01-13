package qq

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
)

type Reason uint8

const (
	Muted Reason = iota
	Nil
	Risked

	MaxRecallTime time.Duration = time.Minute*2 - time.Second*5 // 提前五秒撤回
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

func SendGroupMessageAndRecall(msg *message.SendingMessage, duration time.Duration) error {
	if ValGroupInfo == nil {
		return fmt.Errorf("群资料尚未加载。")
	}
	return SendGroupMessageAndRecallByGroup(ValGroupInfo.Uin, msg, duration)
}

func SendGroupMessageAndRecallByGroup(gp int64, msg *message.SendingMessage, duration time.Duration) (err error) {

	// 超過兩分鐘就無法撤回消息了
	if duration > MaxRecallTime {
		logger.Warnf("撤回消息的時間超過兩分鐘，將會無法撤回消息, 已改為兩分鐘。")
		duration = MaxRecallTime // 提前五秒撤回
	}

	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("致命错误 => %v", recovered)
			debug.PrintStack()
		}
		if err != nil {
			logger.Errorf("向群 %d 發送訊息時出現錯誤: %v", gp, err)
			logger.Errorf("厡訊息: %s", ParseMsgContent(msg.Elements))
		}
	}()

	msg.Append(NewTextf("(本条消息将在 %.f 秒后撤回)", duration.Seconds()))

	result, err := sendGroupMessage(gp, msg)

	if err != nil {
		return
	}

	if result == nil || result.Id == -1 {
		err = &MessageSendError{
			Msg:    "群消息发送失败，该消息可能被风控",
			Reason: Risked,
		}
		return
	}

	go func() {
		time.Sleep(duration)
		err = bot.Instance.RecallGroupMessage(gp, result.Id, result.InternalId)
		if err != nil {
			logger.Errorf("撤回群消息时出现错误: %v", err)
			reply := CreateReply(result)
			reply.Append(NewTextf("撤回失败: %v", err))
			_ = SendGroupMessageByGroup(gp, reply)
		}
	}()

	return
}

func SendGroupForwardMessage(msg *message.ForwardMessage) error {
	if ValGroupInfo == nil {
		return fmt.Errorf("群资料尚未加载。")
	}
	return SendGroupForwardMessageByGroup(ValGroupInfo.Uin, msg)
}

func SendGroupForwardMessageByGroup(gp int64, msg *message.ForwardMessage) (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("致命错误 => %v", recovered)
			debug.PrintStack()
		}
		if err != nil {
			logger.Errorf("向群 %d 發送合并轉發訊息時出現錯誤: %v", gp, err)
			logger.Errorf("厡訊息: %s", msg.Brief())
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
			Msg:    fmt.Sprintf("机器人在群 %d 被禁言，无法发送合并转发消息", gp),
			Reason: Muted,
		}
		return
	}

	builder := bot.Instance.NewForwardMessageBuilder(gp)
	fe := builder.Main(msg)

	if fe == nil {
		err = &MessageSendError{
			Msg:    "合并转发讯息为 NULL",
			Reason: Nil,
		}
		return
	}

	result := bot.Instance.SendGroupForwardMessage(gp, fe)

	if result == nil || result.Id == -1 {
		err = &MessageSendError{
			Msg:    "群合并转发消息发送失败，该消息可能被风控",
			Reason: Risked,
		}
	}
	return
}

func SendPrivateForwardMessage(uid int64, msg *message.ForwardMessage) (err error) {
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

	builder := bot.Instance.NewForwardMessageBuilder(0)
	fe := builder.Main(msg)

	if fe == nil {
		err = &MessageSendError{
			Msg:    "合并转发讯息为 NULL",
			Reason: Nil,
		}
		return
	}

	result := bot.Instance.SendPrivateMessage(uid, message.NewSendingMessage().Append(fe))

	if result == nil || result.Id == -1 {
		err = &MessageSendError{
			Msg:    "私人消息发送失败，帐号可能被风控",
			Reason: Risked,
		}
	}
	return
}

func SendGroupMessageByGroup(gp int64, msg *message.SendingMessage) (err error) {

	defer func() {
		if recovered := recover(); recovered != nil {
			err = fmt.Errorf("致命错误 => %v", recovered)
			debug.PrintStack()
		}
		if err != nil {
			logger.Errorf("向群 %d 發送訊息時出現錯誤: %v", gp, err)
			logger.Errorf("厡訊息: %s", ParseMsgContent(msg.Elements))
		}
	}()

	result, err := sendGroupMessage(gp, msg)

	if err != nil {
		return
	}

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
	result.ToString()
	return
}

func retry(maxTry int, seconds int64, do func(int) error, catch func(error) error, stillRiskFunc func()) {
	try, stillRisky := 0, true
	for try < maxTry+1 {
		if err := do(try); err != nil {
			if catch(err) != nil {
				try += 1
				if try == maxTry+1 {
					break
				}
				logger.Warnf("執行重試操作時出现錯誤，现正等候 %d 秒后重新发送 (第 %d 次重试)", seconds, try)
				<-time.After(time.Second * time.Duration(seconds))
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
		logger.Errorf("尝试执行 %d 次后依然有錯誤，放弃执行。", try-1)
		if stillRiskFunc != nil {
			stillRiskFunc()
		}
	}
}

// SendRiskyMessage 发送风控几率大的消息並实行重试机制
func SendRiskyMessage(maxTry int, seconds int64, f func(currentTry int) error) {
	SendRiskyMessageWithFunc(maxTry, seconds, f, nil)
}

// SendRiskyMessageWithFunc 发送风控几率大的消息並实行重试机制，並在重试失败后执行回调函数
func SendRiskyMessageWithFunc(maxTry int, seconds int64, f func(currentTry int) error, stillRiskFunc func()) {
	retry(maxTry, seconds, f, func(err error) error {
		if sendErr, ok := err.(*MessageSendError); ok && sendErr.Reason == Risked {
			logger.Warnf("嘗試发送消息時出現風控: %v", err)
			return err
		} else {
			return nil
		}
	}, stillRiskFunc)
}

func SendWithRandomRiskyFunc(msg *message.SendingMessage, stillRisky func()) (err error) {
	go SendRiskyMessageWithFunc(5, 60, func(try int) error {
		clone := CloneMessage(msg)
		alt := GetRandomMessageByTry(try)
		if len(alt) > 0 {
			clone.Append(NextLn())
			for _, element := range alt {
				clone.Append(element)
				clone.Append(NextLn())
			}
		}
		return SendGroupMessage(clone)
	}, stillRisky)
	return
}

func SendWithRandomRiskyStrategy(msg *message.SendingMessage) (err error) {
	return SendWithRandomRiskyFunc(msg, nil)
}

func SendWithRandomRiskyStrategyRemind(msg *message.SendingMessage, source *message.GroupMessage) (err error) {
	return SendWithRandomRiskyFunc(msg, func() {
		// 重试失败后，提示信息被风控
		remind := CreateAtReply(source)
		remind.Append(message.NewText("回应发送失败，可能被风控咯 😔"))
		_ = SendGroupMessageByGroup(source.GroupCode, remind)
	})
}

func CloneMessage(msg *message.SendingMessage) *message.SendingMessage {
	clone := message.NewSendingMessage()
	for _, element := range msg.Elements {
		clone.Append(element)
	}
	return clone
}

func GetRandomMessageByTry(try int) []*message.TextElement {

	extras := make([]*message.TextElement, 0)

	// 新增随机发过的群消息

	if try > 0 {

		random, err := GetRandomGroupMessage(ValGroupInfo.Uin)

		// 大于 2 次重试时，根据重试次数发送更多随机消息
		for i := 0; i < try-2; i++ {
			extras = append(extras, GetRandomMessageByTry(1)...)
		}

		if err == nil && random != nil {

			for _, element := range random.Elements {
				switch e := element.(type) {
				case *message.TextElement:
					extras = append(extras, e)
				case *message.AtElement:
					extras = append(extras, message.NewText(e.Display))
				case *message.FaceElement:
					extras = append(extras, message.NewText(e.Name))
				default:
					break
				}
			}

			// 随机消息没有文本
			if len(extras) == 0 {

				logger.Warnf("为被风控的广播插入一条新消息再发送: %s", random.ToString())

				sendFirst := message.NewSendingMessage()
				for _, element := range random.Elements {

					switch element.(type) {
					case *message.ReplyElement:
						continue
					case *message.ForwardElement:
						continue
					default:
						break
					}

					sendFirst.Append(element)
				}
				_ = SendGroupMessage(sendFirst)
				<-time.After(time.Second * 5)     // 发送完等待五秒
				return GetRandomMessageByTry(try) // 再獲取一則隨機消息

			} else {

				logger.Warnf("为被风控的广播新增如下的内容: %s", random.ToString())

			}

		} else { // 随机消息获取失败

			if err != nil {
				logger.Warnf("获取随机消息时出现错误: %v, 将改为发送风控次数", err)
			} else if random == nil {
				logger.Warnf("获取随机消息时出现错误: 訊息為 nil , 将改为发送风控次数")
			}

			// 则发送风控次数?
			extras = append(extras, NewTextf("此广播已被风控 %d 次 QAQ!!", try))

		}

	}

	return extras
}

func sendGroupMessage(gp int64, msg *message.SendingMessage) (result *message.GroupMessage, err error) {
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
	result = bot.Instance.SendGroupMessage(gp, msg)
	return
}
