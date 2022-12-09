package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/modules/game"
)

func startGame(args []string, source *command.MessageSource) error {
	name := args[0]
	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextf(game.StartGame(name)))
	return qq.SendGroupMessage(msg)
}

func stopGame(args []string, source *command.MessageSource) error {
	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextf(game.StopGame()))
	return qq.SendGroupMessage(msg)
}

func dice(args []string, source *command.MessageSource) (err error) {

	guess := -1
	if len(args) > 0 {
		guess, err = strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("无效数字: %s", args[0])
		} else if guess < 1 || guess > 6 {
			return fmt.Errorf("数字必须在 1 ~ 6 之间")
		}
	}

	rand.Seed(time.Now().UnixNano())

	actual := rand.Intn(6) + 1
	msg := message.NewSendingMessage()
	msg.Append(message.NewDice(int32(actual)))

	defer func() {
		if guess < 0 {
			return
		}
		txt := "猜对了!"
		if actual != guess {
			txt = "猜错咯"
		}
		<-time.After(time.Second * 2)
		err = qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText(txt)))
	}()

	return qq.SendGroupMessage(msg)
}

var fingerGuessingId = map[string]int32{
	"石头": 0,
	"剪刀": 1,
	"布":  2,
}

func guessFinger(args []string, source *command.MessageSource) (err error) {

	guess := int32(-1)
	if len(args) > 0 {
		if num, ok := fingerGuessingId[args[0]]; ok {
			guess = num
		} else {
			return fmt.Errorf("无效的选项: %s", args[0])
		}
	}

	rand.Seed(time.Now().UnixNano())

	actual := int32(rand.Intn(3))

	
	msg := message.NewSendingMessage()
	msg.Append(message.NewFingerGuessing(actual))

	defer func() {
		if guess < 0 {
			return
		}
		var txt string
		if actual == guess {
			txt = "这波是平手"
		} else {
			win := false
			switch guess {
			case 0:
				win = actual == 1
			case 1:
				win = actual == 2
			case 2:
				win = actual == 0
			default:
				return
			}

			if win {
				txt = "你赢了"
			} else {
				txt = "你输了"
			}

		}
		<-time.After(time.Second * 2)
		err = qq.SendGroupMessage(qq.CreateReply(source.Message).Append(message.NewText(txt)))
	}()

	return qq.SendGroupMessage(msg)
}

var (
	startGameCommand   = command.NewNode([]string{"start", "开始", "启动"}, "开始一个游戏", false, startGame, "<游戏名称>")
	stopGameCommand    = command.NewNode([]string{"stop", "中止", "关闭"}, "中止目前游戏", false, stopGame)
	diceCommand        = command.NewNode([]string{"dice", "骰子", "掷骰子"}, "掷骰子", false, dice, "[数字]")
	guessFingerCommand = command.NewNode([]string{"finger", "剪刀石头布", "出拳"}, "剪刀石头布", false, guessFinger, "[剪刀/石头/布]")
)

var gameCommand = command.NewParent([]string{"game", "游戏"}, "文字游戏指令",
	startGameCommand,
	stopGameCommand,
	diceCommand,
	guessFingerCommand,
)

func init() {
	command.AddCommand(gameCommand)
}
