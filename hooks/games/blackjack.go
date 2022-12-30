package games

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/game"
)

// create a game of blackjack

var (
	cards = []string{
		"A", "2", "3", "4", "5", "6", "7", "8", "9", "10", "J", "Q", "K",
	}
	suits = []string{
		"♠", "♥", "♣", "♦",
	}
)

type blackjack struct {
	cards  map[int64][]string
	bet    map[int64]int64
	ran    *rand.Rand
	joined [6]*client.GroupMemberInfo
	ctx    context.Context
	stop   context.CancelFunc
	turn   int
	raised bool
}

func (p *blackjack) Start(args []string) error {

	p.ran = rand.New(rand.NewSource(time.Now().UnixNano()))
	p.joined = [6]*client.GroupMemberInfo{}
	p.cards = make(map[int64][]string)
	p.bet = make(map[int64]int64)
	p.turn = -1

	// bot joined the game
	p.joined[5] = qq.FindGroupMember(bot.Instance.Uin)
	if p.joined[5] == nil {
		return fmt.Errorf("机器人不在瓦群内")
	}

	p.ctx, p.stop = context.WithTimeout(context.Background(), time.Second*30)

	go func() {
		<-p.ctx.Done()
		p.stop()
		if p.joined[0] == nil {
			reply := message.NewSendingMessage()
			reply.Append(qq.NewTextfLn("人数不足"))
			reply.Append(qq.NewTextfLn(game.StopGame()))
			_ = qq.SendGroupMessage(reply)
			return
		}
		p.gameStart()
	}()

	sending := message.NewSendingMessage()
	sending.Append(qq.NewTextfLn("三十秒后开始21点，@我输入 加入 [筹码] 参与游戏 (默认筹码为100点)"))
	return qq.SendGroupMessage(sending)
}

func (p *blackjack) Handle(msg *message.GroupMessage) *game.Result {
	args := strings.Split(strings.TrimSpace(strings.Join(qq.ParseMsgContent(msg.Elements).Texts, " ")), " ")
	reply := qq.CreateAtReply(msg)
	if len(args) == 0 {
		reply.Append(qq.NewTextf("你在港咩也?"))
		_ = qq.SendGroupMessage(reply)
		return game.ContinueResult
	}
	select {
	case <-p.ctx.Done():
		return p.handleOption(args, msg)
	default:
		res := p.handleGameJoin(args, msg)
		reply.Append(message.NewText(res))
		_ = qq.SendGroupMessage(reply)
	}

	return game.ContinueResult
}

func (p *blackjack) handleOption(args []string, msg *message.GroupMessage) *game.Result {
	reply := qq.CreateAtReply(msg)

	if p.joined[p.turn].Uin != msg.Sender.Uin {
		reply.Append(qq.NewTextf("现在不是你的回合"))
		_ = qq.SendGroupMessage(reply)
		return game.ContinueResult
	}

	if args[0] == "叫牌" {
		if len(args) > 1 && args[1] == "翻倍" {
			if p.raised {
				reply.Append(qq.NewTextf("叫牌后无法再加注, 请重新输入"))
				_ = qq.SendGroupMessage(reply)
				return game.ContinueResult
			} else {
				if game.WithdrawPoint(msg.Sender.Uin, p.bet[msg.Sender.Uin]) {
					p.bet[msg.Sender.Uin] *= 2
					reply.Append(qq.NewTextf("加注成功, 当前筹码: %d", p.bet[msg.Sender.Uin]))
				} else {
					reply.Append(qq.NewTextf("加注失败, 你的点数不足"))
				}
				_ = qq.SendGroupMessage(reply)
				reply = qq.CreateAtReply(msg)
			}
		}
		p.raised = true
		card := p.pickOneCardFor(msg.Sender.Uin)
		reply.Append(qq.NewTextfLn("你叫了一张牌: %v", card))
		score := p.caculatePoints(msg.Sender.Uin)
		if score > 21 {
			reply.Append(qq.NewTextf("你的点数已超过21点(%d), 鉴定为爆牌", score))
			_ = qq.SendGroupMessage(reply)
			return p.nextTurnResult()
		} else if score == 21 {
			reply.Append(qq.NewTextf("你的点数为21点, 鉴定为黑杰克"))
			_ = qq.SendGroupMessage(reply)
			return p.endGame()
		} else {
			reply.Append(qq.NewTextf("你的点数目前为%d点", score))
			_ = qq.SendGroupMessage(reply)
		}
	} else if args[0] == "停牌" {
		reply.Append(qq.NewTextf("你停牌了"))
		_ = qq.SendGroupMessage(reply)
		return p.nextTurnResult()
	} else {
		reply.Append(qq.NewTextf("未知操作类型: %v, 可用操作: 叫牌 [翻倍], 停牌", args[0]))
		_ = qq.SendGroupMessage(reply)
	}
	return game.ContinueResult
}

