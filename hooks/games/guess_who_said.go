package games

import (
	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/game"
)

type guessWhoSaid struct {
	scores        map[int64]int
	failed        int
	maxFailed     int
	currentSender *message.Sender
}

func (g *guessWhoSaid) Start() {

	g.scores = make(map[int64]int)
	g.failed = 0
	g.maxFailed = 7
	g.currentSender = nil

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextLn("我将说一句话，你们猜猜是谁发的，直接回复我 @ta 就行"))
	_ = qq.SendGroupMessage(msg)

	g.sendNextQuestion()
}

func (g *guessWhoSaid) Handle(msg *message.GroupMessage) game.Result {

	reply := qq.CreateAtReply(msg)

	if g.currentSender == nil {
		reply.Append(qq.NewTextf("题目尚未出现!"))
		return game.ContinueResult
	}

	at := qq.ParseMsgContent(msg.Elements).At
	ans := int64(0)
	for _, a := range at {
		if a == bot.Instance.Uin {
			continue
		}
		ans = a
		break
	}

	if ans == 0 {
		reply.Append(qq.NewTextf("你没有@任何人"))
		_ = qq.SendGroupMessage(reply)
		return game.ContinueResult
	}

	if ans == g.currentSender.Uin {
		reply.Append(qq.NewTextf("恭喜答对! 请听下一题"))
		_ = qq.SendGroupMessage(reply)
		if score, ok := g.scores[msg.Sender.Uin]; ok {
			g.scores[msg.Sender.Uin] = score + 1
		} else {
			g.scores[msg.Sender.Uin] = 1
		}
		g.sendNextQuestion()
		return game.ContinueResult
	} else {
		reply.Append(qq.NewTextfLn("答错了，正确答案是 %q", g.currentSender.DisplayName()))
		g.failed += 1
		if g.maxFailed-g.failed > 0 {
			reply.Append(qq.NewTextf("你群还有 %d/%d 次机会", g.maxFailed-g.failed, g.maxFailed))
			_ = qq.SendGroupMessage(reply)
			g.sendNextQuestion()
			return game.ContinueResult
		} else {
			_ = qq.SendGroupMessage(reply)
			return g.calculateFinalResult()
		}
	}
}

func (g *guessWhoSaid) sendNextQuestion() {
	_ = qq.SendGroupMessage(message.NewSendingMessage().Append(qq.NewTextf("猜猜下面的信息是谁发的:")))
	_ = qq.SendGroupMessage(g.nextQuestion())
}

func (g *guessWhoSaid) calculateFinalResult() game.Result {
	winner, s := int64(0), 0
	for uid, score := range g.scores {
		if score > s {
			winner = uid
		}
	}
	result := game.TerminateResult
	if winner == 0 {
		if member := qq.FindGroupMember(winner); member != nil {
			result.Winner = member
			result.Score = s
		} else {
			logger.Warnf("找不到群成员: %d", winner)
		}
	}
	return result
}

func (g *guessWhoSaid) nextQuestion() *message.SendingMessage {
	random, err := qq.GetRandomGroupMessage(qq.ValGroupInfo.Code)
	if err != nil {
		logger.Errorf("获取随机信息出现错误: %v, 正在重新获取...", err)
		return g.nextQuestion()
	}
	msg := message.NewSendingMessage()
	for _, ele := range random.Elements {
		if _, ok := ele.(*message.ReplyElement); ok {
			continue
		}

		if _, ok := ele.(*message.ForwardElement); ok {
			continue
		}

		msg.Append(ele)
	}

	if len(msg.Elements) == 0 {
		logger.Warnf("获取到的随机消息为空, 正在重新获取...")
		return g.nextQuestion()
	}

	g.currentSender = random.Sender
	return msg
}

func (g *guessWhoSaid) Stop() {

}

func init() {
	game.AddGame("谁发的", &guessWhoSaid{})
}
