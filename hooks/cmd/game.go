package cmd

import (
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

var (
	startGameCommand = command.NewNode([]string{"start", "开始", "启动"}, "开始一个游戏", false, startGame, "<游戏名称>")
	stopGameCommand  = command.NewNode([]string{"stop", "中止", "关闭"}, "中止目前游戏", false, stopGame)
)

var gameCommand = command.NewParent([]string{"parent", "游戏"}, "文字游戏指令")

func init() {
	command.AddCommand(gameCommand)
}