func (p *blackjack) nextTurn() bool {
	p.turn++
	if p.turn >= len(p.joined) {
		return false
	} else if p.joined[p.turn] == nil {
		return p.nextTurn()
	}
	return true
}

func (p *blackjack) nextTurnResult() *game.Result {
	if p.nextTurn() {
		turner := p.joined[p.turn]
		if turner.Uin == bot.Instance.Uin {
			return p.botTurn()
		}
		reply := message.NewSendingMessage()
		reply.Append(message.NewText("现在轮到 "))
		reply.Append(message.NewAt(turner.Uin, turner.DisplayName()))
		reply.Append(message.NewText(" 的回合, 请输入操作: 叫牌 或 叫牌 翻倍 (叫牌前双倍加注), 停牌"))
		_ = qq.SendGroupMessage(reply)
		p.raised = false // 重设加注状态
		return game.ContinueResult
	}
	return p.endGame()
}

func (p *blackjack) botTurn() *game.Result {
	reply := message.NewSendingMessage()
	reply.Append(message.NewText("现在轮到庄家的回合"))
	_ = qq.SendGroupMessage(reply)

	// 庄家叫牌
	for {
		score := p.caculatePoints(bot.Instance.Uin)
		if score > 21 {
			break
		} else if score == 21 {
			break
		} else if score >= 17 {
			break
		}
		card := p.pickOneCardFor(bot.Instance.Uin)
		reply := message.NewSendingMessage()
		reply.Append(qq.NewTextf("庄家叫了一张牌: %v", card))
		_ = qq.SendGroupMessage(reply)
	}

	// 庄家停牌
	reply = message.NewSendingMessage()
	reply.Append(message.NewText("庄家停牌"))
	_ = qq.SendGroupMessage(reply)
	return p.nextTurnResult()
}

func (p *blackjack) endGame() *game.Result {
	result := message.NewSendingMessage()
	// 先计算庄家的点数
	ownerScore := p.caculatePoints(bot.Instance.Uin)
	result.Append(qq.NewTextfLn("游戏结果如下:"))
	result.Append(qq.NewTextfLn("庄家的点数为 %d [%s]", ownerScore, strings.Join(p.cards[bot.Instance.Uin], " | ")))
	if ownerScore > 21 {
		result.Append(qq.NewTextfLn("(庄家爆牌)"))
		ownerScore = 0
	}

	for _, v := range p.joined {
		if v == nil {
			continue
		} else if v.Uin == bot.Instance.Uin { // 庄家额外计算
			continue
		}

		pt := p.caculatePoints(v.Uin)
		if pt > 21 {
			// lose
			p.bet[v.Uin] = 0
			result.Append(qq.NewTextfLn("%v 爆牌, 现有筹码为 %d (x0) [ %s ]", v.DisplayName(), p.bet[v.Uin], strings.Join(p.cards[v.Uin], " | ")))
		} else if pt == 21 {
			// if ownerScore == 21 , draw
			if ownerScore < 21 {
				// win
				p.bet[v.Uin] *= 3
				result.Append(qq.NewTextfLn("%v 为黑杰克, 现有筹码为 %d (x3) [ %s ]", v.DisplayName(), p.bet[v.Uin], strings.Join(p.cards[v.Uin], " | ")))
			} else {
				// draw
				result.Append(qq.NewTextfLn("%v 和庄家都是黑杰克, 现有筹码为 %d (不变) [ %s ]", v.DisplayName(), p.bet[v.Uin], strings.Join(p.cards[v.Uin], " | ")))
			}
		} else {
			if pt > ownerScore {
				p.bet[v.Uin] *= 2
				result.Append(qq.NewTextfLn("%v 赢过庄家, 现有筹码为 %d (x2) [ %s ]", v.DisplayName(), p.bet[v.Uin], strings.Join(p.cards[v.Uin], " | ")))
			} else {
				p.bet[v.Uin] = 0
				result.Append(qq.NewTextfLn("%v 输给庄家, 现有筹码为 %d (x0) [ %s ]", v.DisplayName(), p.bet[v.Uin], strings.Join(p.cards[v.Uin], " | ")))
			}
		}
	}
	_ = qq.SendGroupMessage(result)
	p.returnBets()
	// restart again and wait for next game
	p.Start(nil)
	return game.ContinueResult
}

