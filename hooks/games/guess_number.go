package games

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/game"
)

type guessNumber struct {
	min       int
	max       int
	maxFailed int
	failed    int
	random *rand.Rand

	guess int
}

func (g *guessNumber) ArgHints() []string {
	return []string{"min", "max", "失败次数"}
}

func (g *guessNumber) Start(args []string) error {

	g.failed = 0
	g.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	g.min = 1
	g.max = 100
	g.maxFailed = 7

	if len(args) > 0 {
		min, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("%s 不是有效的数字", args[0])
		}
		g.min = min
	}

	if len(args) > 1 {
		max, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("%s 不是有效的数字", args[1])
		}
		g.max = max
	}

	if len(args) > 2 {
		maxFailed, err := strconv.Atoi(args[2])
		if err != nil {
			return fmt.Errorf("%s 不是有效的数字", args[2])
		}
		g.maxFailed = maxFailed
	}

	g.guess = g.random.Intn(g.max) + g.min

	msg := message.NewSendingMessage()
	msg.Append(qq.NewTextf("猜 %d ~ %d 内的一个数字，最多可以猜 %d 次，@我回答!", g.min, g.max, g.maxFailed))
	return qq.SendGroupMessage(msg)
}

func (g *guessNumber) Handle(msg *message.GroupMessage) *game.Result {
	reply := qq.CreateAtReply(msg)
	txt := strings.TrimSpace(qq.ParseMsgContent(msg.Elements).Texts[0])

	guess, err := strconv.Atoi(txt)
	if err != nil {
		reply.Append(qq.NewTextf("%q 不是有效的数字", txt))
		_ = qq.SendGroupMessage(reply)
		return game.ContinueResult
	}

	if guess == g.guess {
		reply.Append(qq.NewTextf("答案正确!!"))
		_ = qq.SendGroupMessage(reply)

		return &game.Result{
			EndGame: true,
			Winner:  msg.Sender.DisplayName(),
		}
	}

	g.failed++
	if g.failed >= g.maxFailed {
		reply.Append(qq.NewTextf("回答错误且回答次数已用完。 正确答案是: %d", g.guess))
		_ = qq.SendGroupMessage(reply)
		return game.TerminateResult
	} else {
		if guess > g.guess {
			g.max = guess
		} else if guess < g.guess {
			g.min = guess
		}
		reply.Append(qq.NewTextf("回答错误, 范围缩小到 %d ~ %d, 回答次数剩余 %d/%d", g.min, g.max, g.maxFailed-g.failed, g.maxFailed))
		_ = qq.SendGroupMessage(reply)
		return game.ContinueResult
	}
}

func (g *guessNumber) Stop() {
}

func init() {
	game.AddGame("猜数字", &guessNumber{})
}
