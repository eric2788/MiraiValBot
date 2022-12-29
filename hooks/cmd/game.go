package cmd

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/eric2788/MiraiValBot/internal/qq"
	"github.com/eric2788/MiraiValBot/modules/command"
	"github.com/eric2788/MiraiValBot/modules/game"
)

func startGame(args []string, source *command.MessageSource) error {
	name := args[0]
	msg := qq.CreateReply(source.Message)
	res := game.StartGame(name, args[1:]...)
	if res == "" {
		return nil
	}
	msg.Append(qq.NewTextf(res))
	return qq.SendGroupMessage(msg)
}

func stopGame(args []string, source *command.MessageSource) error {
	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextf(game.StopGame()))
	return qq.SendGroupMessage(msg)
}

func listGames(args []string, source *command.MessageSource) error {
	msg := qq.CreateReply(source.Message)
	msg.Append(qq.NewTextf("可用游戏 + 参数: \n%s", strings.Join(game.ListGames(), "\n")))
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

func addPoint(args []string, source *command.MessageSource) (err error) {
	pt, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("无效数字: %s", args[0])
	}
	if pt < 1 {
		return fmt.Errorf("数字必须大于 0")
	}
	user, display := source.Message.Sender.Uin, source.Message.Sender.DisplayName()
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) > 0 {
		user, display = ats[0].Target, ats[0].Display
	}
	game.DepositPoint(user, pt)
	return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(qq.NewTextf("成功给 %d 添加 %d 点", display, pt)))
}

func listPoint(args []string, source *command.MessageSource) (err error) {
	user, display := source.Message.Sender.Uin, source.Message.Sender.DisplayName()
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) > 0 {
		user, display = ats[0].Target, ats[0].Display
	}
	return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(qq.NewTextf("%d 点数: %d", display, game.GetPoint(user))))
}

func removePoint(args []string, source *command.MessageSource) (err error) {
	pt, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("无效数字: %s", args[0])
	}
	if pt < 1 {
		return fmt.Errorf("数字必须大于 0")
	}
	user, display := source.Message.Sender.Uin, source.Message.Sender.DisplayName()
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) > 0 {
		user, display = ats[0].Target, ats[0].Display
	}
	msg := qq.CreateReply(source.Message)
	if game.WithdrawPoint(user, pt) {
		msg.Append(qq.NewTextf("成功从 %d 扣除 %d 点", display, pt))
	} else {
		msg.Append(qq.NewTextf("扣除失败, %d 点数不足", display))
	}
	return qq.SendGroupMessage(msg)
}

func setPoint(args []string, source *command.MessageSource) (err error) {
	pt, err := strconv.ParseInt(args[0], 10, 64)
	if err != nil {
		return fmt.Errorf("无效数字: %s", args[0])
	}
	user, display := source.Message.Sender.Uin, source.Message.Sender.DisplayName()
	ats := qq.ExtractMessageElement[*message.AtElement](source.Message.Elements)
	if len(ats) > 0 {
		user, display = ats[0].Target, ats[0].Display
	}
	game.SetPoint(user, pt)
	return qq.SendGroupMessage(qq.CreateReply(source.Message).Append(qq.NewTextf("成功将 %d 点数设置为 %d", display, pt)))
}

var (
	startGameCommand   = command.NewNode([]string{"start", "开始", "启动"}, "开始一个游戏", false, startGame, "<游戏名称>", "[参数]")
	stopGameCommand    = command.NewNode([]string{"stop", "中止", "关闭"}, "中止目前游戏", false, stopGame)
	diceCommand        = command.NewNode([]string{"dice", "骰子", "掷骰子"}, "掷骰子", false, dice, "[数字]")
	guessFingerCommand = command.NewNode([]string{"finger", "剪刀石头布", "出拳"}, "剪刀石头布", false, guessFinger, "[剪刀/石头/布]")
	listGameCommand    = command.NewNode([]string{"list", "游戏列表"}, "可用游戏列表+参数", false, listGames)
	pointsCommand      = command.NewParent([]string{
		"points", "point", "点数", "积分", "金币",
	}, "点数管理指令",
		command.NewNode([]string{"add", "添加"}, "添加点数", true, addPoint, "<点数>", "[@用户]"),
		command.NewNode([]string{"remove", "扣除"}, "扣除点数", true, removePoint, "<点数>", "[@用户]"),
		command.NewNode([]string{"set", "设置"}, "设置点数", true, setPoint, "<点数>", "[@用户]"),
		command.NewNode([]string{"list", "查看"}, "查看点数", false, listPoint, "[@用户]"),
	)
)

var gameCommand = command.NewParent([]string{"game", "游戏"}, "文字游戏指令",
	startGameCommand,
	stopGameCommand,
	diceCommand,
	guessFingerCommand,
	listGameCommand,
	pointsCommand,
)

func init() {
	command.AddCommand(gameCommand)
}
