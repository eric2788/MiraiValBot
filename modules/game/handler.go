package game

import (
	"fmt"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"golang.org/x/exp/maps"
)

type (
	Handler interface {
		Start(args []string) error
		Handle(msg *message.GroupMessage) *Result
		Stop()
		ArgHints() []string
	}

	Result struct {
		EndGame bool
		Winner  string
		Score   int
	}
)

var (
	games               = make(map[string]Handler)
	currentGame Handler = nil

	ContinueResult  = &Result{EndGame: false}
	TerminateResult = &Result{EndGame: true}
)

func StartGame(name string, args ...string) string {
	if currentGame != nil {
		return "有游戏已经启动"
	}
	if game, ok := games[name]; ok {
		currentGame = game
		err := currentGame.Start(args)
		if err != nil {
			return fmt.Sprintf("启动游戏 %s 失败: %v", name, err)
		}
	} else {
		return fmt.Sprintf("找不到游戏: %s, 可用游戏: %v", name, strings.Join(maps.Keys(games), ", "))
	}

	return ""
}

func ListGames() []string {
	names := make([]string, 0, len(games))
	for name, g := range games {
		names = append(names, fmt.Sprintf("%s => 可用参数: %s", name, strings.Join(g.ArgHints(), ", ")))
	}
	return names
}

func IsInGame() bool {
	return currentGame != nil
}

func StopGame() string {
	if currentGame == nil {
		return "没有游戏被启动"
	}
	currentGame.Stop()
	currentGame = nil
	return "成功停止游戏"
}

func AddGame(name string, handler Handler) {
	games[name] = handler
	logger.Infof("成功添加游戏 %s", name)
}
