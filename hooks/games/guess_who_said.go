package games

import (
	"fmt"
	"strconv"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/game"
)

type guessWhoSaid struct {
	totalQuestions int
	scores         map[int64]int
	failed         int
	maxFailed      int
	currentSender  *message.Sender
}

func (g *guessWhoSaid) Start(args []string) error {

	g.scores = make(map[int64]int)
	g.failed = 0
	g.maxFailed = 7
	g.currentSender = nil
	g.totalQuestions = 0

	if len(args) > 0 {
		max, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("%s 不是有效的数字", args[0])
		}
		g.maxFailed = max
	}

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextLn("我将说一句话，你们猜猜是谁发的，直接回复我 @ta 就行"))
	defer g.sendNextQuestion()
	return sendGameMsg(msg)
}

func (g *guessWhoSaid) Handle(msg *message.GroupMessage) *game.Result {

	reply := qq.CreateAtReply(msg)

	if g.currentSender == nil {
		reply.Append(qq.NewTextf("题目尚未出现!"))
		return game.ContinueResult
	}

	at := qq.ParseMsgContent(msg.Elements).At
	answers := make([]int64, 0)
	for _, a := range at {
		if a == bot.Instance.Uin {
			continue
		}
		answers = append(answers, a)
		break
	}

	if len(answers) == 0 {
		reply.Append(qq.NewTextf("你没有@任何人"))
		_ = sendGameMsg(reply)
		return game.ContinueResult
	} else if len(answers) > 1 {
		reply.Append(qq.NewTextf("你@太多人啦，只能@一个"))
		_ = sendGameMsg(reply)
		return game.ContinueResult
	}

	ans := answers[0]

	if ans == g.currentSender.Uin {
		reply.Append(qq.NewTextf("恭喜答对! 请听下一题"))
		_ = sendGameMsg(reply)
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
			_ = sendGameMsg(reply)
			g.sendNextQuestion()
			return game.ContinueResult
		} else {
			_ = sendGameMsg(reply)
			return g.calculateFinalResult()
		}
	}
}

func (g *guessWhoSaid) sendNextQuestion() {

	// send pre-sending message
	preSend := message.NewSendingMessage()
	preSend.Append(qq.NewTextfLn("猜猜下面的信息是谁发的 👇\n"))
	_ = sendGameMsg(preSend)

	_ = sendGameMsg(g.nextQuestion())
	g.totalQuestions += 1
}

func (g *guessWhoSaid) calculateFinalResult() *game.Result {
	winner, s := int64(0), 0
	for uid, score := range g.scores {
		if score > s {
			winner, s = uid, score
		}
	}
	result := &game.Result{EndGame: true}
	if winner > 0 {
		if member := qq.FindGroupMember(winner); member != nil {
			result.Winner = member.DisplayName()
			result.Score = s
		} else {
			logger.Warnf("找不到群成员: %d", winner)
		}
	}
	defer g.summaryScoreBoard()
	return result
}

func (g *guessWhoSaid) summaryScoreBoard() {
	summary := message.NewSendingMessage()
	summary.Append(qq.NewTextLn("各成员的分数如下:"))
	for uid, score := range g.scores {
		var userName string
		if member := qq.FindGroupMember(uid); member != nil {
			userName = member.DisplayName()
		} else {
			userName = fmt.Sprintf("(用戶: %d)", uid)
		}
		summary.Append(qq.NewTextfLn("%s: %d/%d", userName, score, g.totalQuestions))
	}
	_ = sendGameMsg(summary)
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

func (g *guessWhoSaid) ArgHints() []string {
	return []string{"失败次数限制"}
}

func init() {
	game.AddGame("谁发的", &guessWhoSaid{})
}
