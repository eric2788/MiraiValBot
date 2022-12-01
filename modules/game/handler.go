package game

import (
	"fmt"

	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
)

type (
	Handler interface {
		Start()
		Handle(msg *message.GroupMessage) Result
		Stop()
	}

	Result struct {
		EndGame bool
		Winner  *client.GroupMemberInfo
		Score   int
	}
)

var (
	games               = make(map[string]Handler)
	currentGame Handler = nil

	ContinueResult = Result{EndGame: false}
	TerminateResult = Result{EndGame: true}
)

func StartGame(name string) string {
	if currentGame != nil {
		return "有游戏已经启动"
	}
	if game, ok := games[name]; ok {
		currentGame = game
		currentGame.Start()
		return fmt.Sprintf("成功启动游戏: %s", name)
	} else {
		return fmt.Sprintf("找不到游戏: %s", name)
	}
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