func (p *blackjack) handleGameJoin(args []string, msg *message.GroupMessage) string {
	if args[0] == "加入" {
		if _, ok := p.bet[msg.Sender.Uin]; ok {
			return "你已经加入了游戏"
		}
		balance := int64(100)
		if len(args) > 1 {
			var err error
			balance, err = strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				return fmt.Sprintf("無效的數字 %v: %v", args[1], err)
			}
			if balance < 0 {
				return "筹码必须大于0"
			}
		}
		for i, v := range p.joined {
			if v == nil {
				p.joined[i] = qq.FindGroupMember(msg.Sender.Uin)
				if p.joined[i] == nil {
					return "你不在瓦群内"
				} else if !game.WithdrawPoint(msg.Sender.Uin, balance) {
					return fmt.Sprintf("你的点数不足%d，无法转换筹码", balance)
				}
				p.bet[msg.Sender.Uin] = balance
				if p.playerFull() {
					p.stop()
				}
				return fmt.Sprintf("加入成功, 已转换 %d 点数 为 筹码", balance)
			}
		}
		return "人满了"
	} else if args[0] == "退出" {
		for i, v := range p.joined {
			if v != nil && v.Uin == msg.Sender.Uin {
				p.joined[i] = nil
				delete(p.bet, msg.Sender.Uin)
				game.DepositPoint(msg.Sender.Uin, 100)
				return "退出成功, 已退还点数 100"
			}
		}
		return "你没加入游戏"
	}

	return fmt.Sprintf("未知操作: %v, 可用操作: 加入, 退出", args[0])
}

func (p *blackjack) gameStart() {
	reply := message.NewSendingMessage()
	reply.Append(qq.NewTextf("正在开始发牌...."))
	_ = qq.SendGroupMessage(reply)

	reply = message.NewSendingMessage()
	// pick 2 cards for each player
	for _, v := range p.joined {
		if v == nil {
			continue
		}
		p.pickOneCardFor(v.Uin)
		p.pickOneCardFor(v.Uin)

		if v.Uin == bot.Instance.Uin {
			reply.Append(qq.NewTextfLn("庄家的牌: [ %s | ? ]", p.cards[v.Uin][0]))
		} else {
			reply.Append(qq.NewTextfLn("%s 的牌: [%s] 为 %d 点", v.DisplayName(), strings.Join(p.cards[v.Uin], " | "), p.caculatePoints(v.Uin)))
		}

	}
	reply.Append(qq.NewTextfLn("发牌结束, 开始回合"))
	_ = qq.SendGroupMessage(reply)
	p.nextTurnResult()
}

func (p *blackjack) Stop() {
	if p.stop != nil {
		p.stop()
	}
}

func (p *blackjack) returnBets() {
	for uid, bet := range p.bet {
		game.DepositPoint(uid, bet)
	}
	_ = qq.SendGroupMessage(message.NewSendingMessage().Append(message.NewText("回合结束，已退还所有筹码至点数")))
}

func (p *blackjack) ArgHints() []string {
	return nil
}

func (p *blackjack) pickOneCardFor(user int64) string {
	// pick a card
	card := cards[p.ran.Intn(len(cards))]
	suit := suits[p.ran.Intn(len(suits))]
	// add to user's cards
	p.cards[user] = append(p.cards[user], card+suit)
	return card + suit
}

func (p *blackjack) caculatePoints(user int64) uint8 {
	points := uint8(0)
	aces := 0
	for _, v := range p.cards[user] {
		// remove suit
		ru := []rune(v)
		v = string(ru[:len(ru)-1])
		switch v {
		case "A":
			aces += 1
		case "J", "Q", "K", "10":
			points += 10
		default:
			points += uint8(v[0] - '0')
		}
	}
	for i := 0; i < aces; i++ {
		if points+11 <= 21 {
			points += 11
		} else {
			points += 1
		}
	}
	return points
}

func (p *blackjack) playerFull() bool {
	for _, v := range p.joined {
		if v == nil {
			return false
		}
	}
	return true
}

func init() {
	game.AddGame("21点", &blackjack{})
}
